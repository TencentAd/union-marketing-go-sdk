package main

import (
	"fmt"
	"git.code.oa.com/going/going/config"
	"golang.org/x/time/rate"
	"sync"
)

const (
	ReportTypeHourly = 1 << 0
	ReportTypeDaily  = 1 << 1
)

var AmsMktConf struct {
	QPSGet struct{
		DailyData  float64 //日报数据的QPS
		HourlyData float64 //小时数据的QPS
	}
	MapAccount       map[string]AccountConf
}
var gCampInfoCache *CampaignInfoCache
var gHourlyLimit *rate.Limiter
var gDailyLimit *rate.Limiter

type AccountConf struct {
	AccountId       int64
	AccessToken     string
	AdCompanyType   string
	OnlineLevel     string
	OnlineGroupBy   []string
	OnlinePageSize  int64
	OfflineLevel    string
	OfflineGroupBy  []string
	OffLineDays     int
	OfflinePageSize int64
	ReTrySec        int64 //失败重试时间间隔
	ReportType      int
	IsDebug         bool
}

func HandleHourlyReport(accountConf *AccountConf) error {
	e := &HourlyReportItem{}
	e.Init(accountConf)
	return e.Run()
}

func HandleDailyReport(accountConf *AccountConf) error {
	e := &DailyReportItem{}
	e.Init(accountConf)
	return e.Run()
}

func main() {
	config.ConfPath = "../conf/ams_mkt_api_debug.conf"
	err := config.Parse(&AmsMktConf)
	if err != nil {
		fmt.Printf("parse conf error, err:%+v.\n", err)
	}

	//预读取广告计划信息
	gCampInfoCache = NewCampaignInfoCache()
	gCampInfoCache.BatchAddAccount(AmsMktConf.MapAccount)

	//限流,每秒有放入 HourlyData 个令牌，初始 HourlyData 个
	gHourlyLimit = rate.NewLimiter(rate.Limit(AmsMktConf.QPSGet.HourlyData), len(AmsMktConf.MapAccount))
	gDailyLimit = rate.NewLimiter(rate.Limit(AmsMktConf.QPSGet.DailyData), len(AmsMktConf.MapAccount))

	wg := sync.WaitGroup{}
	for _, conf := range AmsMktConf.MapAccount {
		if gCampInfoCache.GetCampaignSize(conf.AccountId) == 0 {
			fmt.Printf("account id:%d has no campaign, filtered.\n", conf.AccountId)
			continue
		}
		wg.Add(1)
		fmt.Printf("init account: %d\n", conf.AccountId)
		go func(tmpConf AccountConf) {
			//实时按小时报
			fmt.Printf("init hourly account: %d\n", tmpConf.AccountId)
			if (tmpConf.ReportType & ReportTypeHourly) != 0 {
				_ = HandleHourlyReport(&tmpConf)
			}
			wg.Done()
		}(conf)
		wg.Add(1)
		go func(tmpConf AccountConf) {
			//离线按天报
			fmt.Printf("init daily account: %d\n", tmpConf.AccountId)
			if (tmpConf.ReportType & ReportTypeDaily) != 0 {
				_ = HandleDailyReport(&tmpConf)
			}
			wg.Done()
		}(conf)
	}
	wg.Wait()

}
