package main

import (
	"context"
    "flag"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"git.code.oa.com/going/going/cat/qzs"
	"git.code.oa.com/going/going/codec/ckv"
	"git.code.oa.com/going/going/codec/qzh"
	"git.code.oa.com/going/going/config"
	"git.code.oa.com/tme/kg_golang_proj/jce_go/app_dcreport"
	"git.code.oa.com/tme/kg_golang_proj/jce_go/proto_gdt_register"
)

type DimensionsData struct {
	Stat_datetime          string `json:"stat_datetime"`
	Advertiser_id          int    `json:"advertiser_id"`
	Campaign_name          string `json:"campaign_name"`
	Campaigin_id           int    `json:"campaigin_id"`
	Ad_name                string `json:"ad_name"`
	Ad_id                  int    `json:"ad_id"`
	Creative_id            int    `json:"creative_id"`
	Bidword                string `json:"bidword"`
	Bidword_id             int    `json:"bidword_id"`
	Query                  string `json:"query"`
	Pricing                string `json:"pricing"`
	Image_mode             string `json:"image_mode"`
	Inventory              string `json:"inventory"`
	Campaign_type          string `json:"campaign_type"`
	Creative_material_mode string `json:"creative_material_mode"`
	External_action        string `json:"external_action"`
	Landing_type           string `json:"landing_type"`
	Pricing_category       string `json:"pricing_category"`
	Province_name          string `json:"province_name"`
	City_name              string `json:"city_name"`
	Gender                 string `json:"gender"`
	Age                    string `json:"age"`
	Platform               string `json:"platform"`
	Ac                     string `json:"ac"`
	Ad_tag                 string `json:"ad_tag"`
	Interest_tag           string `json:"interest_tag"`
	Material_id            int    `json:"material_id"`
	Playable_id            int    `json:"playable_id"`
	Playable_name          string `json:"playable_name"`
	Playable_url           string `json:"playable_url"`
	Playable_orientation   string `json:"playable_orientation"`
	Playable_preview_url   string `json:"playable_preview_url"`
}

type MetricsData struct {
	Cost                                            float64 `json:"cost"`
	Show                                            int     `json:"show"`
	Avg_show_cost                                   float64 `json:"avg_show_cost"`
	Click                                           int     `json:"click"`
	Avg_click_cost                                  float64 `json:"avg_click_cost"`
	Ctr                                             float64 `json:"ctr"`
	Convert                                         int     `json:"convert"`
	Convert_cost                                    float64 `json:"convert_cost"`
	Convert_rate                                    float64 `json:"convert_rate"`
	Deep_convert                                    int     `json:"deep_convert"`
	Deep_convert_cost                               float64 `json:"deep_convert_cost"`
	Deep_convert_rate                               float64 `json:"deep_convert_rate"`
	Attribution_convert                             int     `json:"attribution_convert"`
	Attribution_convert_cost                        float64 `json:"attribution_convert_cost"`
	Attribution_deep_convert                        int     `json:"attribution_deep_convert"`
	Attribution_deep_convert_cost                   float64 `json:"attribution_deep_convert_cost"`
	Download_start                                  int     `json:"download_start"`
	Download_start_cost                             float64 `json:"download_start_cost"`
	Download_start_rate                             float64 `json:"download_start_rate"`
	Download_finish                                 int     `json:"download_finish"`
	Download_finish_cost                            float64 `json:"download_finish_cost"`
	Download_finish_rate                            float64 `json:"download_finish_rate"`
	Click_install                                   int     `json:"click_install"`
	Install_finish                                  int     `json:"install_finish"`
	Install_finish_cost                             float64 `json:"install_finish_cost"`
	Install_finish_rate                             float64 `json:"install_finish_rate"`
	Active                                          int     `json:"active"`
	Active_cost                                     float64 `json:"active_cost"`
	Active_rate                                     float64 `json:"active_rate"`
	Register                                        int     `json:"register"`
	Active_register_cost                            float64 `json:"active_register_cost"`
	Active_register_rate                            float64 `json:"active_register_rate"`
	Active_pay_amount                               float64 `json:"active_pay_amount"`
	Next_day_open                                   int     `json:"next_day_open"`
	Next_day_open_cost                              float64 `json:"next_day_open_cost"`
	Next_day_open_rate                              float64 `json:"next_day_open_rate"`
	Attribution_next_day_open_cnt                   int     `json:"attribution_next_day_open_cnt"`
	Attribution_next_day_open_cost                  float64 `json:"attribution_next_day_open_cost"`
	Attribution_next_day_open_rate                  float64 `json:"attribution_next_day_open_rate"`
	Game_addiction                                  int     `json:"game_addiction"`
	Game_addiction_cost                             float64 `json:"game_addiction_cost"`
	Game_addiction_rate                             float64 `json:"game_addiction_rate"`
	Pay_count                                       int     `json:"pay_count"`
	Active_pay_cost                                 float64 `json:"active_pay_cost"`
	Active_pay_rate                                 float64 `json:"active_pay_rate"`
	Loan_completion                                 int     `json:"loan_completion"`
	Loan_completion_cost                            float64 `json:"loan_completion_cost"`
	Loan_completion_rate                            float64 `json:"loan_completion_rate"`
	Pre_loan_credit                                 int     `json:"pre_loan_credit"`
	Pre_loan_credit_cost                            float64 `json:"pre_loan_credit_cost"`
	Loan_credit                                     int     `json:"loan_credit"`
	Loan_credit_cost                                float64 `json:"loan_credit_cost"`
	Loan_credit_rate                                float64 `json:"loan_credit_rate"`
	In_app_uv                                       int     `json:"in_app_uv"`
	In_app_detail_uv                                int     `json:"in_app_detail_uv"`
	In_app_cart                                     int     `json:"in_app_cart"`
	In_app_pay                                      int     `json:"in_app_pay"`
	In_app_order                                    int     `json:"in_app_order"`
	Attribution_game_pay_7d_count                   int     `json:"attribution_game_pay_7d_count"`
	Attribution_game_pay_7d_cost                    float64 `json:"attribution_game_pay_7d_cost"`
	Attribution_active_pay_7d_per_count             int     `json:"attribution_active_pay_7d_per_count"`
	Game_pay_cost                                   int     `json:"game_pay_cost"`
	Game_pay_count                                  int     `json:"game_pay_count"`
	Phone                                           int     `json:"phone"`
	Form                                            int     `json:"form"`
	Map_search                                      int     `json:"map_search"`
	Button                                          int     `json:"button"`
	View                                            int     `json:"view"`
	Download                                        int     `json:"download"`
	QQ                                              int     `json:"qq"`
	Lottery                                         int     `json:"lottery"`
	Vote                                            int     `json:"vote"`
	Message                                         int     `json:"message"`
	Redirect                                        int     `json:"redirect"`
	Shopping                                        int     `json:"shopping"`
	Consult                                         int     `json:"consult"`
	Wechat                                          int     `json:"wechat"`
	Phone_confirm                                   int     `json:"phone_confirm"`
	Phone_connect                                   int     `json:"phone_connect"`
	Consult_effective                               int     `json:"consult_effective"`
	Coupon                                          int     `json:"coupon"`
	Coupon_single_page                              int     `json:"coupon_single_page"`
	Redirect_to_shop                                int     `json:"redirect_to_shop"`
	Poi_collect                                     int     `json:"poi_collect"`
	Poi_address_click                               int     `json:"poi_address_click"`
	Luban_order_cnt                                 int     `json:"luban_order_cnt"`
	Luban_order_stat_amount                         float64 `json:"luban_order_stat_amount"`
	Luban_order_roi                                 float64 `json:"luban_order_roi"`
	Luban_live_enter_cnt                            int     `json:"luban_live_enter_cnt"`
	Live_watch_one_minute_count                     int     `json:"live_watch_one_minute_count"`
	Luban_live_follow_cnt                           int     `json:"luban_live_follow_cnt"`
	Live_fans_club_join_cnt                         int     `json:"live_fans_club_join_cnt"`
	Luban_live_comment_cnt                          int     `json:"luban_live_comment_cnt"`
	Luban_live_share_cnt                            int     `json:"luban_live_share_cnt"`
	Luban_live_gift_cnt                             int     `json:"luban_live_gift_cnt"`
	Luban_live_gift_amount                          float64 `json:"luban_live_gift_amount"`
	Luban_live_slidecart_click_cnt                  int     `json:"luban_live_slidecart_click_cnt"`
	Luban_live_click_product_cnt                    int     `json;"luban_live_click_product_cnt"`
	Luban_live_pay_order_count                      int     `json:"luban_live_pay_order_count"`
	Luban_live_pay_order_stat_cost                  float64 `json:"luban_live_pay_order_stat_cost"`
	Luban_live_pay_order_count_by_author_3days      float64 `json:"luban_live_pay_order_count_by_author_3days"`
	Luban_live_pay_order_stat_cost_by_author_3days  float64 `json:"luban_live_pay_order_stat_cost_by_author_3days"`
	Luban_live_pay_order_count_by_author_7days      float64 `json:"luban_live_pay_order_count_by_author_7days"`
	Luban_live_pay_order_stat_cost_by_author_7days  float64 `json:"luban_live_pay_order_stat_cost_by_author_7days"`
	Luban_live_pay_order_count_by_author_15days     float64 `json:"luban_live_pay_order_count_by_author_15days"`
	Luban_live_pay_order_stat_cost_by_author_15days float64 `json:"luban_live_pay_order_stat_cost_by_author_15days"`
	Luban_live_pay_order_count_by_author_30days     float64 `json:"luban_live_pay_order_count_by_author_30days"`
	Luban_live_pay_order_stat_cost_by_author_30days float64 `json:"luban_live_pay_order_stat_cost_by_author_30days"`
	Wechat_login_count                              int     `json:"wechat_login_count"`
	Attribution_wechat_login_30d_count              int     `json:"attribution_wechat_login_30d_count"`
	Wechat_login_cost                               float64 `json:"wechat_login_cost"`
	Attribution_wechat_login_30d_cost               float64 `json:"attribution_wechat_login_30d_cost"`
	Wechat_first_pay_count                          int     `json:"wechat_first_pay_count"`
	Attribution_wechat_first_pay_30d_count          int     `json:"attribution_wechat_first_pay_30d_count"`
	Wechat_first_pay_cost                           float64 `json:"wechat_first_pay_cost"`
	Attribution_wechat_first_pay_30d_cost           float64 `json:"attribution_wechat_first_pay_30d_cost"`
	Wechat_first_pay_rate                           float64 `json:"wechat_first_pay_rate"`
	Attribution_wechat_first_pay_30d_rate           float64 `json:"attribution_wechat_first_pay_30d_rate"`
	Wechat_pay_amount                               float64 `json:"wechat_pay_amount"`
	Attribution_wechat_pay_30d_amount               float64 `json:"attribution_wechat_pay_30d_amount"`
	Attribution_wechat_pay_30d_roi                  float64 `json:"attribution_wechat_pay_30d_roi"`
	Phone_effective                                 int     `json:"phone_effective"`
	Total_play                                      int     `json:"total_play"`
	Valid_play                                      int     `json:"valid_play"`
	Valid_play_cost                                 float64 `json:"valid_play_cost"`
	Valid_play_rate                                 float64 `json:"valid_play_rate"`
	Play_25_feed_break                              int     `json:"play_25_feed_break"`
	Play_50_feed_break                              int     `json:"play_50_feed_break"`
	Play_75_feed_break                              int     `json:"play_75_feed_break"`
	Play_100_feed_break                             int     `json:"play_100_feed_break"`
	Average_play_time_per_play                      float64 `json:"average_play_time_per_play"`
	Play_over_rate                                  float64 `json:"play_over_rate"`
	Wifi_play_rate                                  float64 `json:"wifi_play_rate"`
	Wifi_play                                       int     `json:"wifi_play"`
	Play_duration_sum                               int     `json:"play_duration_sum"`
	Advanced_creative_phone_click                   int     `json:"advanced_creative_phone_click"`
	Advanced_creative_counsel_click                 int     `json:"advanced_creative_counsel_click"`
	Advanced_creative_form_click                    int     `json:"advanced_creative_form_click"`
	Advanced_creative_coupon_addition               int     `json:"advanced_creative_coupon_addition"`
	Advanced_creative_form_submit                   int     `json:"advanced_creative_form_submit"`
	Card_show                                       int     `json:"card_show"`
	Share                                           int     `json:"share"`
	Comment                                         int     `json:"comment"`
	Like                                            int     `json:"like"`
	Follow                                          int     `json:"follow"`
	Home_visited                                    int     `json:"home_visited"`
	Ies_challenge_click                             int     `json:"ies_challenge_click"`
	Ies_music_click                                 int     `json:"ies_music_click"`
	Location_click                                  int     `json:"location_click"`
	Message_action                                  int     `json:"message_action"`
	Click_landing_page                              int     `json:"click_landing_page"`
	Click_shopwindow                                int     `json:"click_shopwindow"`
	Click_website                                   int     `json:"click_website"`
	Click_download                                  int     `json:"click_download"`
	Click_call_dy                                   int     `json:"click_call_dy"`
	Cpc                                             float64 `json:"cpc"`
	Cpm                                             float64 `json:"cpm"`
	Cpa                                             float64 `json:"cpa"`
	Interact_per_cost                               float64 `json:"interact_per_cost"`
	Convert_show_rate                               float64 `json:"convert_show_rate"`
	Report                                          int     `json:"report"`
	Dislike                                         int     `json:"dislike"`
	Play_duration_10s_rate                          float64 `json:"play_duration_10s_rate"`
	Play_50_feed_break_rate                         float64 `json:"play_50_feed_break_rate"`
	Play_25_feed_break_rate                         float64 `json:"play_25_feed_break_rate"`
	Play_duration_5s_rate                           float64 `json:"play_duration_5s_rate"`
	Play_duration_2s_rate                           float64 `json:"play_duration_2s_rate"`
	Average_play_progress                           float64 `json:"average_play_progress"`
	Play_100_feed_break_rate                        float64 `json:"play_100_feed_break_rate"`
	Play_duration_3s_rate                           float64 `json:"play_duration_3s_rate"`
	Play_75_feed_break_rate                         float64 `json:"play_75_feed_break_rate"`
	Play_duration                                   float64 `json:"play_duration"`
	Play_duration_3s                                float64 `json:"play_duration_3s"`
	Play_duration_2s                                float64 `json:"play_duration_2s"`
	Play_duration_10s                               float64 `json:"play_duration_10s"`
	Average_video_play                              float64 `json:"average_video_play"`
	Avg_rank                                        float64 `json:"avg_rank"`
}

type ListData struct {
	Dimensions DimensionsData `json:"dimensions"`
	Metrics    MetricsData    `json:"metrics"`
}

type PageData struct {
	Page         int `json:"page"`
	Page_size    int `json:"page_info"`
	Total_number int `json:"total_number"`
	Total_page   int `json:"total_page"`
}

type RspData struct {
	List      []ListData `json:"list"`
	Page_info PageData   `json:"page_info"`
}

type Rsp struct {
	Code       int     `json:"code"`
	Message    string  `json:"message"`
	Data       RspData `json:"data"`
	Request_id string  `json:"request_id"`
}

var conf = struct {
	MarketingAPI struct {
		DebugMode                string
		AdvertiserIdList         []string
        SleepInterval            config.Duration
    }
}{}

var (
    bTime = false
    DataTime int                    // 0表示今天，-1昨天，-2前天，-30一个月以前
    GroupByCondition  []string      // group_by condition
    bSleep bool                     // 是否需要sleep，实时数据不能sleep
)

func init() {
    flag.IntVar(&DataTime, "datetime", 0, "which day's data to be fetched")
    flag.BoolVar(&bSleep, "sleep", true, "sleep or not")
    flag.Var(newSliceValue([]string{}, &GroupByCondition), "condition", "GROUP BY CONDITION LIST SEPERATED BY COMMA")
}

func main() {
    flag.Parse()
	config.ConfPath = "../conf/marketing_api_tool.toml"
	config.Parse(&conf)
	handleData()
}

func handleData() {
	qzsCtx := qzs.NewContext(context.Background())
	qzsCtx.Debug("handleData starts")
    defer qzsCtx.WriteLog()

	// [step 1] 读取ckv中的token信息
	c := getCkvConn(qzsCtx)
	var stTokenInfo proto_gdt_register.TokenInfo
	key := proto_gdt_register.MARKETING_API_TOKEN_INFO_KEY
	err := c.DoGetJce(qzsCtx, key, &stTokenInfo)
	if err != nil {
		qzsCtx.Error("ckv get key:%v err:%+v", key, err)
		return
	}
	qzsCtx.Debug("ckv get key:%v succ, stTokenInfo:%+v", key, stTokenInfo)

	// [step 2] 并发拉取多个广告主的数据报表信息
	var wg sync.WaitGroup
	for _, AdvertiserId := range conf.MarketingAPI.AdvertiserIdList {
		wg.Add(2)
		go handleDataRealtime(qzsCtx, &wg, AdvertiserId, stTokenInfo)
		go handleDataOffline(qzsCtx, &wg, AdvertiserId, stTokenInfo)
        // sleep一段时间再拉下一个广告主的数据
        if (bSleep) {
            time.Sleep(conf.MarketingAPI.SleepInterval.Duration())
        }
	}
	wg.Wait()
}

func handleDataOffline(qzsCtx *qzs.Context, wg *sync.WaitGroup, AdvertiserId string, stTokenInfo proto_gdt_register.TokenInfo) {
	// [step 1] 先拉取和上报一次数据，获取 total_page
	var firstRsp Rsp
	FetchAndReportData(qzsCtx, AdvertiserId, stTokenInfo, "1", &firstRsp)
    qzsCtx.Debug("FetchAndReportData offline starts, data time:%+v, AdvertiserId:%+v, total_page:%v",
            DataTime, AdvertiserId, firstRsp.Data.Page_info.Total_page)

	// [step 2] 分页拉取
	for i := 2; i <= firstRsp.Data.Page_info.Total_page; i++ {
		var rsp Rsp
		FetchAndReportData(qzsCtx, AdvertiserId, stTokenInfo, strconv.Itoa(i), &rsp)
	}
	qzsCtx.Debug("FetchAndReportData offline ends, AdvertiserId:%+v, total_page:%v", AdvertiserId, firstRsp.Data.Page_info.Total_page)

	wg.Done()
}

func handleDataRealtime(qzsCtx *qzs.Context, wg *sync.WaitGroup, AdvertiserId string, stTokenInfo proto_gdt_register.TokenInfo) {
	// [step 1] 先拉取和上报一次数据，获取 total_page
	var firstRsp Rsp
	FetchAndReportData(qzsCtx, AdvertiserId, stTokenInfo, "1", &firstRsp)
	qzsCtx.Debug("FetchAndReportData realtime starts, AdvertiserId:%+v, total_page:%v", AdvertiserId, firstRsp.Data.Page_info.Total_page)

	// [step 2] 分页拉取
	for i := 2; i <= firstRsp.Data.Page_info.Total_page; i++ {
		var rsp Rsp
		FetchAndReportData(qzsCtx, AdvertiserId, stTokenInfo, strconv.Itoa(i), &rsp)
	}
	qzsCtx.Debug("FetchAndReportData realtime ends, AdvertiserId:%+v, total_page:%v", AdvertiserId, firstRsp.Data.Page_info.Total_page)

	wg.Done()
}

func FetchAndReportData(qzsCtx *qzs.Context, AdvertiserId string, stTokenInfo proto_gdt_register.TokenInfo, currPage string, pRsp *Rsp) {
	// [step 1] 填充数据
	req, err := http.NewRequest("GET", "https://ad.oceanengine.com/open_api/2/report/integrated/get/", nil)
	if err != nil {
		qzsCtx.Error("http.NewRequest err:%+v", err)
		return
	}

	params := req.URL.Query()
	params.Add("advertiser_id", AdvertiserId)

    var curTime = time.Now()
	params.Add("start_date", curTime.AddDate(0, 0, DataTime).Format("2006-01-02"))
	params.Add("end_date", curTime.AddDate(0, 0, DataTime).Format("2006-01-02"))
	params.Add("page", currPage)

	var SingleCondition string
	TotalConditions := "["

	for _, val := range GroupByCondition {
		SingleCondition = "\"" + val + "\"" + ","
		TotalConditions += SingleCondition
		qzsCtx.Debug("SingleCondition:%v TotalConditions:%v", SingleCondition, TotalConditions)
	}
	if len(TotalConditions) == 1 {
		qzsCtx.Error("group_by param is required, invalid request")
		return
	}

	ResCondition := strings.TrimRight(TotalConditions, ",")
	ResCondition += "]"
	params.Add("group_by", ResCondition)
	qzsCtx.Debug("TotalConditions:%v, ResCondition:%+v", TotalConditions, ResCondition)

	// condition不能urlencode
	RawQuery := params.Encode()
	req.URL.RawQuery, _ = url.QueryUnescape(RawQuery)
	// req.URL.RawQuery = params.Encode()

	req.Header.Set("Access-Token", stTokenInfo.StrAccessToken)
	req.Header.Set("X-Debug-Mode", conf.MarketingAPI.DebugMode)

	// [step 2] 拉取数据
	client := http.Client{}
	resp, err := client.Do(req) //Do 方法发送请求，返回 HTTP 回复
	if err != nil {
		qzsCtx.Error("Client.Do err:%+v, req:%+v", err, req)
		return
	}
	defer resp.Body.Close()
	qzsCtx.Debug("Client.Do succ, req:%+v", req)

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		qzsCtx.Error("ioutil.ReadAll err:%+v, token:%+v", err, stTokenInfo)
		return
	}
	err = json.Unmarshal(respBytes, pRsp)
	if err != nil {
		qzsCtx.Error("json.Unmarshal err:%+v, TokenInfo:%+v", err, stTokenInfo)
		return
	}
	qzsCtx.Debug("json.Unmarshal succ, TokenInfo:%+v, page:%+v, rsp:%+v", stTokenInfo, currPage, *pRsp)

	qzsCtx.Error("code: %+v, msg: %+v, reqID: %+v", pRsp.Code, pRsp.Message, pRsp.Request_id)

	//  [step 3] 数据上报
	DoReport(qzsCtx, *pRsp)
}

// 创建ckv连接
func getCkvConn(ctx *qzs.Context) *ckv.Ckv {
	conn := ckv.New("token_info")
	ctx.Trace("conn:%s", conn.String())
	return conn
}

func DoReport(ctx *qzs.Context, rsp Rsp) {
	qcReq := &app_dcreport.DataReportReq{}
	qcRsp := &app_dcreport.DataReportRsp{}
	item := app_dcreport.DataReportItem{}

	// 填充上报信息
	for _, val := range rsp.Data.List {
		item.MData = make(map[string]string)
		if DataTime == 0 {
			item.MData["key"] = "all_page#all_module#null#ad_realtime#0"
		} else {
			item.MData["key"] = "all_page#all_module#null#ad_offline#0"
        }

		report_time := time.Now().Unix()
		if bTime == false {
			bTime = true
			ctx.Error("report time: %+v", report_time)
		}

        item.MData["start_date"] = val.Dimensions.Stat_datetime
		item.MData["end_date"] = val.Dimensions.Stat_datetime

		item.MData["ad_company"] = "1"
		item.MData["advertiser_id"] = strconv.Itoa(val.Dimensions.Advertiser_id)
		item.MData["report_time"] = strconv.FormatInt(report_time, 10)
		item.MData["campaign_id"] = strconv.Itoa(val.Dimensions.Campaigin_id)
		item.MData["ad_id"] = strconv.Itoa(val.Dimensions.Ad_id)
		item.MData["creative_id"] = strconv.Itoa(val.Dimensions.Creative_id)
		item.MData["campaign_name"] = val.Dimensions.Campaign_name
		item.MData["ad_name"] = val.Dimensions.Ad_name
        item.MData["inventory_type"] = val.Dimensions.Inventory

		// appid_渠道号_拉新/拉活_素材一级标签ID_一级标签_二级标签_三级标签_版位名称_出价模式_定向_上线日期_自定义
		if item.MData["campaign_name"] != "" {
			elements := strings.Split(item.MData["campaign_name"], "_")
			if len(elements) >= 2 {
				item.MData["channelid"] = elements[1]
			}
			if len(elements) >= 3 {
				if elements[2] == "拉新" {
					item.MData["account_type"] = "1"
				} else {
					item.MData["account_type"] = "2"
				}
			}
			if len(elements) >= 5 {
				item.MData["first_type"] = elements[4]
			}
			if len(elements) >= 6 {
				item.MData["second_type"] = elements[5]
			}
			if len(elements) >= 7 {
				item.MData["third_type"] = elements[6]
			}
		}

		// songmid/userkid_歌曲名/kol名_素材一级标签_二级标签_三级标签_素材中文名_代理_自定义
		if item.MData["ad_name"] != "" {
			elements := strings.Split(item.MData["ad_name"], "_")
			if len(elements) >= 7 {
				item.MData["agency"] = elements[6]
			}
			if len(elements) >= 6 {
				item.MData["material_name"] = elements[5]
			}
			if len(elements) >= 1 {
				materialInfo := elements[0]
				bMid := false
				for _, char := range materialInfo {
					// 出现了非数字的字符，则为mid；否则为uid
					if !unicode.IsDigit(char) {
						bMid = true
						break
					}
				}
				if bMid {
					item.MData["songmid"] = elements[0]
				} else if elements[0] != "0" {
					item.MData["userkid"] = elements[0]
				}
			}
		}

		item.MData["cost"] = strconv.FormatFloat(val.Metrics.Cost, 'f', 2, 64)
		item.MData["convert"] = strconv.Itoa(val.Metrics.Convert)
		item.MData["convert_cost"] = strconv.FormatFloat(val.Metrics.Convert_cost, 'f', 2, 64)
		item.MData["deep_convert"] = strconv.Itoa(val.Metrics.Deep_convert)
		item.MData["deep_convert_rate"] = strconv.FormatFloat(val.Metrics.Deep_convert_rate, 'f', 2, 64)
		item.MData["show"] = strconv.Itoa(val.Metrics.Show)
		item.MData["avg_show_cost"] = strconv.FormatFloat(val.Metrics.Avg_show_cost, 'f', 2, 64)
		item.MData["click"] = strconv.Itoa(val.Metrics.Click)
		item.MData["avg_click_cost"] = strconv.FormatFloat(val.Metrics.Avg_click_cost, 'f', 2, 64)
		item.MData["ctr"] = strconv.FormatFloat(val.Metrics.Ctr, 'f', 2, 64)
		item.MData["convert_rate"] = strconv.FormatFloat(val.Metrics.Convert_rate, 'f', 2, 64)

		qcReq.VecReportItem = append(qcReq.VecReportItem, item)
	}
	callDesc := qzh.CallDesc{
		CmdId:       app_dcreport.NUM_CMD_DCREPORT,
		SubCmdId:    app_dcreport.CMD_DCREPORT_NEW,
		AppProtocol: "qza",
	}
	var authInfo qzh.AuthInfo
	qc := qzh.NewQzClient(callDesc, authInfo, "dcreport", qcReq, qcRsp)
	err := qc.Do(ctx)
	if err != nil {
		err := err.(qzh.QzError)
		ctx.Error("dcreport err%+v", err)
		return
	}
	ctx.Debug("dcreport succ, req:%+v, rsp:%+v", qcReq, qcRsp)
}
