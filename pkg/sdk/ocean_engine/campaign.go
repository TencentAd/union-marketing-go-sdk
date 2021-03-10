package ocean_engine

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
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
	CampaignStatus     string  `json:"status,omitempty"`               // 广告组状态过滤，默认为返回“所有不包含已删除”
}

// getFilter 获取过滤信息
func (s *CampaignService) getFilter(input *sdk.CampaignGetInput) (string, error) {
	if input.Filtering == nil {
		return "", nil
	}

	campaignFilterInfo := &CampaignGetFilterInfo{}

	if len(input.Filtering.CampaignIDList) > 0 {
		campaignIDList := input.Filtering.CampaignIDList
		for i := 0; i < len(campaignIDList); i++ {
			campaignFilterInfo.CampaignIds = append(campaignFilterInfo.CampaignIds, campaignIDList[i])
		}
	}

	if len(input.Filtering.CampaignName) > 0 {
		campaignFilterInfo.CampaignName = input.Filtering.CampaignName
	}

	if len(input.Filtering.LandingType) > 0 {
		campaignFilterInfo.LandingType = string(input.Filtering.LandingType)
	}

	if len(input.Filtering.CampaignStatus) > 0 {
		campaignFilterInfo.CampaignStatus = string(input.Filtering.CampaignStatus)
	}

	// CreatedEndTime
	if len(input.Filtering.CreateTime) > 0 {
		campaignFilterInfo.CampaignCreateTime = input.Filtering.CreateTime
	}

	filterJson, _ := json.Marshal(campaignFilterInfo)
	return string(filterJson), nil
}

type GetCampaignResponse struct {
	Code      int              `json:"code"`
	Message   string           `json:"message"`
	Data      *GetCampaignList `json:"data"`
	RequestId string           `json:"request_id"`
}
type GetCampaignList struct {
	CampaignList []*GetCampaignInfo `json:"list"`
	PageInfo     *sdk.PageConf      `json:"page_info"`
}

type GetCampaignInfo struct {
	ID                 int64  `json:"id"`                   //广告组ID
	Name               string `json:"name"`                 //广告组名称
	Budget             float32  `json:"budget"`               //广告组预算
	BudgetMode         string `json:"budget_mode"`          //广告组预算类型, 详见【附录-预算类型】
	LandingType        string `json:"landing_type"`         //广告组推广目的，详见【附录-推广目的类型】
	ModifyTime         string `json:"modify_time"`          //广告组时间戳,用于更新时提交,服务端判断是否基于最新信息修改
	Status             string `json:"status"`               //广告组状态,详见【附录-广告组状态】
	CampaignCreateTime string `json:"campaign_create_time"` //广告组创建时间, 格式：yyyy-mm-dd hh:MM:ss
	CampaignModifyTime string `json:"campaign_modify_time"` //广告组修改时间, 格式：yyyy-mm-dd hh:MM:ss
	DeliveryRelatedNum string `json:"delivery_related_num"` //广告组商品类型，详见【附录-广告组商品类型】
	DeliveryMode       string `json:"delivery_mode"`        //投放类型，允许值：MANUAL（手动）、PROCEDURAL（自动，投放管家）
}

func (s *CampaignService) GetCampaignList(input *sdk.CampaignGetInput) (*sdk.CampaignGetOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("GetImage get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	method := http_tools.GET
	// create path and map variables
	path := s.config.HttpConfig.BasePath + "/2/campaign/get/"

	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["Access-Token"] = authAccount.AccessToken
	header["Content-Type"] = "application/json"

	query := url.Values{}

	var request *http.Request
	query["advertiser_id"] = []string{input.BaseInput.AccountId}

	campaignFilter, err := s.getFilter(input)
	if err != nil {
		return nil, err
	}
	if len(campaignFilter) > 0 {
		query["filtering"] = []string{campaignFilter}
	}

	fieldList := []string{"id", "name", "budget", "budget_mode", "landing_type", "status", "modify_time",
		"modify_time", "campaign_modify_time", "campaign_create_time"}
	field, _ := json.Marshal(fieldList)
	query["field"] = []string{string(field)}

	if input.Page > 0 {
		query["page"] = []string{strconv.FormatInt(input.Page, 10)}
	}
	if input.PageSize > 0 {
		query["page_size"] = []string{strconv.FormatInt(input.PageSize, 10)}
	}

	request, err = s.httpClient.PrepareRequest(context.Background(), path, method, nil, header,
		query, nil, "", nil, "")
	if err != nil {
		return nil, err
	}

	response := &GetCampaignResponse{}
	respErr := s.httpClient.DoProcess(request, response)
	if respErr != nil {
		return nil, respErr
	}
	if response.Code != 0 {
		return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", response.Code,
			response.Message,
			response.RequestId)
	}

	campaignData := response.Data
	campaignGetOutput := &sdk.CampaignGetOutput{}
	s.copyCampaignDataToOutput(campaignData, campaignGetOutput)
	return campaignGetOutput, err
}

// copyAccountReportToOutput 拷贝账户报表
func (s *CampaignService) copyCampaignDataToOutput(input *GetCampaignList,
	output *sdk.CampaignGetOutput) {
	if len(input.CampaignList) == 0 {
		return
	}
	rList := make([]*sdk.CampaignGetInfo, 0, len(input.CampaignList))
	for i := 0; i < len(input.CampaignList); i++ {
		campaignInfo := (input.CampaignList)[i]
		rList = append(rList, &sdk.CampaignGetInfo{
			CampaignId:         campaignInfo.ID,
			CampaignName:       campaignInfo.Name,
			ConfiguredStatus:   sdk.CampaignStatus(campaignInfo.Status),
			PromotedObjectType: sdk.LandingType(campaignInfo.LandingType),
			DailyBudget:        campaignInfo.Budget,
			CreatedTime:        campaignInfo.CampaignCreateTime,
		})
	}
	output.List = rList
	output.PageInfo = &sdk.PageConf{
		Page:        input.PageInfo.Page,
		PageSize:    input.PageInfo.PageSize,
		TotalNumber: input.PageInfo.TotalNumber,
		TotalPage:   input.PageInfo.TotalPage,
	}
}
