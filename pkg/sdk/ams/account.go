package ams

import (
	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/account"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/config"
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
