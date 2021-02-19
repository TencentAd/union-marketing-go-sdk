package tencent

import (
	"fmt"
	"git.code.oa.com/tme-server-component/kg_growth_open/api"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	tapi "github.com/tencentad/marketing-api-go-sdk/pkg/api"
	tconfig "github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
	"strings"
)

type TencentReportService struct {
}

// TencentReport constructor
func (t *TencentReportService) TencentReportInit(sdkConfig *api.MarketingSDKConfig, reqParam *api.ReportInputParam) *ads.SDKClient {
	tSdkConfig := tconfig.SDKConfig{
		Configuration: tconfig.Configuration{
			BasePath:      sdkConfig.Configuration.BasePath,
			Host:          sdkConfig.Configuration.Host,
			Scheme:        sdkConfig.Configuration.Scheme,
			DefaultHeader: sdkConfig.Configuration.DefaultHeader,
			UserAgent:     sdkConfig.Configuration.UserAgent,
			HTTPClient:    sdkConfig.Configuration.HTTPClient,
		},
		AccessToken:  reqParam.AccessToken,
		IsDebug:      sdkConfig.IsDebug,
		DebugFile:    sdkConfig.DebugFile,
		SkipMonitor:  sdkConfig.SkipMonitor,
		IsStrictMode: sdkConfig.IsStrictMode,
		GlobalConfig: tconfig.GlobalConfig{
			ServiceName: tconfig.ServiceName{
				Name:   sdkConfig.GlobalConfig.ServiceName.Name,
				Schema: sdkConfig.GlobalConfig.ServiceName.Schema,
			},
			HttpOption: tconfig.HttpOption{
				Header: sdkConfig.GlobalConfig.HttpOption.Header,
			},
		},
	}
	tClient := ads.Init(&tSdkConfig)
	tClient.UseProduction()
	return tClient
}

const TFilterMax = 5

func (t *TencentReportService) getReportAdLevel(reqParam *api.ReportInputParam, adLevel *string) (bool, error) {

	if reqParam.AccountType <= api.AccountTypeInvalid || reqParam.AccountType >= api.AccountTypeMax {
		return false, fmt.Errorf("newtencentreport invalid account type = %d, id = %d", reqParam.AccountType, reqParam.AccountId)
	}
	if reqParam.AccountType == api.AccountTypeTencent {
		switch reqParam.ReportAdLevel {
		case api.LevelAccount:
			*adLevel = "REPORT_LEVEL_ADVERTISER"
			break
		case api.LevelCampaign:
			*adLevel = "REPORT_LEVEL_CAMPAIGN"
			break
		case api.LevelAd:
			*adLevel = "REPORT_LEVEL_ADGROUP"
			break
		case api.LevelCreative:
			*adLevel = "REPORT_LEVEL_AD"
			break
		default:
			*adLevel = ""
			return false, nil
		}
		return true, nil
	} else {
		switch reqParam.ReportAdLevel {
		case api.LevelAccount:
			*adLevel = "REPORT_LEVEL_ADVERTISER_WECHAT"
			break
		case api.LevelCampaign:
			*adLevel = "REPORT_LEVEL_CAMPAIGN_WECHAT"
			break
		case api.LevelAd:
			*adLevel = "REPORT_LEVEL_ADGROUP_WECHAT"
			break
		case api.LevelCreative:
			*adLevel = "REPORT_LEVEL_AD_WECHAT"
			break
		default:
			*adLevel = ""
			return false, nil
		}
		return true, nil
	}
}

/**
获取天级别的广告数据
*/
func (t *TencentReportService) getDailyReport(sdkConfig *api.MarketingSDKConfig, reqParam *api.ReportInputParam) (*api.ReportOutput, error) {
	tClient := t.TencentReportInit(sdkConfig, reqParam)
	var level string
	isSucc, err := t.getReportAdLevel(reqParam, &level)
	if !isSucc {
		return nil, err
	}
	dateRange := model.ReportDateRange{
		StartDate: reqParam.ReportDateRange.StartDate,
		EndDate:   reqParam.ReportDateRange.EndDate,
	}

	var dailyReportsGetOpts tapi.DailyReportsGetOpts
	// Filtering
	if reqParam.ReportFiltering != nil {
		TFiltering := make([]model.FilteringStruct, 0, TFilterMax)
		// TODO 验证 Operator: "IN"是否可以完全替代EQUALS
		// campaign_id
		mFiltering := reqParam.ReportFiltering.(api.ReportFiltering)
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
	if len(reqParam.ReportGroupBy) > 0 {
		dailyReportsGetOpts.GroupBy = optional.NewInterface(strings.Split(string(reqParam.ReportGroupBy), ","))
	}

	// OrderBy
	dailyReportsGetOpts.OrderBy = optional.NewInterface([]model.OrderByStruct{{
		reqParam.ReportOrderBy.SortField,
		model.Sortord(string(reqParam.ReportOrderBy.SortType)),
	}})

	// Page,Page_size
	dailyReportsGetOpts.Page = reqParam.Page
	dailyReportsGetOpts.PageSize = reqParam.PageSize

	// Fields
	dailyReportsGetOpts.Fields = optional.NewInterface(reqParam.Fields)

	// 获取天级别广告数据
	result, _, err := tClient.DailyReports().Get(*tClient.Ctx, reqParam.AccountId, level, dateRange, &dailyReportsGetOpts)
	response := api.ReportOutput{
		TencentReportResponse: &result,
	}

	return &response, err
}

/**
获取小时级别的广告数据
*/
func (t *TencentReportService) getHourlyReport(sdkConfig *api.MarketingSDKConfig, reqParam *api.ReportInputParam) (*api.ReportOutput, error) {
	tClient := t.TencentReportInit(sdkConfig, reqParam)
	var level string
	isSucc, err := t.getReportAdLevel(reqParam, &level)
	if !isSucc {
		return nil, err
	}
	dateRange := model.DateRange{
		StartDate: reqParam.ReportDateRange.StartDate,
		EndDate:   reqParam.ReportDateRange.EndDate,
	}

	var hourlyReportsGetOpts tapi.HourlyReportsGetOpts
	// Filtering
	if reqParam.ReportFiltering != nil {
		TFiltering := make([]model.FilteringStruct, 0, TFilterMax)
		// TODO 验证 Operator: "IN"是否可以完全替代EQUALS
		// campaign_id
		mFiltering := reqParam.ReportFiltering.(api.ReportFiltering)
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
	if len(reqParam.ReportGroupBy) > 0 {
		hourlyReportsGetOpts.GroupBy = optional.NewInterface(strings.Split(string(reqParam.ReportGroupBy), ","))
	}

	// OrderBy
	hourlyReportsGetOpts.OrderBy = optional.NewInterface([]model.OrderByStruct{{
		reqParam.ReportOrderBy.SortField,
		model.Sortord(string(reqParam.ReportOrderBy.SortType)),
	}})

	// Page,Page_size
	hourlyReportsGetOpts.Page = reqParam.Page
	hourlyReportsGetOpts.PageSize = reqParam.PageSize

	// Fields
	hourlyReportsGetOpts.Fields = optional.NewInterface(reqParam.Fields)

	// 获取天级别广告数据
	result, _, err := tClient.HourlyReports().Get(*tClient.Ctx, reqParam.AccountId, level, dateRange, &hourlyReportsGetOpts)
	response := api.ReportOutput{
		TencentHourlyReportResponse: &result,
	}

	return &response, err
}
