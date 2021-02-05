package main

import (
	"context"
	"encoding/json"
	"git.code.oa.com/going/going/cat/qzs"
	"git.code.oa.com/going/mic"
	"git.code.oa.com/tme/kg_golang_proj/jce_go/proto_gdt_register"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	mktSDKConfig "github.com/tencentad/marketing-api-go-sdk/pkg/config"
	TadsErr "github.com/tencentad/marketing-api-go-sdk/pkg/errors"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type DailyReportItem struct {
	TAds                *ads.SDKClient
	AccessToken         string
	AccountId           int64
	AdCompanyType       string
	Level               string
	PageSize            int64 //单次拉取页面大小
	ReTrySec            int64 //失败重试时间间隔
	DateRange           model.ReportDateRange
	DailyReportsGetOpts *api.DailyReportsGetOpts
}

func (e *DailyReportItem) Init(conf *AccountConf) {
	e.AccessToken = conf.AccessToken
	e.TAds = ads.Init(&mktSDKConfig.SDKConfig{
		AccessToken: conf.AccessToken,
		IsDebug:     conf.IsDebug,
	})
	e.TAds.UseProduction()
	e.AccountId = conf.AccountId
	e.Level = conf.OfflineLevel
	e.AdCompanyType = conf.AdCompanyType
	e.PageSize = conf.OfflinePageSize
	e.ReTrySec = conf.ReTrySec
	e.DateRange = model.ReportDateRange{
		StartDate: time.Now().AddDate(0, 0, -conf.OffLineDays).Format("2006-01-02"),
		EndDate:   time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
	}
	e.DailyReportsGetOpts = &api.DailyReportsGetOpts{
		Page:     optional.NewInt64(1),
		PageSize: optional.NewInt64(e.PageSize),
		GroupBy:  optional.NewInterface(conf.OfflineGroupBy),
		OrderBy: optional.NewInterface([]model.OrderByStruct{
			{
				SortField: "date",
				SortType:  "ASCENDING",
			},
		}),
		Fields: optional.NewInterface([]string{"account_id", "date", "campaign_id", "adgroup_id",
			"ad_id", "view_count", "adgroup_name", "ad_name", "cost", "activated_count", "activated_cost",
			"retention_count", "retention_rate", "view_count", "thousand_display_price", "valid_click_count",
			"cpc", "ctr", "click_activated_rate"}),
	}
}

func (e *DailyReportItem) Run() error {
	qzsCtx := qzs.NewContext(context.Background())
	qzsCtx.Error("handle Daily report starts: account id:%d, start day:%s end day:%s.", e.AccountId, e.DateRange.StartDate, e.DateRange.EndDate)

	tads := e.TAds
	ctx := *tads.Ctx
	e.setPage(1, e.PageSize)
	for true {
		_ = gDailyLimit.Wait(context.Background()) //不可取消，没有超时
		//1 http拉取日报表
		start := time.Now()
		qzsCtx.Error("account id:%d pulling daily data, page: %d, page size:%d",
			e.AccountId, e.DailyReportsGetOpts.Page.Value(), e.DailyReportsGetOpts.PageSize.Value())
		dailyRsp, _, err := tads.DailyReports().Get(ctx, e.AccountId, e.Level, e.DateRange, e.DailyReportsGetOpts)
		cost := time.Since(start)
		if err != nil {
			errCode := -1
			if resErr, ok := err.(TadsErr.ResponseError); ok {
				errStr, _ := json.Marshal(resErr)
				errCode = int(resErr.Code)
				qzsCtx.Error("pulling daily data ResponseError:%s", string(errStr))
			} else {
				qzsCtx.Error("pulling daily data Unknown error:%+v", err)
			}
			mic.Report(proto_gdt_register.GDT_REGISTER_MOD_ID, proto_gdt_register.IF_GET_DAILY_REPORT_AMS,
				"", errCode, 1, cost)
			qzsCtx.Error("tads.DailyReports().Get error, err:%+v, Item:%+v", err, e)
			time.Sleep(time.Duration(e.ReTrySec) * time.Second)
			continue
		} else {
			mic.Report(proto_gdt_register.GDT_REGISTER_MOD_ID, proto_gdt_register.IF_GET_DAILY_REPORT_AMS,
				"", 0, 0, cost)
		}
		qzsCtx.Error("accout id:%d get daily data succ, page info:%+v.", e.AccountId, dailyRsp.PageInfo)

		//2 数据转换
		listData := *dailyRsp.List
		for _, singleRspData := range listData {
			mData := e.TransDailyReports(singleRspData) //提取需要上报的数据
			//3 上报罗盘
			for i := 0; i < 3; i++ { //3次重试
				err = DoReport(qzsCtx, mData) //call api 上报
				if err != nil {
					qzsCtx.Error("DoReport error, err:%+v, mData:%+v", err, mData)
				} else {
					qzsCtx.Debug("DoReport succ, mData:%+v.", mData)
					break
				}
			}
		}

		//4 检查是否还有下一页
		if hasNext, nextPage := HasNextPage(dailyRsp.PageInfo); hasNext {
			e.setPage(nextPage, e.PageSize)
		} else {
			break
		}
	}

	return nil
}

func (e *DailyReportItem) setPage(page, pageSize int64) {
	e.DailyReportsGetOpts.Page = optional.NewInt64(page)
	e.DailyReportsGetOpts.PageSize = optional.NewInt64(pageSize)
}

func (e *DailyReportItem) TransDailyReports(reportData model.DailyReportsGetListStruct) map[string]string {
	mData := make(map[string]string)
	mData["ad_company"] = e.AdCompanyType
	mData["key"] = "all_page#all_module#null#ad_offline#0"
	mData["start_date"] = reportData.Date + " 00:00:00"
	mData["end_date"] = reportData.Date + " 00:00:00"
	mData["advertiser_id"] = strconv.FormatInt(reportData.AccountId, 10)
	mData["report_time"] = strconv.FormatInt(time.Now().Unix(), 10)
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
