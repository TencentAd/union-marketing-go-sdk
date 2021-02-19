package marketingsdk

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/marketingsdk/oceanengine"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/marketingsdk/tencent"
)

type MarketingSDK struct {
	// config

	TencentSDKService *tencent.TencentMarketingService
	OceanEngineSDKService *oceanengine.OceanEngineReport
	//Map[string, service]
}

func (msc *MarketingSDK) GetReport(reqParam *api.ReportInputParam) (*api.ReportOutput, error) {
	// accountId, accountType, AccessToken
	// tencent getReport

	// oceanengine getReport
	return nil, nil
}