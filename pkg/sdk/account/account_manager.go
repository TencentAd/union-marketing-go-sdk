package account

import "git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"

var m *manager

// manager 授权账号管理器
type manager struct {
	cache   *cache
	storage Storage
}

// Init 初始化账号管理器
func Init(storage Storage) error {
	m = newManager(storage)
	return m.init()
}

func newManager(storage Storage) *manager {
	cache := newCache()
	return &manager{
		cache:   cache,
		storage: storage,
	}
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

// RefreshToken refresh后，更新token
func RefreshToken(authAccount *sdk.AuthAccount) error {
	current, err := m.cache.refreshToken(authAccount)
	if err != nil {
		return err
	}
	if m.storage != nil {
		return m.storage.Update(current)
	}

	return nil
}

// GetAuthAccount 获取授权账号
func GetAuthAccount(accountID int64) *sdk.AuthAccount {
	return m.cache.get(accountID)
}

// GetAllAuthAccount 获取所有授权账号
func  GetAllAuthAccount() []*sdk.AuthAccount {
	return m.cache.getAll()
}