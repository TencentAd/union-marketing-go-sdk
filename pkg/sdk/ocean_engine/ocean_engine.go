package ocean_engine

import (
	sdkconfig "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
)

type OceanEngineService struct {
	config *sdkconfig.Config
}

// Name 名称
func (t *OceanEngineService) Name() string {
	return "OceanEngine"
}

