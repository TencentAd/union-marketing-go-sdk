package ocean_engine

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/account"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/config"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/http_tools"
)

type AdGroupService struct {
	config     *config.Config
	httpClient *http_tools.HttpClient
}


// NewAdGroupService 获取广告组服务
func NewAdGroupService(sConfig *config.Config) *AdGroupService {
	return &AdGroupService{
		config:     sConfig,
		httpClient: http_tools.Init(sConfig.HttpConfig),
	}
}

type AdGroupGetFilterInfo struct {
	AdGroupIds        []int64 `json:"ids,omitempty"`                  // 广告组ID过滤，数组，不超过100个
	AdGroupName       string  `json:"campaign_name,omitempty"`        // 广告组name过滤，长度为1-30个字符，其中1个中文字符算2位
	LandingType       string  `json:"landing_type,omitempty"`         // 广告组推广目的过滤
	AdGroupCreateTime string  `json:"campaign_create_time,omitempty"` // 广告组创建时间，格式yyyy-mm-dd，表示过滤出当天创建的广告组
	AdGroupStatus     string  `json:"status,omitempty"`               // 广告组状态过滤，默认为返回“所有不包含已删除”
}

// getFilter 获取过滤信息
func (s *AdGroupService) getFilter(input *sdk.AdGroupGetInput) string {
	if input.Filtering == nil {
		return ""
	}

	adGroupFilterInfo := &AdGroupGetFilterInfo{}

	if len(input.Filtering.AdGroupIDList) > 0 {
		adGroupIDList := input.Filtering.AdGroupIDList
		for i := 0; i < len(adGroupIDList); i++ {
			adGroupFilterInfo.AdGroupIds = append(adGroupFilterInfo.AdGroupIds, adGroupIDList[i])
		}
	}

	if len(input.Filtering.AdGroupName) > 0 {
		adGroupFilterInfo.AdGroupName = input.Filtering.AdGroupName
	}

	// CreatedEndTime
	if len(input.Filtering.CreateTime) > 0 {
		adGroupFilterInfo.AdGroupCreateTime = input.Filtering.CreateTime
	}

	filterJson, _ := json.Marshal(adGroupFilterInfo)
	return string(filterJson)
}

type GetAdGroupResponse struct {
	Code      int             `json:"code"`
	Message   string          `json:"message"`
	Data      *GetAdGroupList `json:"data"`
	RequestId string          `json:"request_id"`
}
type GetAdGroupList struct {
	AdGroupList []*GetAdGroupInfo `json:"list"`
	PageInfo    *sdk.PageConf     `json:"page_info"`
}

type GetAdGroupInfo struct {
	ID                 int64   `json:"id"`                   //广告组ID
	Name               string  `json:"name"`                 //广告组名称
	CampaignId         int64   `json:"campaign_id"`          //广告计划ID
	Budget             float32 `json:"budget"`               //广告组预算
	BudgetMode         string  `json:"budget_mode"`          //广告组预算类型, 详见【附录-预算类型】
	LandingType        string  `json:"landing_type"`         //广告组推广目的，详见【附录-推广目的类型】
	ModifyTime         string  `json:"modify_time"`          //广告组时间戳,用于更新时提交,服务端判断是否基于最新信息修改
	Status             string  `json:"status"`               //广告组状态,详见【附录-广告组状态】
	AdGroupCreateTime  string  `json:"campaign_create_time"` //广告组创建时间, 格式：yyyy-mm-dd hh:MM:ss
	AdGroupModifyTime  string  `json:"campaign_modify_time"` //广告组修改时间, 格式：yyyy-mm-dd hh:MM:ss
	DeliveryRelatedNum string  `json:"delivery_related_num"` //广告组商品类型，详见【附录-广告组商品类型】
	DeliveryMode       string  `json:"delivery_mode"`        //投放类型，允许值：MANUAL（手动）、PROCEDURAL（自动，投放管家）
}

func (s *AdGroupService) GetAdGroupList(input *sdk.AdGroupGetInput) (*sdk.AdGroupGetOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("GetImage get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	method := http_tools.GET
	// create path and map variables
	path := s.config.HttpConfig.BasePath + "/2/ad/get/"

	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["Access-Token"] = authAccount.AccessToken
	header["Content-Type"] = "application/json"

	query := url.Values{}

	var request *http.Request
	query["advertiser_id"] = []string{input.BaseInput.AccountId}

	adGroupFilter := s.getFilter(input)
	if len(adGroupFilter) > 0 {
		query["filtering"] = []string{adGroupFilter}
	}

	fieldList := []string{"campaign_id", "id", "name", "status", "budget", "ad_create_time", "ad_modify_time"}

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

	response := &GetAdGroupResponse{}
	respErr := s.httpClient.DoProcess(request, response)
	if respErr != nil {
		return nil, respErr
	}
	if response.Code != 0 {
		return nil, fmt.Errorf("response : code = %d, message = %s, request_id = %s ", response.Code,
			response.Message,
			response.RequestId)
	}

	adGroupData := response.Data
	adGroupGetOutput := &sdk.AdGroupGetOutput{}
	s.copyAdGroupDataToOutput(adGroupData, adGroupGetOutput)
	return adGroupGetOutput, err
}

// copyAccountReportToOutput 拷贝账户报表
func (s *AdGroupService) copyAdGroupDataToOutput(input *GetAdGroupList,
	output *sdk.AdGroupGetOutput) {
	if len(input.AdGroupList) == 0 {
		return
	}

	rList := make([]*sdk.AdGroupGetInfo, 0, len(input.AdGroupList))
	for i := 0; i < len(input.AdGroupList); i++ {
		adGroupInfo := (input.AdGroupList)[i]
		rList = append(rList, &sdk.AdGroupGetInfo{
			CampaignId:       adGroupInfo.CampaignId,
			AdGroupId:        adGroupInfo.ID,
			AdGroupName:      adGroupInfo.Name,
			AdGroupStatus:    sdk.AdGroupStatus(adGroupInfo.Status),
			DailyBudget:      adGroupInfo.Budget,
			CreatedTime:      adGroupInfo.AdGroupCreateTime,
			LastModifiedTime: adGroupInfo.AdGroupModifyTime,
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
