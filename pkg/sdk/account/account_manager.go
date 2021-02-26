package account

import (
	"fmt"
	"time"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	log "github.com/sirupsen/logrus"
)

var m *manager

var (
	getTokenRefreshTimeMethods = make(map[sdk.MarketingPlatformType]sdk.GetTokenRefreshTime)
	refreshTokenMethods        = make(map[sdk.MarketingPlatformType]sdk.RefreshToken)
)

func init() {
	m = &manager{
		cache: newCache(),
	}
}

// RegisterGetTokenRefreshTime 注册GetTokenRefreshTime方法
func RegisterGetTokenRefreshTime(t sdk.MarketingPlatformType, f sdk.GetTokenRefreshTime) {
	getTokenRefreshTimeMethods[t] = f
}

func RegisterRefreshToken(t sdk.MarketingPlatformType, f sdk.RefreshToken) {
	refreshTokenMethods[t] = f
}

// isExpired 判断账号过期时间，只判断token过期
func isExpired(account *sdk.AuthAccount) bool {
	f, ok := getTokenRefreshTimeMethods[account.Platform]
	if ok {
		return f(account).After(time.Now())
	} else {
		return sdk.GetRefreshTimeDefault(account).After(time.Now())
	}
}

func refreshToken(account *sdk.AuthAccount) (*sdk.RefreshTokenOutput, error) {
	if f, ok := refreshTokenMethods[account.Platform]; !ok {
		return nil, fmt.Errorf("platform[%s] not register refreshToken method", account.Platform)
	} else {
		return f(account)
	}
}

// manager 授权账号管理器
type manager struct {
	cache   *cache
	storage Storage
	lock    RefreshLock
}

// Init 初始化账号管理器
func Init(storage Storage, lock RefreshLock) error {
	m.storage = storage
	m.lock = lock
	return m.init()
}

// init 初始化token管理器
func (m *manager) init() error {
	authAccounts, err := m.storage.List()
	if err != nil {
		return err
	}

	for _, info := range authAccounts {
		m.cache.insert(info)
	}

	go refreshRoutine()

	return nil
}

// Upsert 插入或者更新授权账户
func Insert(authAccount *sdk.AuthAccount) error {
	m.cache.insert(authAccount)
	if m.storage != nil {
		return m.storage.Upsert(authAccount)
	}
	return nil
}

// GetAuthAccount 获取授权账号
func GetAuthAccount(id string) (*sdk.AuthAccount, error) {
	account := m.cache.get(id)
	if account != nil && !isExpired(account) {
		return account, nil
	} else {
		dbAccount, err := m.storage.Take(id)
		if err != nil {
			return nil, err
		}

		m.cache.insert(dbAccount)
		return dbAccount, nil
	}
}

func refreshRoutine() {
	for {
		time.Sleep(time.Second * 10)
		if err := m.lock.Lock(); err != nil {
			continue
		}
		refresh()
		if err := m.lock.Unlock(); err != nil {
			log.Errorf("failed to unlock, err: %v", err)
		}
	}
}

func refresh() {
	account, err := m.storage.List()

	if err != nil {
		log.Errorf("failed to list all account, err: %v", err)
		return
	}

	for _, a := range account {
		if isExpired(a) {
			out, err := refreshToken(a)
			if err != nil {
				log.Errorf("failed to refresh token, err: %v", err)
				continue
			}
			if err = applyRefreshTokenOutput(out); err != nil {
				log.Errorf("failed to apply refreshed token, err: %v", err)
			}
		}
	}
}

// RefreshToken refresh后，更新token
func applyRefreshTokenOutput(out *sdk.RefreshTokenOutput) error {
	if m.storage != nil {
		return m.storage.UpdateToken(out)
	}

	return nil
}
