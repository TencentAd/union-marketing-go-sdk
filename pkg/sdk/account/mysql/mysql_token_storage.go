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
func (s *tokenStorage) Upsert(authAccount *sdk.AuthAccount) error {
	return orm.AuthAccountUpsert(orm.GetDB(), authAccount)
}

// UpdateToken
func (s *tokenStorage) UpdateToken(out *sdk.RefreshTokenOutput) error {
	original, err := s.Take(out.ID)
	if err != nil {
		return err
	}

	sdk.UpdateToken(original, out)

	return orm.AuthAccountUpdate(orm.GetDB(), original)
}

// List
func (s *tokenStorage) List() ([]*sdk.AuthAccount, error) {
	return orm.AuthAccountGetAll(orm.GetDB())
}

// Take
func (s *tokenStorage) Take(id string) (*sdk.AuthAccount, error) {
	return orm.AuthAccountTake(orm.GetDB(), id)
}
