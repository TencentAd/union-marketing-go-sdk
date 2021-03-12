package ams

import (
	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	config "github.com/tencentad/union-marketing-go-sdk/pkg/sdk/config"
)

// AMService AMS处理服务
type AMService struct {
	config *config.Config

	*AccountService  // 账户模块
	*CampaignService // 计划模块
	*AdGroupService   // 广告组模块
	*ReportService   // 报表模块
	*AuthService     // 权限认证模块
	*MaterialService // 物料管理模块
}

// Name 名称
func (t *AMService) Name() sdk.MarketingPlatformType {
	return sdk.AMS
}

func (t *AMService) GetConfig() *config.Config {
	return t.config
}

// NewAMSService 创建AMS服务
func NewAMSService(config *config.Config) *AMService {
	return &AMService{
		config:          config,
		AccountService:  NewAccountService(config),
		CampaignService: NewCampaignService(config),
		AdGroupService: NewAdGroupService(config),
		ReportService:   NewReportService(config),
		MaterialService: NewMaterialService(config),
		AuthService:     NewAuthService(config),
	}
}
