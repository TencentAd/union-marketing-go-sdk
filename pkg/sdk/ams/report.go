package ams

import (
	"fmt"
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	sdkconfig "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	tapi "github.com/tencentad/marketing-api-go-sdk/pkg/api"
	tconfig "github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
	"strings"
)

type AMSReportService struct {
	config *sdkconfig.Config
}

func NewAMSReportService(sConfig *sdkconfig.Config) *AMSReportService {
	return &AMSReportService{
		config: sConfig,
	}
}

// 获取报表接口
func (t *AMSReportService) GetReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	if reportInput.TimeGranularity == sdk.ReportTimeDaily {
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
	if reportInput.BaseInput.AccountType == sdk.AccountTypeAMS {
		switch reportInput.AdLevel {
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
		case sdk.LevelVideo:
			*adLevel = "REPORT_LEVEL_MATERIAL_VIDEO"
			break
		case sdk.LevelImage:
			*adLevel = "REPORT_LEVEL_MATERIAL_IMAGE"
			break
		default:
			*adLevel = ""
			return false, fmt.Errorf("getReportAdLevel invalid adLevel= %s", reportInput.AdLevel)
		}
		return true, nil
	} else {
		switch reportInput.AdLevel {
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
			return false, fmt.Errorf("getReportAdLevel invalid adLevel= %s", reportInput.AdLevel)
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
		StartDate: reportInput.DateRange.StartDate,
		EndDate:   reportInput.DateRange.EndDate,
	}

	var dailyReportsGetOpts tapi.DailyReportsGetOpts
	// Filtering
	if reportInput.Filtering != nil {
		TFiltering := make([]model.FilteringStruct, 0, TFilterMax)
		// TODO 验证 Operator: "IN"是否可以完全替代EQUALS
		// campaign_id
		mFiltering := reportInput.Filtering.(sdk.Filtering)
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
	if len(reportInput.GroupBy) > 0 {
		dailyReportsGetOpts.GroupBy = optional.NewInterface(strings.Split(string(reportInput.GroupBy), ","))
	}

	// OrderBy
	dailyReportsGetOpts.OrderBy = optional.NewInterface([]model.OrderByStruct{{
		reportInput.OrderBy.SortField,
		model.Sortord(string(reportInput.OrderBy.SortType)),
	}})

	// Page,Page_size
	dailyReportsGetOpts.Page = reportInput.Page
	dailyReportsGetOpts.PageSize = reportInput.PageSize

	// Fields
	dailyReportsGetOpts.Fields = optional.NewInterface(reportInput.Fields_AMS)

	// 获取天级别广告数据
	result, _, err := tClient.DailyReports().Get(*tClient.Ctx, reportInput.BaseInput.AccountId, level, dateRange, &dailyReportsGetOpts)
	if err != nil {
		return nil, err
	}
	reportOutput := &sdk.GetReportOutput{}
	t.copyDailyReportToOutput(&result, reportOutput)
	return reportOutput, err
}

func (t *AMSReportService) copyDailyReportToOutput(dailyResponseData *model.DailyReportsGetResponseData, reportOutput *sdk.GetReportOutput) {
	if len(*dailyResponseData.List) == 0 {
		return
	}
	rList := make([]sdk.ReportOutputListStruct, 0, len(*dailyResponseData.List))
	for i := 0; i < len(*dailyResponseData.List); i++ {
		dailyData := (*dailyResponseData.List)[i]
		rList = append(rList, sdk.ReportOutputListStruct{
			AccountId:            dailyData.AccountId,
			Date:                 dailyData.Date,
			ViewCount:            dailyData.ViewCount,
			DownloadCount:        dailyData.DownloadCount,
			ActivatedCount:       dailyData.ActivatedCount,
			ActivatedRate:        dailyData.ActivatedRate,
			ThousandDisplayPrice: dailyData.ThousandDisplayPrice,
			ValidClickCount:      dailyData.ValidClickCount,
			Ctr:                  dailyData.Ctr,
			Cpc:                  dailyData.Cpc,
			Cost:                 dailyData.Cost,
			KeyPageViewCost:      dailyData.KeyPageViewCost,
			CouponClickCount:     dailyData.CouponClickCount,
			CouponIssueCount:     dailyData.CouponIssueCount,
			CouponGetCount:       dailyData.CouponGetCount,
		})
	}
	reportOutput.List = &rList
	reportOutput.PageInfo = &sdk.PageConf{
		Page:        dailyResponseData.PageInfo.Page,
		PageSize:    dailyResponseData.PageInfo.PageSize,
		TotalNumber: dailyResponseData.PageInfo.TotalNumber,
		TotalPage:   dailyResponseData.PageInfo.TotalPage,
	}
}

// getHourlyReport 获取小时级别的广告数据
func (t *AMSReportService) getHourlyReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	tClient := t.getAMSReportClient(reportInput)
	var level string
	isSucc, err := t.getReportAdLevel(reportInput, &level)
	if !isSucc {
		return nil, err
	}
	dateRange := model.DateRange{
		StartDate: reportInput.DateRange.StartDate,
		EndDate:   reportInput.DateRange.EndDate,
	}

	var hourlyReportsGetOpts tapi.HourlyReportsGetOpts
	// Filtering
	if reportInput.Filtering != nil {
		TFiltering := make([]model.FilteringStruct, 0, TFilterMax)
		// TODO 验证 Operator: "IN"是否可以完全替代EQUALS
		// campaign_id
		mFiltering := reportInput.Filtering.(sdk.Filtering)
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
	if len(reportInput.GroupBy) > 0 {
		hourlyReportsGetOpts.GroupBy = optional.NewInterface(strings.Split(string(reportInput.GroupBy), ","))
	}

	// OrderBy
	hourlyReportsGetOpts.OrderBy = optional.NewInterface([]model.OrderByStruct{{
		reportInput.OrderBy.SortField,
		model.Sortord(string(reportInput.OrderBy.SortType)),
	}})

	// Page,Page_size
	hourlyReportsGetOpts.Page = reportInput.Page
	hourlyReportsGetOpts.PageSize = reportInput.PageSize

	// Fields
	hourlyReportsGetOpts.Fields = optional.NewInterface(reportInput.Fields_AMS)

	// 获取天级别广告数据
	result, _, err := tClient.HourlyReports().Get(*tClient.Ctx, reportInput.BaseInput.AccountId, level, dateRange, &hourlyReportsGetOpts)
	if err != nil {
		return nil, err
	}
	reportOutput := &sdk.GetReportOutput{}
	t.copyHourReportToOutput(reportInput.DateRange.StartDate, &result, reportOutput)
	return reportOutput, err
}

func (t *AMSReportService) copyHourReportToOutput(date string, hourResponseData *model.HourlyReportsGetResponseData, reportOutput *sdk.GetReportOutput) {
	if len(*hourResponseData.List) == 0 {
		return
	}
	rList := make([]sdk.ReportOutputListStruct, 0, len(*hourResponseData.List))
	for i := 0; i < len(*hourResponseData.List); i++ {
		hourlyData := (*hourResponseData.List)[i]
		rList = append(rList, sdk.ReportOutputListStruct{
			AccountId:            hourlyData.AccountId,
			Date:                 date,
			Hour:                 hourlyData.Hour,
			ViewCount:            hourlyData.ViewCount,
			DownloadCount:        hourlyData.DownloadCount,
			ActivatedCount:       hourlyData.ActivatedCount,
			ActivatedRate:        hourlyData.ActivatedRate,
			ThousandDisplayPrice: hourlyData.ThousandDisplayPrice,
			ValidClickCount:      hourlyData.ValidClickCount,
			Ctr:                  hourlyData.Ctr,
			Cpc:                  hourlyData.Cpc,
			Cost:                 hourlyData.Cost,
			KeyPageViewCost:      hourlyData.KeyPageViewCost,
			CouponClickCount:     hourlyData.CouponClickCount,
			CouponIssueCount:     hourlyData.CouponIssueCount,
			CouponGetCount:       hourlyData.CouponGetCount,
		})
	}
	reportOutput.List = &rList
	reportOutput.PageInfo = &sdk.PageConf{
		Page:        hourResponseData.PageInfo.Page,
		PageSize:    hourResponseData.PageInfo.PageSize,
		TotalNumber: hourResponseData.PageInfo.TotalNumber,
		TotalPage:   hourResponseData.PageInfo.TotalPage,
	}
}

func (t *AMSReportService) GetVideoReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	if reportInput.BaseInput.AccountType != sdk.AccountTypeAMS {
		return nil, fmt.Errorf("GetDailyVideoReport invalid account type = %d, id = %d", reportInput.BaseInput.AccountType, reportInput.BaseInput.AccountId)
	}
	reportInput.TimeGranularity = sdk.ReportTimeDaily
	reportInput.AdLevel = sdk.LevelVideo
	return t.getDailyReport(reportInput)
}
func (t *AMSReportService) GetImageReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	if reportInput.BaseInput.AccountType != sdk.AccountTypeAMS {
		return nil, fmt.Errorf("GetDailyVideoReport invalid account type = %d, id = %d", reportInput.BaseInput.AccountType, reportInput.BaseInput.AccountId)
	}
	reportInput.TimeGranularity = sdk.ReportTimeDaily
	reportInput.AdLevel = sdk.LevelImage
	return t.getDailyReport(reportInput)
}
