package account

import "git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"

var ManagerSingleton *manager

// manager 授权账号管理器
type manager struct {
	cache   *cache
	storage Storage
}

func Init(storage Storage) error {
	ManagerSingleton = newManager(storage)
	return ManagerSingleton.init()
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
func (m *manager) Insert(authAccount *sdk.AuthAccount) error {
	m.cache.insert(authAccount)
	return m.storage.Insert(authAccount)
}

func (m *manager) RefreshToken(authAccount *sdk.AuthAccount) error {
	current, err := m.cache.refreshToken(authAccount)
	if err != nil {
		return err
	}

	return m.storage.Update(current)
}

// GetAuthAccount 获取授权账号
func (m *manager) GetAuthAccount(accountID int64) *sdk.AuthAccount {
	return m.cache.get(accountID)
}

func (m *manager) GetAllAuthAccount() []*sdk.AuthAccount {
	return m.cache.getAll()
}