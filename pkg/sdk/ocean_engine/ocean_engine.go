package ocean_engine

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	sdkconfig "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
)

type OceanEngineService struct {
	config                    *sdkconfig.Config
	*OceanEngineReportService // 报表模块
	*AuthService
	//*OceanEngineMaterialService // 物料管理模块
}

// Name 名称
func (t *OceanEngineService) Name() sdk.MarketingPlatformType {
	return sdk.OceanEngine
}

func (t *OceanEngineService) GetConfig() *sdkconfig.Config {
	return t.config
}

// NewAMSService 创建AMS服务
func NewOceanEngineService(sconfig *sdkconfig.Config) *OceanEngineService {
	return &OceanEngineService{
		config: sconfig,
		//AMSReportService:   NewAMSReportService(sconfig),
		//AMSMaterialService: NewAMSMaterialService(sconfig),
		AuthService: NewAuthService(sconfig),
	}
}
