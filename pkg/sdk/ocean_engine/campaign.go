package ocean_engine

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/http_tools"
)

type CampaignService struct {
	config     *config.Config
	httpClient *http_tools.HttpClient
}

// NewCampaignService 获取广告组服务
func NewCampaignService(sConfig *config.Config) *CampaignService {
	return &CampaignService{
		config:     sConfig,
		httpClient: http_tools.Init(sConfig.HttpConfig),
	}
}

type CampaignGetFilterInfo struct {
	CampaignIds        []int64 `json:"ids,omitempty"`                  // 广告组ID过滤，数组，不超过100个
	CampaignName       string  `json:"campaign_name,omitempty"`        // 广告组name过滤，长度为1-30个字符，其中1个中文字符算2位
	LandingType        string  `json:"landing_type,omitempty"`         // 广告组推广目的过滤
	CampaignCreateTime string  `json:"campaign_create_time,omitempty"` // 广告组创建时间，格式yyyy-mm-dd，表示过滤出当天创建的广告组
}

//// getFilter 获取过滤信息
//func (s *CampaignService) getFilter(input *sdk.CampaignGetInput) (string, error) {
//	if input.Filtering == nil {
//		return "", nil
//	}
//
//	campaignFilterInfo := &CampaignGetFilterInfo{}
//
//	if len(input.Filtering.CampaignIDList) > 0 {
//		campaignIDList := input.Filtering.CampaignIDList
//		for i := 0; i < len(campaignIDList); i++ {
//			campaignFilterInfo.CampaignIds = append(campaignFilterInfo.CampaignIds, campaignIDList[i])
//		}
//	}
//
//	if len(input.Filtering.CampaignName) > 0 {
//		campaignFilterInfo.CampaignName = input.Filtering.CampaignName
//	}
//
//	if len(input.Filtering.LandingType) > 0 {
//		campaignFilterInfo.LandingType = string(input.Filtering.LandingType)
//	}
//
//	// Width
//	if input.Filtering.Width > 0 {
//		campaignFilterInfo.Width = input.Filtering.Width
//	}
//	// Height
//	if input.Filtering.Height > 0 {
//		campaignFilterInfo.Height = input.Filtering.Height
//	}
//	// CreatedStartTime
//	if len(input.Filtering.CreatedStartTime) > 0 {
//		campaignFilterInfo.StartTime = input.Filtering.CreatedStartTime
//	}
//	// CreatedEndTime
//	if len(input.Filtering.CreatedEndTime) > 0 {
//		campaignFilterInfo.EndTime = input.Filtering.CreatedEndTime
//	}
//
//	filterJson, _ := json.Marshal(campaignFilterInfo)
//	return string(filterJson), nil
//}

func (s *CampaignService) GetCampaignList(input *sdk.CampaignGetInput) (*sdk.CampaignGetOutput, error) {
	return nil, nil
	//id := formatAuthAccountID(input.BaseInput.AccountId)
	//authAccount, err := account.GetAuthAccount(id)
	//if err != nil {
	//	return nil, fmt.Errorf("GetImage get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	//}
	//
	//method := http_tools.POST
	//// create path and map variables
	//path := s.config.HttpConfig.BasePath + "/2/campaign/get/"
	//
	//header := make(map[string]string)
	//header["Accept"] = "application/json"
	//header["Access-Token"] = authAccount.AccessToken
	//header["Content-Type"] = "application/json"
	//
	//query := url.Values{}
	//
	//var request *http.Request
	//query["advertiser_id"] = []string{input.BaseInput.AccountId}
	//
	//imageFilter, err := s.getFilter(input)
	//if err != nil {
	//	return nil, err
	//}
	//if len(imageFilter) > 0 {
	//	query["filtering"] = []string{imageFilter}
	//}
	//
	//if input.Page > 0 {
	//	query["page"] = []string{strconv.FormatInt(input.Page, 10)}
	//}
	//if input.PageSize > 0 {
	//	query["page_size"] = []string{strconv.FormatInt(input.PageSize, 10)}
	//}
	//
	//request, err = s.httpClient.PrepareRequest(context.Background(), path, method, nil, header,
	//	query, nil, "", nil, "")
	//if err != nil {
	//	return nil, err
	//}
	//
	//response := &GetMaterialData{}
	//respErr := s.httpClient.DoProcess(request, response)
	//if respErr != nil {
	//	return nil, respErr
	//}
	//if response.Code != 0 {
	//	return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", response.Code,
	//		response.Message,
	//		response.RequestId)
	//}
	//
	//videoInfo := response.Data
	//videoOutput := &sdk.VideoGetOutput{}
	//s.copyVideoInfoToOutput(videoInfo, videoOutput)
	//return videoOutput, err
}
