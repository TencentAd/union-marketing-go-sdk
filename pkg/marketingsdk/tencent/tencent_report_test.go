package tencent

import (
	"encoding/json"
	"fmt"
	"git.code.oa.com/tme-server-component/kg_growth_open/api"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/errors"
	"testing"
)

// 测试天级别数据
func TestGetDailyReport(t *testing.T) {
	reqInput := api.ReportInputParam{
		BaseConfig: api.BaseConfig{
			AccountId:   25610,
			AccountType: api.AccountTypeTencent,
			AccessToken: "4b647781310e83b001408e3ce092e48e",
		},
		ReportAdLevel:         api.LevelAccount,
		ReportTimeGranularity: api.ReportTimeDaily,
		ReportDateRange: api.ReportDateRange{
			StartDate: "2021-02-01",
			EndDate:   "2021-02-08",
		},
		ReportFiltering: nil,
		ReportGroupBy:   api.ADVERTISER_DATE_TENCENT,
		ReportOrderBy: api.ReportOrderBy{
			SortField: "view_count",
			SortType:  api.DESCENDING_TENCENT,
		},
		Page:     optional.NewInt64(1),
		PageSize: optional.NewInt64(100),
		Fields:   []string{"date", "view_count", "valid_click_count", "ctr", "cpc", "cost"},
	}

	reportService := &TencentReportService{}

	toutput, err := reportService.getDailyReport(&api.MarketingSDKConfig{}, &reqInput)
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
