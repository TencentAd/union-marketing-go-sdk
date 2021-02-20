package sdk

import "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/ams"

type Config struct {
	AMS *ams.Config `json:"ams"`
}
