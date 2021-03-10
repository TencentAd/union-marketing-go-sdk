package sdk

// AdGroupGetInput 获取推广组
type AdGroupGetInput struct {
	BaseInput BaseInput         `json:"base_input,omitempty"` // 账户信息
	Filtering *AdGroupFiltering `json:"filtering,omitempty"`  // 过滤信息
	Page      int64             `json:"page,omitempty"`       // 搜索页码，默认值：1 最小值 1，最大值 99999
	PageSize  int64             `json:"page_size,omitempty"`  // 一页显示的数据条数，默认值：10。最小值 1，最大值 500
}

// AdGroupFiltering 过滤粒度
type AdGroupFiltering struct {
	// 共有
	AdGroupIDList []int64       `json:"adgroup_ids,omitempty"`  // 计划id列表
	AdGroupName   []string      `json:"adgroup_name,omitempty"` // 广告组name过滤,长度为1-30个字符
	LandingType   []LandingType `json:"landing_type,omitempty"`  // 广告组推广目的过滤
	CreateTime    string        `json:"create_time,omitempty"`   // 广告组创建时间，格式yyyy-mm-dd,表示过滤出当天创建的广告组
	IsDeletedAMS   bool          `json:"is_deleted,omitempty"`    // 是否已删除，AMS特有参数
}

// CampaignGetOutput 推广计划列表
type AdGroupGetOutput struct {
	List     []*AdGroupGetInfo `json:"list,omitempty"`
	PageInfo *PageConf         `json:"page_info,omitempty"`
}
type AdGroupGetInfo struct {
	CampaignId         int64           `json:"campaign_id,omitempty"`
	AdGroupId          int64           `json:"adgroup_id,omitempty"`
	AdGroupName       string          `json:"adgroup_name,omitempty"`
	ConfiguredStatus   CampaignStatus  `json:"configured_status,omitempty"`
	PromotedObjectType LandingType     `json:"promoted_object_type,omitempty"`
	DailyBudget        int64           `json:"daily_budget,omitempty"`
	CreatedTime        string          `json:"created_time,omitempty"`
	LastModifiedTime   string          `json:"last_modified_time,omitempty"`
	IsDeleted          *bool           `json:"is_deleted,omitempty"`
}
