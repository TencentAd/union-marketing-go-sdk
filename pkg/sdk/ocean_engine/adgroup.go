package ocean_engine

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/http_tools"
)

type AdGroupService struct {
	config     *config.Config
	httpClient *http_tools.HttpClient
}

func (s *AdGroupService) GetAdGroupList(input *sdk.AdGroupGetInput) (*sdk.AdGroupGetOutput, error) {
	panic("implement me")
}
