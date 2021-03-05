package ams

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/http_tools"
)

type AMSAccountService struct {
	config     *config.Config
	httpClient *http_tools.HttpClient
}

// NewAccountService
func NewAccountService(config *config.Config) *AMSAccountService {
	s := &AMSAccountService{
		config:     config,
		httpClient: http_tools.Init(config.HttpConfig),
	}
	return s
}

func (s *AMSAccountService) GetAuthAccount(id string) (*sdk.AuthAccount, error) {
	return account.GetAuthAccount(id)
}
