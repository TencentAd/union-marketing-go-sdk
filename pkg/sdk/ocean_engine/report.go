package ocean_engine

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/http_tools"
)

// ReportService
type ReportService struct {
	config     *config.Config
	httpClient *http_tools.HttpClient
}

// NewReportService 创建报表服务
func NewReportService(sConfig *config.Config) *ReportService {
	return &ReportService{
		config:     sConfig,
		httpClient: http_tools.Init(sConfig.HttpConfig),
	}
}

// getReportGroupBy 获取报表GroupBy字段
func (s *ReportService) getReportGroupBy(authAccount *sdk.AuthAccount, reportInput *sdk.GetReportInput) string {
	var groupByResult []string
	groupByInput := reportInput.GroupBy
	for i := 0; i < len(groupByInput); i++ {
		// 账户维度
		switch groupByInput[i] {
		case sdk.AdvertiserId:
			groupByResult = append(groupByResult, "STAT_GROUP_BY_ADVERTISER_ID")
			break
		case sdk.CampaignId:
			groupByResult = append(groupByResult, "STAT_GROUP_BY_CAMPAIGN_ID")
			break
		case sdk.AdId:
			groupByResult = append(groupByResult, "STAT_GROUP_BY_AD_ID")
			break
		case sdk.CreativeId:
			groupByResult = append(groupByResult, "STAT_GROUP_BY_CREATIVE_ID")
			break
		case sdk.MaterialId:
			groupByResult = append(groupByResult, "STAT_GROUP_BY_MATERIAL_ID")
			break
		default:
			// do nothing
		}

		// 时间维度
		switch groupByInput[i] {
		case sdk.Date:
			groupByResult = append(groupByResult, "STAT_GROUP_BY_TIME_DAY")
			break
		case sdk.Hour:
			groupByResult = append(groupByResult, "STAT_GROUP_BY_TIME_HOUR")
			break
		default:
			// do nothing
		}

		// 其他维度
		if groupByInput[i] == sdk.InventoryOceanEngine {
			groupByResult = append(groupByResult, "STAT_GROUP_BY_INVENTORY")
		}
	}
	result, _ := json.Marshal(groupByResult)
	return string(result)
}

func (s *ReportService) getReportTime(granularity sdk.TimeGranularity) string {
	switch granularity {
	case sdk.ReportTimeDaily:
		return "STAT_TIME_GRANULARITY_DAILY"
	case sdk.ReportTimeHour:
		return "STAT_TIME_GRANULARITY_HOURLY"
	}
	return ""
}

type FilterInfo struct {
	CampaignIDList []int64 `json:"campaign_ids,omitempty"` // 计划id列表
	AdIDList       []int64 `json:"ad_id,omitempty"`        // 广告id列表
	CreativeIDList []int64 `json:"creative_id,omitempty"`  // 广告创意id列表
}

// getReportFilter 获取报表过滤字段
func (s *ReportService) getReportFilter(input *sdk.GetReportInput) (*FilterInfo, error) {
	if input.Filtering == nil {
		return nil, nil
	}
	// campaign_id
	filterInfo := &FilterInfo{}
	if len(input.Filtering.CampaignIDList) > 0 {
		camIDList := input.Filtering.CampaignIDList
		for i := 0; i < len(camIDList); i++ {
			camID, err := strconv.ParseInt(camIDList[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("getReportFilter error accID = %s", camIDList[i])
			}
			filterInfo.CampaignIDList = append(filterInfo.CampaignIDList, camID)
		}
	}

	// ad_id
	if len(input.Filtering.AdIDList) > 0 {
		adIDList := input.Filtering.AdIDList
		for i := 0; i < len(adIDList); i++ {
			adID, err := strconv.ParseInt(adIDList[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("getReportFilter error adID = %s", adIDList[i])
			}
			filterInfo.CampaignIDList = append(filterInfo.CampaignIDList, adID)
		}
	}

	// creative_id
	if len(input.Filtering.CreativeIDList) > 0 {
		creativeIDList := input.Filtering.CreativeIDList
		for i := 0; i < len(creativeIDList); i++ {
			creativeID, err := strconv.ParseInt(creativeIDList[i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("getReportFilter error creativeID = %s", creativeIDList[i])
			}
			filterInfo.CampaignIDList = append(filterInfo.CampaignIDList, creativeID)
		}
	}

	return filterInfo, nil
}

func (s *ReportService) getReportField(input *sdk.GetReportInput) string {
	var fieldList []string
	switch input.AdLevel {
	case sdk.LevelAccount:
	case sdk.LevelCampaign:
	case sdk.LevelAd:
	case sdk.LevelCreative:
		fieldList = []string{"cost", "show", "avg_show_cost", "click", "avg_click_cost", "ctr", "convert", "convert_cost",
			"convert_rate", "deep_convert", "deep_convert_cost", "deep_convert_rate", "attribution_convert",
			"attribution_convert_cost", "attribution_deep_convert", "attribution_deep_convert_cost", "download_start",
			"download_start_cost", "download_start_rate", "download_finish", "download_finish_cost",
			"download_finish_rate", "install_finish", "install_finish_cost", "install_finish_rate", "click_install",
			"active", "active_cost", "active_rate", "register", "active_register_cost", "active_register_rate",
			"pay_count", "active_pay_cost", "active_pay_rate", "next_day_open", "next_day_open_cost",
			"next_day_open_rate", "attribution_next_day_open_cnt", "attribution_next_day_open_cost",
			"attribution_next_day_open_rate", "attribution_game_pay_7d_count", "attribution_game_pay_7d_cost",
			"attribution_active_pay_7d_per_count", "game_addiction", "game_addiction_cost", "game_addiction_rate",
			"game_pay_cost", "game_pay_count", "loan_completion", "loan_completion_cost", "loan_completion_rate",
			"pre_loan_credit", "pre_loan_credit_cost", "loan_credit", "loan_credit_cost", "loan_credit_rate", "in_app_uv",
			"in_app_detail_uv", "in_app_cart", "in_app_pay", "in_app_order", "phone", "form", "button", "map_search", "view",
			"download", "qq", "lottery", "vote", "message", "redirect", "shopping", "consult", "wechat", "phone_confirm",
			"phone_connect", "phone_effective", "consult_effective", "coupon", "coupon_single_page", "redirect_to_shop",
			"poi_collect", "poi_address_click", "luban_order_cnt", "luban_order_stat_amount", "luban_order_roi",
			"luban_live_enter_cnt", "live_watch_one_minute_count", "luban_live_follow_cnt", "live_fans_club_join_cnt",
			"luban_live_comment_cnt", "luban_live_share_cnt", "luban_live_gift_cnt", "luban_live_gift_amount",
			"luban_live_slidecart_click_cnt", "luban_live_click_product_cnt", "luban_live_pay_order_count",
			"luban_live_pay_order_stat_cost", "wechat_login_count", "attribution_wechat_login_30d_count",
			"wechat_login_cost", "attribution_wechat_login_30d_cost", "wechat_first_pay_count",
			"attribution_wechat_first_pay_30d_count", "wechat_first_pay_cost", "attribution_wechat_first_pay_30d_cost",
			"wechat_first_pay_rate", "attribution_wechat_first_pay_30d_rate", "wechat_pay_amount",
			"attribution_wechat_pay_30d_amount", "attribution_wechat_pay_30d_roi", "total_play", "valid_play",
			"valid_play_cost", "valid_play_rate", "play_25_feed_break", "play_50_feed_break", "play_75_feed_break",
			"play_100_feed_break", "average_play_time_per_play", "play_over_rate", "wifi_play", "play_duration_sum",
			"wifi_play_rate", "card_show", "advanced_creative_phone_click", "advanced_creative_counsel_click",
			"advanced_creative_form_click", "advanced_creative_coupon_addition", "advanced_creative_form_submit",
			"like", "share", "comment", "follow", "home_visited", "ies_challenge_click", "ies_music_click",
			"location_click", "message_action", "click_landing_page", "click_shopwindow", "click_website",
			"click_download", "click_call_dy", "submit_certification_count", "approval_count", "first_order_count",
			"first_rental_order_count", "commute_first_pay_count"}
		break
	case sdk.LevelVideo:
	case sdk.LevelImage:
		fieldList = []string{"cost", "show", "click", "convert", "download_start", "download_finish",
			"install_finish", "active", "register", "next_day_open", "game_addiction", "pay_count", "phone", "form",
			"map_search", "button", "view", "download", "qq", "lottery", "vote", "message", "redirect", "shopping",
			"consult", "wechat", "phone_confirm", "phone_connect", "consult_effective", "coupon",
			"coupon_single_page", "total_play", "valid_play", "play_25_feed_break", "play_50_feed_break",
			"play_75_feed_break", "play_100_feed_break", "advanced_creative_phone_click",
			"advanced_creative_counsel_click", "advanced_creative_form_click", "advanced_creative_coupon_addition",
			"share", "comment", "like", "follow", "home_visited", "ies_challenge_click", "ies_music_click",
			"location_click", "play_duration_sum", "active_pay_amount", "phone_effective", "wifi_play",
			"play_duration", "play_duration_3s", "in_app_uv", "in_app_detail_uv", "in_app_cart", "in_app_pay",
			"in_app_order", "play_duration_2s", "play_duration_10s", "active_cost", "active_pay_cost",
			"active_pay_rate", "active_rate", "active_register_cost", "active_register_rate",
			"average_play_time_per_play", "average_video_play", "download_start_cost", "download_start_rate",
			"convert_cost", "convert_rate", "cpa", "cpc", "cpm", "ctr", "download_finish_cost",
			"download_finish_rate", "game_addiction_cost", "game_addiction_rate", "install_finish_cost",
			"install_finish_rate", "next_day_open_cost", "next_day_open_rate", "play_over_rate", "wifi_play_rate",
			"valid_play_cost", "valid_play_rate", "convert_show_rate"}
		break
	default:
		fieldList = []string{"cost", "show", "avg_show_cost", "click", "avg_click_cost", "ctr", "convert",
			"convert_cost", "convert_rate", "deep_convert", "deep_convert_cost", "deep_convert_rate",
			"attribution_convert", "attribution_convert_cost", "attribution_deep_convert", "attribution_deep_convert_cost"}
		break
	}

	field, _ := json.Marshal(fieldList)
	return string(field)
}

func (s *ReportService) GetReport(input *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("getAccountReport get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	method := http_tools.GET
	// create path and map variables
	path := s.config.HttpConfig.BasePath + "/2/report/integrated/get/"

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Accept"] = "application/json"
	header["Access-Token"] = authAccount.AccessToken
	header["X-Debug-Mode"] = "1"

	if len(input.BaseInput.AccountId) == 0 || len(input.DateRange.StartDate) == 0 || len(input.DateRange.
		EndDate) == 0 {
		return nil, fmt.Errorf("some field is empth accID = %s, startDate = %s, endDate = %s",
			input.BaseInput.AccountId, input.DateRange.StartDate, input.DateRange.EndDate)
	}
	query := url.Values{}
	query["advertiser_id"] = []string{input.BaseInput.AccountId}
	query["start_date"] = []string{input.DateRange.StartDate}
	query["end_date"] = []string{input.DateRange.EndDate}
	query["group_by"] = []string{s.getReportGroupBy(authAccount, input)}
	// fields
	query["fields"] = []string{s.getReportField(input)}
	// filtering
	filterInfo, err := s.getReportFilter(input)
	if err != nil {
		return nil, err
	}
	if filterInfo != nil {
		filtering, _ := json.Marshal(filterInfo)
		query["filtering"] = []string{string(filtering)}
	}

	// OrderBy
	if len(input.OrderBy.SortField) > 0 {
		query["order_field"] = []string{input.OrderBy.SortField}
		query["order_type"] = []string{string(input.OrderBy.SortType)}
	}

	if input.Page > 0 {
		query["page"] = []string{strconv.FormatInt(input.Page, 10)}
	}
	if input.PageSize > 0 {
		query["page_size"] = []string{strconv.FormatInt(input.PageSize, 10)}
	}

	request, err := s.httpClient.PrepareRequest(context.Background(), path, method, nil, header,
		query, nil, "", nil, "")
	if err != nil {
		return nil, err
	}

	reportsData := &OceanEngineReportsData{}
	respErr := s.httpClient.DoProcess(request, reportsData)
	if respErr != nil {
		return nil, respErr
	}
	if reportsData.Code != 0 {
		return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", reportsData.Code,
			reportsData.Message,
			reportsData.RequestId)
	}
	output := &sdk.GetReportOutput{}
	s.copyReportToOutput(reportsData, output)
	return output, nil
}

// copyAccountReportToOutput 拷贝账户报表
func (s *ReportService) copyReportToOutput(input *OceanEngineReportsData,
	reportOutput *sdk.GetReportOutput) {
	if len(input.Data.AdList) == 0 {
		return
	}
	rList := make([]*sdk.ReportOutputListStruct, 0, len(input.Data.AdList))
	for i := 0; i < len(input.Data.AdList); i++ {
		adInfo := (input.Data.AdList)[i]
		rList = append(rList, &sdk.ReportOutputListStruct{
			AccountId:    adInfo.AdvertiserId,
			CampaignId:   adInfo.CampaignId,
			CampaignName: adInfo.CampaignName,
			AdId:         adInfo.AdId,
			AdName:       adInfo.AdName,
			Ctr:          adInfo.Ctr,
		})
	}
	reportOutput.List = rList
	reportOutput.PageInfo = &sdk.PageConf{
		Page:        input.Data.PageInfo.Page,
		PageSize:    input.Data.PageInfo.PageSize,
		TotalNumber: input.Data.PageInfo.TotalNumber,
		TotalPage:   input.Data.PageInfo.TotalPage,
	}
}

// GetVideoReport 获取视频报表
func (s *ReportService) GetVideoReport(input *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	return s.GetReport(input)
}

// GetImageReport 获取图片报表
func (s *ReportService) GetImageReport(input *sdk.GetReportInput) (*sdk.GetReportOutput, error) {
	return s.GetReport(input)
}
