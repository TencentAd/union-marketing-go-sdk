package ams

import (
	"fmt"
	"strings"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	sdkconfig "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	tapi "github.com/tencentad/marketing-api-go-sdk/pkg/api"
	tconfig "github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
)

type AMSReportService struct {
	config *sdkconfig.Config
}

func NewAMSReportService(sConfig *sdkconfig.Config) *AMSReportService {
	return &AMSReportService{
		config: sConfig,
	}
}

// GetReport 获取报表接口
func (t *AMSReportService) GetReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	if reportInput.ReportTimeGranularity == sdk.ReportTimeDaily {
		return t.getDailyReport(reportInput)
	} else {
		return t.getHourlyReport(reportInput)
	}
}

// TencentReport constructor
func (t *AMSReportService) getAMSReportClient(reportInput *sdk.GetReportInput) *ads.SDKClient {
	tSdkConfig := tconfig.SDKConfig{
		AccessToken: reportInput.BaseInput.AccessToken,
	}
	tClient := ads.Init(&tSdkConfig)
	tClient.UseProduction()
	return tClient
}

const TFilterMax = 5

func (t *AMSReportService) getReportAdLevel(reportInput *sdk.GetReportInput, adLevel *string) (bool, error) {

	if reportInput.BaseInput.AccountType <= sdk.AccountTypeInvalid || reportInput.BaseInput.AccountType >= sdk.AccountTypeMax {
		return false, fmt.Errorf("getReportAdLevel invalid account type = %d, id = %d", reportInput.BaseInput.AccountType, reportInput.BaseInput.AccountId)
	}
	if reportInput.BaseInput.AccountType == sdk.AccountTypeTencent {
		switch reportInput.ReportAdLevel {
		case sdk.LevelAccount:
			*adLevel = "REPORT_LEVEL_ADVERTISER"
			break
		case sdk.LevelCampaign:
			*adLevel = "REPORT_LEVEL_CAMPAIGN"
			break
		case sdk.LevelAd:
			*adLevel = "REPORT_LEVEL_ADGROUP"
			break
		case sdk.LevelCreative:
			*adLevel = "REPORT_LEVEL_AD"
			break
		default:
			*adLevel = ""
			return false, nil
		}
		return true, nil
	} else {
		switch reportInput.ReportAdLevel {
		case sdk.LevelAccount:
			*adLevel = "REPORT_LEVEL_ADVERTISER_WECHAT"
			break
		case sdk.LevelCampaign:
			*adLevel = "REPORT_LEVEL_CAMPAIGN_WECHAT"
			break
		case sdk.LevelAd:
			*adLevel = "REPORT_LEVEL_ADGROUP_WECHAT"
			break
		case sdk.LevelCreative:
			*adLevel = "REPORT_LEVEL_AD_WECHAT"
			break
		default:
			*adLevel = ""
			return false, nil
		}
		return true, nil
	}
}

// getDailyReport 获取天级别的广告数据
func (t *AMSReportService) getDailyReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	tClient := t.getAMSReportClient(reportInput)
	var level string
	isSucc, err := t.getReportAdLevel(reportInput, &level)
	if !isSucc {
		return nil, err
	}
	dateRange := model.ReportDateRange{
		StartDate: reportInput.ReportDateRange.StartDate,
		EndDate:   reportInput.ReportDateRange.EndDate,
	}

	var dailyReportsGetOpts tapi.DailyReportsGetOpts
	// Filtering
	if reportInput.ReportFiltering != nil {
		TFiltering := make([]model.FilteringStruct, 0, TFilterMax)
		// TODO 验证 Operator: "IN"是否可以完全替代EQUALS
		// campaign_id
		mFiltering := reportInput.ReportFiltering.(sdk.ReportFiltering)
		if len(mFiltering.CampaignIDList) > 0 {
			TFiltering = append(TFiltering, model.FilteringStruct{
				Field:    "campaign_id",
				Operator: "IN",
				Values:   &mFiltering.CampaignIDList,
			})
		}
		// adgroup_id
		if len(mFiltering.GroupIDList) > 0 {
			TFiltering = append(TFiltering, model.FilteringStruct{
				Field:    "adgroup_id",
				Operator: "IN",
				Values:   &mFiltering.CampaignIDList,
			})
		}
		// ad_id
		if len(mFiltering.CreativeIDList) > 0 {
			TFiltering = append(TFiltering, model.FilteringStruct{
				Field:    "ad_id",
				Operator: "IN",
				Values:   &mFiltering.CreativeIDList,
			})
		}
		dailyReportsGetOpts.Filtering = optional.NewInterface(TFiltering)
	}

	// GroupBy 逗号分割
	if len(reportInput.ReportGroupBy) > 0 {
		dailyReportsGetOpts.GroupBy = optional.NewInterface(strings.Split(string(reportInput.ReportGroupBy), ","))
	}

	// OrderBy
	dailyReportsGetOpts.OrderBy = optional.NewInterface([]model.OrderByStruct{{
		reportInput.ReportOrderBy.SortField,
		model.Sortord(string(reportInput.ReportOrderBy.SortType)),
	}})

	// Page,Page_size
	dailyReportsGetOpts.Page = reportInput.Page
	dailyReportsGetOpts.PageSize = reportInput.PageSize

	// Fields
	dailyReportsGetOpts.Fields = optional.NewInterface(reportInput.Fields)

	// 获取天级别广告数据
	result, _, err := tClient.DailyReports().Get(*tClient.Ctx, reportInput.BaseInput.AccountId, level, dateRange, &dailyReportsGetOpts)
	response := sdk.GetReportOutput{
		TencentReportResponse: &result,
	}

	return &response, err
}

/**
获取小时级别的广告数据
*/
func (t *AMSReportService) getHourlyReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	tClient := t.getAMSReportClient(reportInput)
	var level string
	isSucc, err := t.getReportAdLevel(reportInput, &level)
	if !isSucc {
		return nil, err
	}
	dateRange := model.DateRange{
		StartDate: reportInput.ReportDateRange.StartDate,
		EndDate:   reportInput.ReportDateRange.EndDate,
	}

	var hourlyReportsGetOpts tapi.HourlyReportsGetOpts
	// Filtering
	if reportInput.ReportFiltering != nil {
		TFiltering := make([]model.FilteringStruct, 0, TFilterMax)
		// TODO 验证 Operator: "IN"是否可以完全替代EQUALS
		// campaign_id
		mFiltering := reportInput.ReportFiltering.(sdk.ReportFiltering)
		if len(mFiltering.CampaignIDList) > 0 {
			TFiltering = append(TFiltering, model.FilteringStruct{
				Field:    "campaign_id",
				Operator: "IN",
				Values:   &mFiltering.CampaignIDList,
			})
		}
		// adgroup_id
		if len(mFiltering.GroupIDList) > 0 {
			TFiltering = append(TFiltering, model.FilteringStruct{
				Field:    "adgroup_id",
				Operator: "IN",
				Values:   &mFiltering.CampaignIDList,
			})
		}
		// ad_id
		if len(mFiltering.CreativeIDList) > 0 {
			TFiltering = append(TFiltering, model.FilteringStruct{
				Field:    "ad_id",
				Operator: "IN",
				Values:   &mFiltering.CreativeIDList,
			})
		}
		hourlyReportsGetOpts.Filtering = optional.NewInterface(TFiltering)
	}

	// GroupBy 逗号分割
	if len(reportInput.ReportGroupBy) > 0 {
		hourlyReportsGetOpts.GroupBy = optional.NewInterface(strings.Split(string(reportInput.ReportGroupBy), ","))
	}

	// OrderBy
	hourlyReportsGetOpts.OrderBy = optional.NewInterface([]model.OrderByStruct{{
		reportInput.ReportOrderBy.SortField,
		model.Sortord(string(reportInput.ReportOrderBy.SortType)),
	}})

	// Page,Page_size
	hourlyReportsGetOpts.Page = reportInput.Page
	hourlyReportsGetOpts.PageSize = reportInput.PageSize

	// Fields
	hourlyReportsGetOpts.Fields = optional.NewInterface(reportInput.Fields)

	// 获取天级别广告数据
	result, _, err := tClient.HourlyReports().Get(*tClient.Ctx, reportInput.BaseInput.AccountId, level, dateRange, &hourlyReportsGetOpts)
	response := sdk.GetReportOutput{
		TencentHourlyReportResponse: &result,
	}

	return &response, err
}
