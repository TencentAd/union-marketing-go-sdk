package ocean_engine

import (
	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/account"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/config"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/http_tools"
)

type AccountService struct {
	config     *config.Config
	httpClient *http_tools.HttpClient
}

// NewOceanEngineAccountService
func NewAccountService(config *config.Config) *AccountService {
	s := &AccountService{
		config:     config,
		httpClient: http_tools.Init(config.HttpConfig),
	}
	return s
}

func (s *AccountService) GetAuthAccount(input *sdk.BaseInput) (*sdk.AuthAccount, error) {
	id := formatAuthAccountID(input.AccountId)
	return account.GetAuthAccount(id)
}