package ocean_engine

import "github.com/tencentad/union-marketing-go-sdk/api/sdk"

type OceanEngineReportsData struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      *OceanEngineReportInfo `json:"data"`
	RequestId string      `json:"request_id"`
}

type OceanEngineReportInfo struct {
	AdList   []*OceanEngineReportAdInfo `json:"list"`
	PageInfo *sdk.PageConf   `json:"page_info"`
}

type OceanEngineReportAdInfo struct {
	AdvertiserId int64   `json:"advertiser_id"`
	CampaignId   int64   `json:"campaign_id"`
	CampaignName string  `json:"campaign_name"`
	AdId         int64   `json:"ad_id"`
	AdName       string  `json:"ad_name"`
	CreativeId   int64   `json:"creative_id"`
	Ctr          float64 `json:"ctr"`
	Cost         float64 `json:"cost"`
}
