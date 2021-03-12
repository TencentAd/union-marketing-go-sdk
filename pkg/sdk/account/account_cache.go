package account

import (
	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	"github.com/modern-go/concurrent"
)

type cache struct {
	m *concurrent.Map
}

func newCache() *cache {
	return &cache{
		m: concurrent.NewMap(),
	}
}

func (c *cache) insert(authAccount *sdk.AuthAccount) {
	c.m.Store(authAccount.ID, authAccount)
}

func (c *cache) get(id string) *sdk.AuthAccount {
	account, ok := c.m.Load(id)
	if !ok {
		return nil
	} else {
		return account.(*sdk.AuthAccount)
	}
}

func (c *cache) getAll() []*sdk.AuthAccount {
	ret := make([]*sdk.AuthAccount, 0)
	c.m.Range(func(key, value interface{}) bool {
		ret = append(ret, value.(*sdk.AuthAccount))
		return true
	})
	return ret
}
