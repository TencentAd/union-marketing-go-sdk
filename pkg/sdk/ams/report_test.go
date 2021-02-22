package ams

import (
	"encoding/json"
	"fmt"
	"testing"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/errors"
)

// 测试天级别数据
func TestGetDailyReport(t *testing.T) {
	reqInput := sdk.GetReportInput{
		BaseInput: sdk.BaseInput{
			AccountId:   25610,
			AccountType: sdk.AccountTypeAMS,
			AccessToken: "4b647781310e83b001408e3ce092e48e",
		},
		AdLevel:         sdk.LevelAccount,
		TimeGranularity: sdk.ReportTimeDaily,
		DateRange: sdk.DateRange{
			StartDate: "2021-02-01",
			EndDate:   "2021-02-08",
		},
		Filtering: nil,
		GroupBy:   sdk.ADVERTISER_DATE_AMS,
		OrderBy: sdk.OrderBy{
			SortField: "view_count",
			SortType:  sdk.DESCENDING_AMS,
		},
		Page:       optional.NewInt64(1),
		PageSize:   optional.NewInt64(100),
		Fields_AMS: []string{"date", "view_count", "valid_click_count", "ctr", "cpc", "cost"},
	}

	amsService := NewAMSService(&config.Config{})

	toutput, err := amsService.getDailyReport(&reqInput)
	if err != nil {
		if resErr, ok := err.(errors.ResponseError); ok {
			errStr, _ := json.Marshal(resErr)
			fmt.Println("Response error:", string(errStr))
		} else {
			fmt.Println("Error:", err)
		}
	}
	responseJson, _ := json.Marshal(toutput)
	fmt.Println("Response data:", string(responseJson))
}
