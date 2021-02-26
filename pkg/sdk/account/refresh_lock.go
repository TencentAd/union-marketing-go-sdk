package account

// RefreshLock 刷新token加锁，防止不同服务同时刷新导致数据冲突
type RefreshLock interface {
	// Lock
	Lock() error

	// Unlock
	Unlock() error
}
