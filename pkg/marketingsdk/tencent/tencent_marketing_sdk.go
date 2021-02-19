package tencent

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api"
)

type TencentMarketingService struct {
	SdkConfig     *api.MarketingSDKConfig // sdk配置
	reportService *TencentReportService   // 报表模块
}

// Name 名称
func (t *TencentMarketingService) Name() string {
	return "TencentMarketingService"
}

func NewTencentMarketingService(sdkConfig *api.MarketingSDKConfig) *TencentMarketingService {
	return &TencentMarketingService{
		SdkConfig:     sdkConfig,
		reportService: &TencentReportService{},
	}
}

// 获取报表接口
func (t *TencentMarketingService) GetReport(reqParam *api.ReportInputParam) (*api.ReportOutput, error) {
	if reqParam.ReportTimeGranularity == api.ReportTimeDaily {
		return t.reportService.getDailyReport(t.SdkConfig, reqParam)
	} else {
		return t.reportService.getHourlyReport(t.SdkConfig, reqParam)
	}
}
