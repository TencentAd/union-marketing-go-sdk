package ams

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	sdkconfig "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
)

type AMService struct {
	config *sdkconfig.Config
	reportService *AMSReportService // 报表模块
}

// Name 名称
func (t *AMService) Name() string {
	return "AMS"
}

func NewAMSService(sconfig *sdkconfig.Config) *AMService {
	return &AMService{
		config:     sconfig,
		reportService: NewAMSReportService(sconfig),
	}
}

// 获取报表接口
func (t *AMService) GetReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	if reportInput.ReportTimeGranularity == sdk.ReportTimeDaily {
		return t.reportService.getDailyReport(reportInput)
	} else {
		return t.reportService.getHourlyReport(reportInput)
	}
}