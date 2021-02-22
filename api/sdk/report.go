package sdk

import (
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
)

// 广告报表请求的参数信息
type GetReportInput struct {
	BaseInput             BaseInput             `json:"base_input,omitempty"`       // 账户信息
	ReportAdLevel         ReportAdLevel         `json:"level,omitempty"`            // 报表类型级别
	ReportTimeGranularity ReportTimeGranularity `json:"time_granularity,omitempty"` // 时间粒度
	ReportDateRange       ReportDateRange       `json:"data_range,omitempty"`       // 日期范围
	ReportFiltering       interface{}           `json:"filtering,omitempty"`        // 过滤条件
	ReportGroupBy         ReportGroupBy         `json:"groupby,omitempty"`          // 聚合条件
	ReportOrderBy         ReportOrderBy         `json:"orderby,omitempty"`          // 排序
	Page                  optional.Int64        `json:"page,omitempty"`             // 搜索页码，默认值：1 最小值 1，最大值 99999
	PageSize              optional.Int64        `json:"page_size,omitempty"`        // 一页显示的数据条数，默认值：10 最小值 1，最大值 1000
	Fields                []string              `json:"fields,omitempty"`           // 指定返回字段
}

type AccountType int

const (
	AccountTypeInvalid       = 0
	AccountTypeTencent       = 1 // 腾讯账户
	AccountTypeTencentWechat = 2 // 腾讯微信账户
	AccountTypeMax           = 3
)

// 报表的类型级别
type ReportAdLevel string

const (
	LevelAccount  ReportAdLevel = "account"
	LevelCampaign ReportAdLevel = "campaign"
	LevelAd       ReportAdLevel = "ad"
	LevelCreative ReportAdLevel = "creative"
)

// 时间粒度
type ReportTimeGranularity string

const (
	ReportTimeDaily ReportTimeGranularity = "daily"
	ReportTimeHour  ReportTimeGranularity = "hour"
)

// 日期范围
type ReportDateRange struct {
	StartDate string `json:"start_date,omitempty"` // 开始日期
	EndDate   string `json:"end_date,omitempty"`   // 结束日期
}

// 过滤粒度
type ReportFiltering struct {
	// 共有
	CampaignIDList []string `json:"campaign_ids,omitempty"`  // 计划id列表
	GroupIDList    []string `json:"adgroup_id,omitempty"`    // 组id列表
	CreativeIDList []string `json:"creative_id,omitempty"`   // 广告创意id列表
	LandingTypes   []string `json:"landing_types,omitempty"` // 推广目的列表
	// Ocean
	InventoryTypes        optional.Interface `json:"inventory_types,omitempty"`         // 投放版位
	PricingTypes          optional.Interface `json:"pricings,omitempty"`                // 出价类型
	ImageModes            optional.Interface `json:"image_modes,omitempty"`             // 素材类型列表
	CreativeMaterialModes optional.Interface `json:"creative_material_modes,omitempty"` // 创意类型列表
	FilterStatus          optional.Interface `json:"filter_status,omitempty"`           // 过滤状态：
}

type ReportGroupBy string

const (
	// Tencent
	ADVERTISER_DATE_TENCENT      ReportGroupBy = "date"                     // 按照date聚合
	CAMPAIGN_DATE_TENCENT        ReportGroupBy = "date,campaign_id"         // 按照date,campaign_id聚合
	ADGROUP_DATE_TENCENT         ReportGroupBy = "date,adgroup_id"          // 支持按date、adgroup_id、site_set 聚合
	ADGROUP_DATE_SITESET_TENCENT ReportGroupBy = "date,adgroup_id,site_set" // 支持按date、adgroup_id、site_set 聚合
	AD_DATE_TENCENT              ReportGroupBy = "date,ad_id"               // 按照date,ad_id聚合
	AD_DATE_SITESET_TENCENT      ReportGroupBy = "date,ad_id,site_set"      // 按照date、ad_id、site_set聚合

	// Oceans
	STAT_GROUP_BY_FIELD_ID ReportGroupBy = "STAT_GROUP_BY_FIELD_ID" // ID 类型-按照 ID 分组
	// TODO
)

// 排序字段结构
type ReportOrderBy struct {
	SortField string  `json:"sort_field,omitempty"`
	SortType  Sortord `json:"sort_type,omitempty"`
}

// Sortord : 排序方式
type Sortord string

const (
	ASCENDING_TENCENT  Sortord = "ASCENDING"
	DESCENDING_TENCENT Sortord = "DESCENDING"
)

// 投放版位
type InventoryTypes string

// https://ad.oceanengine.com/openapi/doc/index.html?id=528 后面补充头条和ams的版位
const (
	INVENTORY_FEED InventoryTypes = "INVENTORY_FEED" // 头条信息流（广告投放）
	// TODO 后续补充全部
)

// 出价类型
type PricingTypes string

const (
	PRICING_CPC PricingTypes = "cpc" // CPC出价
	// TODO 后续补充全部
)

// 素材类型
type ImageModes string

const (
	CREATIVE_IMAGE_MODE_SMALL ImageModes = "small_image"
	// TODO
)

// 创意类型过滤
type CreativeMaterialModes string

const (
	OCEAN_STATIC_ASSEMBLE CreativeMaterialModes = "STATIC_ASSEMBLE" // 程序化创意
	OCEAN_INTERVENE       CreativeMaterialModes = "INTERVENE"       // 自定义创意
)

// 推广目标
type LandingTypes string

const (
	LINK LandingTypes = "link_ocean" // 销售线索收集
)

// 状态，包括计划，组，创意状态
type FilterStatus string

const (
	OCEAN_CREATIVE_STATUS_DELIVERY_OK FilterStatus = "creative_status_ok" // 创意投放中
)

// 报表回包
type GetReportOutput struct {
	// Tencent response
	TencentReportResponse *model.DailyReportsGetResponseData
	// Tencent hour response  TODO 需要跟天级别的数据合并
	TencentHourlyReportResponse *model.HourlyReportsGetResponseData
}
