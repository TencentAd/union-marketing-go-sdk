package sdk

// 广告报表请求的参数信息
type GetReportInput struct {
	BaseInput       BaseInput       `json:"base_input,omitempty"`       // 账户信息
	AdLevel         AdLevel         `json:"level,omitempty"`            // 报表类型级别
	TimeGranularity TimeGranularity `json:"time_granularity,omitempty"` // 时间粒度
	DateRange       DateRange       `json:"data_range,omitempty"`       // 日期范围
	Filtering       *Filtering      `json:"filtering,omitempty"`        // 过滤条件
	GroupBy         *[]GroupBy      `json:"group_by,omitempty"`         // 聚合条件
	OrderBy         OrderBy         `json:"order_by,omitempty"`         // 排序
	Page            int64           `json:"page,omitempty"`             // 搜索页码，默认值：1 最小值 1，最大值 99999
	PageSize        int64           `json:"page_size,omitempty"`        // 一页显示的数据条数，默认值：10 最小值 1，最大值 1000
}

// 报表的类型级别
type AdLevel string

const (
	LevelAccount  AdLevel = "account"
	LevelCampaign AdLevel = "campaign"
	LevelAd       AdLevel = "ad"
	LevelCreative AdLevel = "creative"
	LevelVideo    AdLevel = "video"
	LevelImage    AdLevel = "image"
)

// 时间粒度
type TimeGranularity string

const (
	ReportTimeDaily TimeGranularity = "daily"
	ReportTimeHour  TimeGranularity = "hour"
)

// 日期范围
type DateRange struct {
	StartDate string `json:"start_date,omitempty"` // 开始日期
	EndDate   string `json:"end_date,omitempty"`   // 结束日期
}

// 过滤粒度
type Filtering struct {
	// 共有
	CampaignIDList []string `json:"campaign_ids,omitempty"` // 计划id列表
	AdIDList       []string `json:"ad_id,omitempty"`        // 广告id列表
	CreativeIDList []string `json:"creative_id,omitempty"`  // 广告创意id列表
}

// GroupBy类型
type GroupBy string

// ID维度和时间维度可以组合查询
const (
	ADVERTISER_ID GroupBy = "ADVERTISER_ID" // advertiser聚合
	CAMPAIGN_ID   GroupBy = "CAMPAIGN_ID"   // 按照campaign_id聚合
	AD_ID         GroupBy = "AD_ID"         // 支持按ad_id聚合
	CREATIVE_ID   GroupBy = "CREATIVE_ID"   // 按照creative_id聚合
	Material_ID   GroupBy = "MATERIAL_ID"   // 按照material_id聚合
	DATE       GroupBy = "DATE"       // 按照DATE聚合
	HOUR       GroupBy = "HOUR"       // 按照HOUR聚合
)

// 排序字段结构
type OrderBy struct {
	SortField string   `json:"sort_field,omitempty"`
	SortType  SortType `json:"sort_type,omitempty"`
}

// Sortord : 排序方式
type SortType string

const (
	ASC  SortType = "asc"
	DESC SortType = "desc"
)
