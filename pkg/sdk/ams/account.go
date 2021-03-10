package ams

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
)

type AccountService struct {
	config     *config.Config
}

// NewAccountService
func NewAccountService(config *config.Config) *AccountService {
	s := &AccountService{
		config:     config,
	}
	return s
}

// GetAuthAccount 获取账户信息
func (s *AccountService) GetAuthAccount(input *sdk.BaseInput) (*sdk.AuthAccount, error) {
	id := formatAuthAccountID(input.AccountId, input.AMSSystemType)
	return account.GetAuthAccount(id)
}
