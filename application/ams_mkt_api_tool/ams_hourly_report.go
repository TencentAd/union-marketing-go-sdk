package main

import (
	"context"
	"encoding/json"
	"git.code.oa.com/going/going/cat/qzs"
	"git.code.oa.com/going/mic"
	"git.code.oa.com/tme/kg_golang_proj/jce_go/proto_gdt_register"
	"git.code.oa.com/tme/kg_golang_proj/plib/code_utils"
	"github.com/tencentad/marketing-api-go-sdk/pkg/errors"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	mktSDKConfig "github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
)

type HourlyReportItem struct {
	TAds                 *ads.SDKClient
	AdCompanyType        string
	AccessToken          string
	AccountId            int64
	Level                string
	PageSize             int64 //单次拉取页面大小
	ReTrySec             int64 //失败重试时间间隔
	DateRange            model.DateRange
	HourlyReportsGetOpts *api.HourlyReportsGetOpts
	mapCampaignName      map[int64]string
}

func (e *HourlyReportItem) Init(conf *AccountConf) {
	e.AccessToken = conf.AccessToken
	e.TAds = ads.Init(&mktSDKConfig.SDKConfig{
		AccessToken: conf.AccessToken,
		IsDebug:     conf.IsDebug,
	})
	e.TAds.UseProduction()
	e.AccountId = conf.AccountId
	e.Level = conf.OnlineLevel
	e.AdCompanyType = conf.AdCompanyType
	e.PageSize = conf.OnlinePageSize
	e.ReTrySec = conf.ReTrySec
	e.DateRange = model.DateRange{
		StartDate: time.Now().Format("2006-01-02"),
		EndDate:   time.Now().Format("2006-01-02"),
	}
	e.HourlyReportsGetOpts = &api.HourlyReportsGetOpts{
		Page:     optional.NewInt64(1),
		PageSize: optional.NewInt64(e.PageSize),
		GroupBy:  optional.NewInterface(conf.OnlineGroupBy),
		OrderBy: optional.NewInterface([]model.OrderByStruct{
			{
				SortField: "hour",
				SortType:  "ASCENDING",
			},
		}),
		Fields: optional.NewInterface([]string{"account_id", "hour", "campaign_id", "adgroup_id",
			"ad_id", "view_count", "adgroup_name", "ad_name", "cost", "activated_count", "activated_cost",
			"retention_count", "retention_rate", "view_count", "thousand_display_price", "valid_click_count",
			"cpc", "ctr", "click_activated_rate"}),
	}
}

func (e *HourlyReportItem) setPage(page, pageSize int64) {
	e.HourlyReportsGetOpts.Page = optional.NewInt64(page)
	e.HourlyReportsGetOpts.PageSize = optional.NewInt64(pageSize)
}

func (e *HourlyReportItem) Run() error {
	qzsCtx := qzs.NewContext(context.Background())
	qzsCtx.Debug("handleData starts")

	tads := e.TAds
	tadCtx := *tads.Ctx
	e.setPage(1, e.PageSize)
	for true {
		_ = gHourlyLimit.Wait(context.Background()) //不可取消，没有超时
		//1 分页拉取 call api 拉取小时报
		start := time.Now()
		qzsCtx.Error("account id:%d pulling hourly data, page: %d, page size:%d",
			e.AccountId, e.HourlyReportsGetOpts.Page.Value(), e.HourlyReportsGetOpts.PageSize.Value())
		rsp, _, err := tads.HourlyReports().Get(tadCtx, e.AccountId, e.Level, e.DateRange, e.HourlyReportsGetOpts)
		cost := time.Since(start)
		if err != nil {
			errCode := -1
			if resErr, ok := err.(errors.ResponseError); ok {
				errStr, _ := json.Marshal(resErr)
				errCode = int(resErr.Code)
				qzsCtx.Error("pulling hourly data ResponseError:%s", string(errStr))
			} else {
				qzsCtx.Error("pulling hourly data Unknown error:%s", err)
			}
			mic.Report(proto_gdt_register.GDT_REGISTER_MOD_ID, proto_gdt_register.IF_GET_HOURLY_REPORT_AMS,
				"", errCode, 1, cost)
			qzsCtx.Error("account id:%d tads.HourlyReports().Get error, err:%+v, Item:%+v", e.AccountId, err, e)
			time.Sleep(time.Duration(e.ReTrySec) * time.Second)
			continue
		} else {
			mic.Report(proto_gdt_register.GDT_REGISTER_MOD_ID, proto_gdt_register.IF_GET_HOURLY_REPORT_AMS,
				"", 0, 0, cost)
		}
		qzsCtx.Error("accout id:%d get hourly data succ, page info:%+v.", e.AccountId, rsp.PageInfo)
		//2 DC上报数据
		listData := *rsp.List
		for _, singleRspData := range listData {
			mData := e.TransHourlyReports(singleRspData) //提取需要上报的数据
			for i := 0; i < 3; i++ {
				err = DoReport(qzsCtx, mData) //call api 上报
				if err != nil {
					qzsCtx.Error("DoReport error, err:%+v, mData:%+v\n", err, code_utils.GetPrettyJsonStr(mData))
				} else {
					qzsCtx.Debug("DoReport succ, err:%+v, mData:%+v\n", err, code_utils.GetPrettyJsonStr(mData))
					break
				}
			}
		}

		//3 检查是否还有下一页
		if hasNext, nextPage := HasNextPage(rsp.PageInfo); hasNext {
			e.setPage(nextPage, e.PageSize)
		} else {
			break
		}
	}
	return nil
}

func (e *HourlyReportItem) TransHourlyReports(reportData model.HourlyReportsGetListStruct) map[string]string {
	tNow := uint32(time.Now().Unix())
	mData := make(map[string]string)
	mData["key"] = "all_page#all_module#null#ad_realtime#0"
	mData["ad_company"] = e.AdCompanyType
	mData["advertiser_id"] = strconv.FormatInt(reportData.AccountId, 10)
	mData["report_time"] = strconv.FormatInt(time.Now().Unix(), 10)
	zeroTime := ZeroClockTime(tNow)
	mData["start_date"] = zeroTime.Add(time.Duration(reportData.Hour) * time.Hour).Format("2006-01-02 15:04:05")
	mData["end_date"] = zeroTime.Add(time.Duration(reportData.Hour+1) * time.Hour).Format("2006-01-02 15:04:05")
	mData["campaign_id"] = strconv.FormatInt(reportData.CampaignId, 10)
	mData["ad_id"] = strconv.FormatInt(reportData.AdgroupId, 10)
	mData["creative_id"] = strconv.FormatInt(reportData.AdId, 10)
	mData["campaign_name"] = gCampInfoCache.GetCampaignName(e.AccountId, reportData.CampaignId)
	mData["ad_name"] = reportData.AdgroupName

	// appid_渠道号_拉新/拉活_素材一级标签ID_一级标签_二级标签_三级标签_版位名称_出价模式_定向_上线日期_自定义
	if mData["campaign_name"] != "" {
		elements := strings.Split(mData["campaign_name"], "_")
		if len(elements) >= 2 {
			mData["channelid"] = elements[1]
		}
		if len(elements) >= 3 {
			if elements[2] == "拉新" {
				mData["account_type"] = "1"
			} else {
				mData["account_type"] = "2"
			}
		}
		if len(elements) >= 5 {
			mData["first_type"] = elements[4]
		}
		if len(elements) >= 6 {
			mData["second_type"] = elements[5]
		}
		if len(elements) >= 7 {
			mData["third_type"] = elements[6]
		}
	}

	// songmid/userkid_歌曲名/kol名_素材一级标签_二级标签_三级标签_素材中文名_代理_自定义
	if mData["ad_name"] != "" {
		elements := strings.Split(mData["ad_name"], "_")
		if len(elements) >= 7 {
			mData["agency"] = elements[6]
		}
		if len(elements) >= 6 {
			mData["material_name"] = elements[5]
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
				mData["songmid"] = elements[0]
			} else if elements[0] != "0" {
				mData["userkid"] = elements[0]
			}
		}
	}

	mData["cost"] = strconv.FormatFloat(float64(reportData.Cost)/100, 'f', 2, 64)
	mData["convert"] = strconv.FormatInt(reportData.ActivatedCount, 10)
	mData["convert_cost"] = strconv.FormatFloat(float64(reportData.ActivatedCost)/100, 'f', 2, 64)
	mData["deep_convert"] = strconv.FormatInt(reportData.RetentionCount, 10)
	mData["deep_convert_rate"] = strconv.FormatFloat(reportData.RetentionRate, 'f', 2, 64)
	mData["show"] = strconv.FormatInt(reportData.ViewCount, 10)
	mData["avg_show_cost"] = strconv.FormatFloat(float64(reportData.ThousandDisplayPrice)/100, 'f', 2, 64)
	mData["click"] = strconv.FormatInt(reportData.ValidClickCount, 10)
	mData["avg_click_cost"] = strconv.FormatFloat(float64(reportData.Cpc), 'f', 2, 64)
	mData["ctr"] = strconv.FormatFloat(reportData.Ctr, 'f', 2, 64)
	mData["convert_rate"] = strconv.FormatFloat(reportData.ClickActivatedRate, 'f', 2, 64)

	return mData
}
