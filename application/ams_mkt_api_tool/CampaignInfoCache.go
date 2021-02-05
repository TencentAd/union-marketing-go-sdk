package main

import (
	"fmt"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	mktSDKConfig "github.com/tencentad/marketing-api-go-sdk/pkg/config"
)

type CampaignInfoCache struct {
	mapCampName map[int64]map[int64]string
}

func NewCampaignInfoCache() *CampaignInfoCache {
	return &CampaignInfoCache{
		mapCampName: make(map[int64]map[int64]string),
	}
}

func (m *CampaignInfoCache) AddAccount(conf *AccountConf) {
	tds := ads.Init(&mktSDKConfig.SDKConfig{
		AccessToken: conf.AccessToken,
		IsDebug:     false,
	})
	tds.UseProduction()
	tmpMap, _ := GetAllCampaignInfo(tds, conf.AccountId)
	if tmpMap == nil {
		tmpMap = make(map[int64]string)
	}
	m.mapCampName[conf.AccountId] = tmpMap
	fmt.Printf("account id:%d campign size:%d\n", conf.AccountId, len(tmpMap))
}

func (m *CampaignInfoCache) BatchAddAccount(allConf map[string]AccountConf) {
	for _, tmpConf := range allConf {
		m.AddAccount(&tmpConf)
	}
}

func (m *CampaignInfoCache) GetCampaignName(accountId, campaignId int64) string {
	if tmpMap, ok := m.mapCampName[accountId]; ok {
		if campaignName, ok := tmpMap[campaignId]; ok {
			return campaignName
		}
	}
	return ""
}

func (m CampaignInfoCache) GetCampaignSize(accountId int64) int {
	if tmpMap, ok := m.mapCampName[accountId]; ok {
		return len(tmpMap)
	}
	return 0
}
