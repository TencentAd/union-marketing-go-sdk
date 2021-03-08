package sdk

// 报表回包
type GetReportOutput struct {
	List     []*ReportOutputListStruct `json:"list,omitempty"`
	PageInfo *PageConf                 `json:"page_info,omitempty"`
}

// 返回结构
type ReportOutputListStruct struct {
	AccountId            int64   `json:"account_id"`
	CampaignId           int64   `json:"campaign_id"`
	CampaignName         string  `json:"campaign_name"`
	AdgroupId            int64   `json:"adgroup_id"`
	AdgroupName          string  `json:"adgroup_name"`
	AdId                 int64   `json:"ad_id"`
	AdName               string  `json:"ad_name"`
	Date                 string  `json:"date"`
	Hour                 int64   `json:"hour"`
	ViewCount            int64   `json:"view_count"`
	DownloadCount        int64   `json:"download_count"`
	ActivatedCount       int64   `json:"activated_count"`
	ActivatedRate        float64 `json:"activated_rate"`
	ThousandDisplayPrice int64   `json:"thousand_display_price"`
	ValidClickCount      int64   `json:"valid_click_count"`
	Ctr                  float64 `json:"ctr"`
	Cpc                  int64   `json:"cpc"`
	Cost                 int64   `json:"cost"`
	KeyPageViewCost      int64   `json:"key_page_view_cost"`
	CouponClickCount     int64   `json:"coupon_click_count"`
	CouponIssueCount     int64   `json:"coupon_issue_count"`
	CouponGetCount       int64   `json:"coupon_get_count"`
	MaterialId           string  `json:"material_id,omitempty"`
}

// 分页配置信息
type PageConf struct {
	Page        int64 `json:"page,omitempty"`
	PageSize    int64 `json:"page_size,omitempty"`
	TotalNumber int64 `json:"total_number,omitempty"`
	TotalPage   int64 `json:"total_page,omitempty"`
}
