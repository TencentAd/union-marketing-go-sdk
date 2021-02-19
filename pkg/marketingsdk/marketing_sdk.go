package marketingsdk

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/marketingsdk/oceanengine"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/marketingsdk/tencent"
)

type Manager struct {
	// config

	TencentSDKService *tencent.TencentMarketingService
	OceanEngineSDKService *oceanengine.OceanEngineReport
	//Map[string, service]

	Implementation map[string]api.MarketingSDK
}


func (m *Manager) Call(method string, platform string, input string) (string, error) {
	return "", nil
}

//func (msc *Manager) GetReport(reqParam *api.ReportInputParam) (*api.ReportOutput, error) {
//	// accountId, accountType, AccessToken
//	// tencent getReport
//
//	// oceanengine getReport
//	return nil, nil
//}