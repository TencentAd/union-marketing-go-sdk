package ams

import (
	sdkconfig "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
)

type AMService struct {
	config         *sdkconfig.Config

	*AMSReportService // 报表模块
	*AuthService
	*AMSMaterialService // 物料管理模块
}

// Name 名称
func (t *AMService) Name() string {
	return "AMS"
}

func NewAMSService(sconfig *sdkconfig.Config) *AMService {
	return &AMService{
		config:     sconfig,
		AMSReportService: NewAMSReportService(sconfig),
		AMSMaterialService: NewAMSMaterialService(sconfig),
		AuthService:      NewAuthService(sconfig),
	}
}