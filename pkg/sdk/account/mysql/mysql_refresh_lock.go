package mysql

import "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/orm"

type RefreshLock struct {

}

const (
	resource int64 = 94234105142214
)

func NewRefreshLock() *RefreshLock {
	return &RefreshLock{}
}

func (l *RefreshLock) Lock() error {
	return orm.LockDB(orm.GetDB(), resource)
}

func (l *RefreshLock) Unlock() error {
	return orm.UnlockDB(orm.GetDB(), resource)
}
