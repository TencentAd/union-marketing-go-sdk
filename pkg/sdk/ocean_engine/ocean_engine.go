package ocean_engine

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	config "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/http_tools"
)

type OceanEngineService struct {
	config           *config.Config
	*AccountService  // 账户模块
	*CampaignService  // 计划模块
	*AdGroupService  // 广告组模块
	*ReportService   // 报表模块
	*AuthService     // 权限模块
	*MaterialService // 物料管理模块
}

// Name 名称
func (s *OceanEngineService) Name() sdk.MarketingPlatformType {
	return sdk.OceanEngine
}

func (s *OceanEngineService) GetConfig() *config.Config {
	return s.config
}

// NewAMSService 创建AMS服务
func NewOceanEngineService(config *config.Config) *OceanEngineService {
	config.HttpConfig = &http_tools.HttpConfig{
		BasePath: "https://ad.oceanengine.com/open_api",
	}
	return &OceanEngineService{
		config:          config,
		AccountService:  NewAccountService(config),
		CampaignService: NewCampaignService(config),
		AdGroupService: NewAdGroupService(config),
		ReportService:   NewReportService(config),
		AuthService:     NewAuthService(config),
		MaterialService: NewMaterialService(config),
	}

}
