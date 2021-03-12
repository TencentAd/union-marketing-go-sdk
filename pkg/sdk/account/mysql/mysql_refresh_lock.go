package mysql

import "github.com/tencentad/union-marketing-go-sdk/pkg/sdk/orm"

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
