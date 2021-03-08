package ams

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	sdkconfig "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
)

// AMService AMS处理服务
type AMService struct {
	config *sdkconfig.Config

	*AMSAccountService // 账户模块
	*ReportService     // 报表模块
	*AuthService
	*MaterialService // 物料管理模块
}

// Name 名称
func (t *AMService) Name() sdk.MarketingPlatformType {
	return sdk.AMS
}

func (t *AMService) GetConfig() *sdkconfig.Config {
	return t.config
}

// NewAMSService 创建AMS服务
func NewAMSService(config *sdkconfig.Config) *AMService {
	return &AMService{
		config:             config,
		ReportService:   NewReportService(config),
		MaterialService: NewMaterialService(config),
		AuthService:        NewAuthService(config),
	}
}
