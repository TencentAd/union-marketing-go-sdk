package ams

import (
	"fmt"
	"strconv"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
	sdkconfig "git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	tapi "github.com/tencentad/marketing-api-go-sdk/pkg/api"
	tconfig "github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
)

// AMSReportService 报表服务
type AMSReportService struct {
	config *sdkconfig.Config
}

// NewAMSReportService 创建AMS报表服务
func NewAMSReportService(sConfig *sdkconfig.Config) *AMSReportService {
	return &AMSReportService{
		config: sConfig,
	}
}

// GetReport 获取报表接口
func (t *AMSReportService) GetReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	if reportInput.TimeGranularity == sdk.ReportTimeDaily {
		return t.getDailyReport(reportInput)
	} else {
		return t.getHourlyReport(reportInput)
	}
}

// Tencent
func getAMSSdkClient(authAccount *sdk.AuthAccount) *ads.SDKClient {
	tSdkConfig := tconfig.SDKConfig{
		AccessToken: authAccount.AccessToken,
		IsDebug:     true,
	}
	tClient := ads.Init(&tSdkConfig)
	tClient.UseProduction()
	return tClient
}

const TFilterMax = 5

// getReportAdLevel 获取报表adlevel
func (t *AMSReportService) getReportAdLevel(authAccount *sdk.AuthAccount, reportInput *sdk.GetReportInput, adLevel *string) (bool, error) {
	if authAccount.AMSSystemType == sdk.AMS {
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
	} else if authAccount.AMSSystemType == sdk.AMSWechat {
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
	return false, fmt.Errorf("Invalid ReportAdLevel adLevel= %s", reportInput.AdLevel)
}

// getReportFilter 获取报表过滤字段
func (t *AMSReportService) getReportFilter(reportInput *sdk.GetReportInput) []model.FilteringStruct {
	if reportInput.Filtering == nil {
		return nil
	}
	TFiltering := make([]model.FilteringStruct, 0, TFilterMax)
	// campaign_id
	mFiltering := reportInput.Filtering
	if len(mFiltering.CampaignIDList) > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "campaign_id",
			Operator: "IN",
			Values:   &mFiltering.CampaignIDList,
		})
	}
	// adgroup_id
	if len(mFiltering.AdIDList) > 0 {
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
	return TFiltering
}

// getReportGroupBy 获取报表GroupBy字段
func (t *AMSReportService) getReportGroupBy(authAccount *sdk.AuthAccount, reportInput *sdk.GetReportInput) []string {
	var result []string
	// AMS 微信账户不支持
	if authAccount.AMSSystemType != sdk.AMS || reportInput.GroupBy == nil {
		return nil
	}
	groupby := *reportInput.GroupBy
	for i := 0; i < len(groupby); i++ {
		// 账户维度
		switch groupby[i] {
		case sdk.ADVERTISER_ID:
			// do nothing
			break
		case sdk.CAMPAIGN_ID:
			result = append(result, "campaign_id")
			break
		case sdk.AD_ID:
			result = append(result, "adgroup_id")
			break
		case sdk.CREATIVE_ID:
			result = append(result, "ad_id")
			break
		case sdk.Material_ID:
			result = append(result, "material_id")
			break
		default:
			// do nothing
		}

		// 时间维度
		switch groupby[i] {
		case sdk.DATE:
			result = append(result, "date")
			break
		case sdk.HOUR:
			result = append(result, "hour")
			break
		default:
			// do nothing
		}
	}
	return result
}

func (t *AMSReportService) getReportOrderType(sortType sdk.SortType) string {
	switch sortType {
	case sdk.ASC:
		return "ASCENDING"
	case sdk.DESC:
		return "DESCENDING"
	}
	return ""
}

func (t *AMSReportService) getReportField() []string {
	return []string{"date", "view_count", "valid_click_count", "ctr", "cpc", "cost", "account_id", "campaign_id",
		"campaign_name", "adgroup_id", "adgroup_name", "ad_id", "ad_name", "material_id"}
}

// getDailyReport 获取天级别的广告数据
func (t *AMSReportService) getDailyReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	id := formatAuthAccountID(reportInput.BaseInput.AccountId, reportInput.BaseInput.AMSSystemType)
	account, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("getDailyReport get AuthAccount Info error accid=%s", reportInput.BaseInput.AccountId)
	}

	tClient := getAMSSdkClient(account)
	var level string
	isSucc, err := t.getReportAdLevel(account, reportInput, &level)
	if !isSucc {
		return nil, err
	}
	dateRange := model.ReportDateRange{
		StartDate: reportInput.DateRange.StartDate,
		EndDate:   reportInput.DateRange.EndDate,
	}

	var dailyReportsGetOpts tapi.DailyReportsGetOpts
	if res := t.getReportFilter(reportInput); res != nil && len(res) > 0 {
		dailyReportsGetOpts.Filtering = optional.NewInterface(res)
	}

	if groupBy := t.getReportGroupBy(account, reportInput); groupBy != nil && len(groupBy) > 0 {
		dailyReportsGetOpts.GroupBy = optional.NewInterface(groupBy)
	}

	// OrderBy
	if len(reportInput.OrderBy.SortField) > 0 {
		dailyReportsGetOpts.OrderBy = optional.NewInterface([]model.
		OrderByStruct{{
			reportInput.OrderBy.SortField,
			model.Sortord(t.getReportOrderType(reportInput.OrderBy.SortType)),
		}})
	}

	// Page,Page_size
	if reportInput.Page > 0 {
		dailyReportsGetOpts.Page = optional.NewInt64(reportInput.Page)
	}
	if reportInput.PageSize > 0 {
		dailyReportsGetOpts.PageSize = optional.NewInt64(reportInput.PageSize)
	}

	// Fields
	dailyReportsGetOpts.Fields = optional.NewInterface(t.getReportField())

	accountid, err := strconv.ParseInt(reportInput.BaseInput.AccountId, 10, 64)
	if err != nil {
		return nil, err
	}
	// 获取天级别广告数据
	result, _, err := tClient.DailyReports().Get(*tClient.Ctx, accountid, level, dateRange,
		&dailyReportsGetOpts)
	if err != nil {
		return nil, err
	}
	reportOutput := &sdk.GetReportOutput{}
	t.copyDailyReportToOutput(&result, reportOutput)
	return reportOutput, err
}

// copyDailyReportToOutput 拷贝天级别报表数据
func (t *AMSReportService) copyDailyReportToOutput(dailyResponseData *model.DailyReportsGetResponseData, reportOutput *sdk.GetReportOutput) {
	if len(*dailyResponseData.List) == 0 {
		return
	}
	rList := make([]sdk.ReportOutputListStruct, 0, len(*dailyResponseData.List))
	for i := 0; i < len(*dailyResponseData.List); i++ {
		dailyData := (*dailyResponseData.List)[i]
		rList = append(rList, sdk.ReportOutputListStruct{
			AccountId:            dailyData.AccountId,
			CampaignId:           dailyData.CampaignId,
			CampaignName:         dailyData.CampaignName,
			AdgroupId:            dailyData.AdgroupId,
			AdgroupName:          dailyData.AdgroupName,
			AdId:                 dailyData.AdId,
			AdName:               dailyData.AdName,
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
			MaterialId:           dailyData.MaterialId,
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

	id := formatAuthAccountID(reportInput.BaseInput.AccountId, reportInput.BaseInput.AMSSystemType)
	account, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("getHourlyReport get AuthAccount Info error accid=%s", reportInput.BaseInput.AccountId)
	}

	tClient := getAMSSdkClient(account)
	var level string
	isSucc, err := t.getReportAdLevel(account, reportInput, &level)
	if !isSucc {
		return nil, err
	}
	dateRange := model.DateRange{
		StartDate: reportInput.DateRange.StartDate,
		EndDate:   reportInput.DateRange.EndDate,
	}

	var hourlyReportsGetOpts tapi.HourlyReportsGetOpts
	if res := t.getReportFilter(reportInput); res != nil && len(res) > 0 {
		hourlyReportsGetOpts.Filtering = optional.NewInterface(res)
	}

	if groupBy := t.getReportGroupBy(account, reportInput); groupBy != nil && len(groupBy) > 0 {
		hourlyReportsGetOpts.GroupBy = optional.NewInterface(groupBy)
	}

	// OrderBy
	if len(reportInput.OrderBy.SortField) > 0 {
		hourlyReportsGetOpts.OrderBy = optional.NewInterface([]model.
		OrderByStruct{{
			reportInput.OrderBy.SortField,
			model.Sortord(t.getReportOrderType(reportInput.OrderBy.SortType)),
		}})
	}

	// Page,Page_size
	if reportInput.Page > 0 {
		hourlyReportsGetOpts.Page = optional.NewInt64(reportInput.Page)
	}
	if reportInput.PageSize > 0 {
		hourlyReportsGetOpts.PageSize = optional.NewInt64(reportInput.PageSize)
	}

	// Fields
	hourlyReportsGetOpts.Fields = optional.NewInterface(t.getReportField())

	accountid, err := strconv.ParseInt(reportInput.BaseInput.AccountId, 10, 64)
	if err != nil {
		return nil, err
	}
	// 获取小时级别广告数据
	result, _, err := tClient.HourlyReports().Get(*tClient.Ctx, accountid, level, dateRange,
		&hourlyReportsGetOpts)
	if err != nil {
		return nil, err
	}
	reportOutput := &sdk.GetReportOutput{}
	t.copyHourReportToOutput(dateRange.StartDate, &result, reportOutput)
	return reportOutput, err
}

// copyHourReportToOutput 拷贝小时级别数据
func (t *AMSReportService) copyHourReportToOutput(date string, hourResponseData *model.HourlyReportsGetResponseData, reportOutput *sdk.GetReportOutput) {
	if len(*hourResponseData.List) == 0 {
		return
	}
	rList := make([]sdk.ReportOutputListStruct, 0, len(*hourResponseData.List))
	for i := 0; i < len(*hourResponseData.List); i++ {
		hourlyData := (*hourResponseData.List)[i]
		rList = append(rList, sdk.ReportOutputListStruct{
			AccountId:            hourlyData.AccountId,
			CampaignId:           hourlyData.CampaignId,
			CampaignName:         hourlyData.CampaignName,
			AdgroupId:            hourlyData.AdgroupId,
			AdgroupName:          hourlyData.AdgroupName,
			AdId:                 hourlyData.AdId,
			AdName:               hourlyData.AdName,
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

// GetVideoReport 获取视频报表
func (t *AMSReportService) GetVideoReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	id := formatAuthAccountID(reportInput.BaseInput.AccountId, reportInput.BaseInput.AMSSystemType)
	account, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("GetVideoReport get AuthAccount Info error accid=%s", reportInput.BaseInput.AccountId)
	}
	if account.AMSSystemType != sdk.AMS {
		return nil, fmt.Errorf("GetDailyVideoReport invalid account type = %d, id = %d",
			account.AMSSystemType, reportInput.BaseInput.AccountId)
	}
	reportInput.TimeGranularity = sdk.ReportTimeDaily
	reportInput.AdLevel = sdk.LevelVideo
	return t.getDailyReport(reportInput)
}

// GetImageReport 获取图片报表
func (t *AMSReportService) GetImageReport(reportInput *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	id := formatAuthAccountID(reportInput.BaseInput.AccountId, reportInput.BaseInput.AMSSystemType)
	account, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("GetImageReport get AuthAccount Info error accid=%s", reportInput.BaseInput.AccountId)
	}
	if account.AMSSystemType != sdk.AMS {
		return nil, fmt.Errorf("GetDailyVideoReport  invalid account type = %d, id = %d",
			account.AMSSystemType, reportInput.BaseInput.AccountId)
	}
	reportInput.TimeGranularity = sdk.ReportTimeDaily
	reportInput.AdLevel = sdk.LevelImage
	return t.getDailyReport(reportInput)
}
