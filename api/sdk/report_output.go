package sdk

// 报表回包
type GetReportOutput struct {
	List     *[]ReportOutputListStruct `json:"list,omitempty"`
	PageInfo *PageConf                 `json:"page_info,omitempty"`
}

// 返回结构
type ReportOutputListStruct struct {
	AccountId            int64   `json:"account_id,omitempty"`
	Date                 string  `json:"date,omitempty"`
	Hour                 int64   `json:"hour,omitempty"`
	ViewCount            int64   `json:"view_count,omitempty"`
	DownloadCount        int64   `json:"download_count,omitempty"`
	ActivatedCount       int64   `json:"activated_count,omitempty"`
	ActivatedRate        float64 `json:"activated_rate,omitempty"`
	ThousandDisplayPrice int64   `json:"thousand_display_price,omitempty"`
	ValidClickCount      int64   `json:"valid_click_count,omitempty"`
	Ctr                  float64 `json:"ctr,omitempty"`
	Cpc                  int64   `json:"cpc,omitempty"`
	Cost                 int64   `json:"cost,omitempty"`
	KeyPageViewCost      int64   `json:"key_page_view_cost,omitempty"`
	CouponClickCount     int64   `json:"coupon_click_count,omitempty"`
	CouponIssueCount     int64   `json:"coupon_issue_count,omitempty"`
	CouponGetCount       int64   `json:"coupon_get_count,omitempty"`
}

// 分页配置信息
type PageConf struct {
	Page        int64 `json:"page,omitempty"`
	PageSize    int64 `json:"page_size,omitempty"`
	TotalNumber int64 `json:"total_number,omitempty"`
	TotalPage   int64 `json:"total_page,omitempty"`
}
