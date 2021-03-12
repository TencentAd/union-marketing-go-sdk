package ams

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/account"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/config"
	"github.com/antihax/optional"
	tapi "github.com/tencentad/marketing-api-go-sdk/pkg/api"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
)

type AdGroupService struct {
	config *config.Config
}

// NewAdGroupService
func NewAdGroupService(config *config.Config) *AdGroupService {
	s := &AdGroupService{
		config: config,
	}
	return s
}

// getFilter 获取过滤字段
func (s *AdGroupService) getFilter(input *sdk.AdGroupGetInput) []model.FilteringStruct {
	if input.Filtering == nil {
		return nil
	}
	var TFiltering []model.FilteringStruct
	// campaign_id
	mFiltering := input.Filtering
	if len(mFiltering.AdGroupIDList) > 0 {
		var groupListString []string
		for i := 0; i < len(mFiltering.AdGroupIDList); i++ {
			groupListString = append(groupListString, strconv.FormatInt(mFiltering.AdGroupIDList[i], 10))
		}
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "adgroup_id",
			Operator: "IN",
			Values:   &groupListString,
		})
	}

	if len(mFiltering.AdGroupName) > 0 {
		TFiltering = append(TFiltering, model.FilteringStruct{
			Field:    "adgroup_name",
			Operator: "CONTAINS",
			Values:   &[]string{mFiltering.AdGroupName},
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

func (s *AdGroupService) GetAdGroupList(input *sdk.AdGroupGetInput) (*sdk.AdGroupGetOutput, error) {
	id := formatAuthAccountID(input.BaseInput.AccountId, input.BaseInput.AMSSystemType)
	authAccount, err := account.GetAuthAccount(id)
	if err != nil {
		return nil, fmt.Errorf("GetAdGroupList get AuthAccount Info error accid=%s", input.BaseInput.AccountId)
	}

	tClient := getAMSSdkClient(authAccount)

	var adGroupsGetOpts tapi.AdgroupsGetOpts
	if res := s.getFilter(input); res != nil && len(res) > 0 {
		adGroupsGetOpts.Filtering = optional.NewInterface(res)
	}

	// Page,Page_size
	if input.Page > 0 {
		adGroupsGetOpts.Page = optional.NewInt64(input.Page)
	}
	if input.PageSize > 0 {
		adGroupsGetOpts.PageSize = optional.NewInt64(input.PageSize)
	}
	if input.Filtering != nil && input.Filtering.IsDeletedAMS {
		adGroupsGetOpts.IsDeleted = optional.NewBool(input.Filtering.IsDeletedAMS)
	}

	// Fields
	adGroupsGetOpts.Fields = optional.NewInterface([]string{"campaign_id", "adgroup_id", "adgroup_name",
		"configured_status", "promoted_object_type", "daily_budget", "created_time",
		"last_modified_time", "is_deleted"})

	accID, err := strconv.ParseInt(input.BaseInput.AccountId, 10, 64)
	if err != nil {
		return nil, err
	}

	response, _, err := tClient.Adgroups().Get(*tClient.Ctx, accID, &adGroupsGetOpts)
	if err != nil {
		return nil, err
	}
	output := &sdk.AdGroupGetOutput{}
	s.copyAdGroupInfoToOutput(&response, output)
	return output, nil
}

// copyAdGroupInfoToOutput 拷贝物料信息
func (s *AdGroupService) copyAdGroupInfoToOutput(adGroupsData *model.AdgroupsGetResponseData,
	adGroupOutput *sdk.AdGroupGetOutput) {
	if adGroupsData == nil {
		return
	}
	rList := make([]*sdk.AdGroupGetInfo, 0, len(*adGroupsData.List))
	adGroupList := *adGroupsData.List
	for i := 0; i < len(adGroupList); i++ {
		adGroupInfo := (adGroupList)[i]
		rList = append(rList, &sdk.AdGroupGetInfo{
			CampaignId:         adGroupInfo.CampaignId,
			AdGroupId:          adGroupInfo.AdgroupId,
			AdGroupName:        adGroupInfo.AdgroupName,
			AdGroupStatus:      sdk.AdGroupStatus(adGroupInfo.ConfiguredStatus),
			PromotedObjectType: sdk.LandingType(adGroupInfo.PromotedObjectType),
			DailyBudget:        float32(adGroupInfo.DailyBudget),
			CreatedTime:        time.Unix(adGroupInfo.CreatedTime, 0).Format("2006-01-02 15:04:05"),
			LastModifiedTime:   time.Unix(adGroupInfo.LastModifiedTime, 0).Format("2006-01-02 15:04:05"),
			IsDeleted:          adGroupInfo.IsDeleted,
		})
	}
	adGroupOutput.List = rList
	adGroupOutput.PageInfo = &sdk.PageConf{
		Page:        adGroupsData.PageInfo.Page,
		PageSize:    adGroupsData.PageInfo.PageSize,
		TotalNumber: adGroupsData.PageInfo.TotalNumber,
		TotalPage:   adGroupsData.PageInfo.TotalPage,
	}
}
