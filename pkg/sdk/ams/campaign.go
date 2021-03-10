package ams

import (
	"fmt"
	"strconv"
	"time"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"github.com/antihax/optional"
	tapi "github.com/tencentad/marketing-api-go-sdk/pkg/api"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
)

type CampaignService struct {
	config *config.Config
}

// NewCampaignService
func NewCampaignService(config *config.Config) *CampaignService {
	s := &CampaignService{
		config: config,
	}
	return s
}

// getFilter 获取过滤字段
func (s *CampaignService) getFilter(input *sdk.CampaignGetInput) []model.FilteringStruct {
	if input.Filtering == nil {
		return nil
	}
	TFiltering := make([]model.FilteringStruct, 0, TFilterMax)
	// campaign_id
	mFiltering := input.Filtering
	if len(mFiltering.CampaignIDList) > 0 {
		var cIDListString []string
		for i := 0; i < len(mFiltering.CampaignIDList); i++ {
			cIDListString = append(cIDListString, strconv.FormatInt(mFiltering.CampaignIDList[i], 10))
		}
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "campaign_id",
			Operator: "IN",
			Values:   &cIDListString,
		})
	}

	if len(mFiltering.CampaignName) > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "campaign_name",
			Operator: "CONTAINS",
			Values:   &[]string{mFiltering.CampaignName},
		})
	}

	if len(mFiltering.LandingType) > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "promoted_object_type",
			Operator: "EQUALS",
			Values:   &[]string{string(input.Filtering.LandingType)},
		})
	}

	if len(mFiltering.CreateTime) > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "created_time",
			Operator: "EQUALS",
			Values:   &[]string{input.Filtering.CreateTime},
		})
	}
	return TFiltering
}

func (s *CampaignService) GetCampaignList(input *sdk.CampaignGetInput) (*sdk.CampaignGetOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId, input.BaseInput.AMSSystemType)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("GetCampaignList get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	tClient := getAMSSdkClient(authAccount)

	var campaignsGetOpts tapi.CampaignsGetOpts
	if res := s.getFilter(input); res != nil && len(res) > 0 {
		campaignsGetOpts.Filtering = optional.NewInterface(res)
	}

	// Page,Page_size
	if input.Page > 0 {
		campaignsGetOpts.Page = optional.NewInt64(input.Page)
	}
	if input.PageSize > 0 {
		campaignsGetOpts.PageSize = optional.NewInt64(input.PageSize)
	}
	if input.Filtering != nil && input.Filtering.IsDeletedAMS {
		campaignsGetOpts.IsDeleted = optional.NewBool(input.Filtering.IsDeletedAMS)
	}

	// Fields
	campaignsGetOpts.Fields = optional.NewInterface([]string{"campaign_id", "campaign_name", "configured_status",
		"campaign_type", "promoted_object_type", "daily_budget", "budget_reach_date", "created_time",
		"last_modified_time", "speed_mode", "is_deleted"})

	accID, err := strconv.ParseInt(input.BaseInput.AccountId, 10, 64)
	if err != nil {
		return nil, err
	}

	response, _, err := tClient.Campaigns().Get(*tClient.Ctx, accID, &campaignsGetOpts)
	if err != nil {
		return nil, err
	}
	output := &sdk.CampaignGetOutput{}
	s.copyCampaignInfoToOutput(&response, output)

	return output, nil
}

// copyImageToOutput 拷贝物料信息
func (s *CampaignService) copyCampaignInfoToOutput(campaignData *model.CampaignsGetResponseData, campaignOutput *sdk.CampaignGetOutput) {
	if campaignData == nil {
		return
	}
	rList := make([]*sdk.CampaignGetInfo, 0, len(*campaignData.List))
	campaignList := *campaignData.List
	for i := 0; i < len(campaignList); i++ {
		camInfo := (campaignList)[i]
		rList = append(rList, &sdk.CampaignGetInfo{
			CampaignId:         camInfo.CampaignId,
			CampaignName:       camInfo.CampaignName,
			ConfiguredStatus:   sdk.CampaignStatus(camInfo.ConfiguredStatus),
			CampaignType:       sdk.CampaignTypeAMS(camInfo.CampaignType),
			PromotedObjectType: sdk.LandingType(camInfo.PromotedObjectType),
			DailyBudget:        float32(camInfo.DailyBudget),
			BudgetReachDate:    camInfo.BudgetReachDate,
			CreatedTime:        time.Unix(camInfo.CreatedTime, 0).Format("2006-01-02 15:04:05"),
			LastModifiedTime:   time.Unix(camInfo.LastModifiedTime, 0).Format("2006-01-02 15:04:05"),
			SpeedMode:          sdk.SpeedModeAMS(camInfo.SpeedMode),
			IsDeleted:          *camInfo.IsDeleted,
		})
	}
	campaignOutput.List = rList
	campaignOutput.PageInfo = &sdk.PageConf{
		Page:        campaignData.PageInfo.Page,
		PageSize:    campaignData.PageInfo.PageSize,
		TotalNumber: campaignData.PageInfo.TotalNumber,
		TotalPage:   campaignData.PageInfo.TotalPage,
	}
}
