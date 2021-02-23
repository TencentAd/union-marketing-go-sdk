package mysql

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/orm"
)

type tokenStorage struct {
}

// NewTokenStorage
func NewTokenStorage() *tokenStorage {
	return &tokenStorage{}
}

// Upsert
func (s *tokenStorage) Ipsert(authInfo *sdk.AuthAccount) error {
	return orm.AuthAccountUpsert(orm.GetDB(), authInfo)
}

// List
func (s *tokenStorage) List() ([]*sdk.AuthAccount, error) {
	return orm.AuthAccountGetAll(orm.GetDB())
}
