package ams

import (
	"encoding/json"
	"fmt"
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/config"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/errors"
	"testing"
)

// 测试天级别数据
func TestGetDailyReport(t *testing.T) {
	reqInput := sdk.GetReportInput{
		BaseReportInput: sdk.BaseReportInput{
			AccountId:   25610,
			AccountType: sdk.AccountTypeTencent,
			AccessToken: "4b647781310e83b001408e3ce092e48e",
		},
		ReportAdLevel:         sdk.LevelAccount,
		ReportTimeGranularity: sdk.ReportTimeDaily,
		ReportDateRange: sdk.ReportDateRange{
			StartDate: "2021-02-01",
			EndDate:   "2021-02-08",
		},
		ReportFiltering: nil,
		ReportGroupBy:   sdk.ADVERTISER_DATE_TENCENT,
		ReportOrderBy: sdk.ReportOrderBy{
			SortField: "view_count",
			SortType:  sdk.DESCENDING_TENCENT,
		},
		Page:     optional.NewInt64(1),
		PageSize: optional.NewInt64(100),
		Fields:   []string{"date", "view_count", "valid_click_count", "ctr", "cpc", "cost"},
	}

	amsReportService := NewAMSReportService(&config.Config{})

	toutput, err := amsReportService.getDailyReport(&reqInput)
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
