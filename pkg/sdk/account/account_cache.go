package account

import (
	"fmt"
	"strconv"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
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
	c.m.Store(authAccount.Key(), authAccount)
}

func (c *cache) refreshToken(authAccount *sdk.AuthAccount) (*sdk.AuthAccount, error) {
	original := c.get(authAccount.AccountId)

	if original == nil {
		return nil, fmt.Errorf("cache corrupted")
	}

	updateToken(original, authAccount)
	return original, nil
}

func (c *cache) get(accountID int64) *sdk.AuthAccount {
	account, ok := c.m.Load(strconv.FormatInt(accountID, 10))
	if !ok {
		return nil
	} else {
		return account.(*sdk.AuthAccount)
	}
}

func (c *cache) getAll() []*sdk.AuthAccount  {
	ret := make([]*sdk.AuthAccount, 0)
	c.m.Range(func(key, value interface{}) bool {
		ret = append(ret, value.(*sdk.AuthAccount))
		return true
	})
	return ret
}

func updateToken(original *sdk.AuthAccount, refreshed *sdk.AuthAccount) {
	original.RefreshTokenExpireAt = refreshed.RefreshTokenExpireAt
	original.RefreshToken = refreshed.RefreshToken
	original.AccessToken = refreshed.AccessToken
	original.AccessTokenExpireAt = refreshed.AccessTokenExpireAt
}