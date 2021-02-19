package tencent

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

type GroupByTypes string

const (
	// Tencent
	ADVERTISER_DATE_TENCENT GroupByTypes = "date" // 按照date聚合
	CAMPAIGN_DATE_TENCENT GroupByTypes = "date,campaign_id" // 按照date,campaign_id聚合
	ADGROUP_DATE_TENCENT GroupByTypes = "date,adgroup_id" // 支持按date、adgroup_id、site_set 聚合
	ADGROUP_DATE_SITESET_TENCENT GroupByTypes = "date,adgroup_id,site_set" // 支持按date、adgroup_id、site_set 聚合
	AD_DATE_TENCENT GroupByTypes = "date,ad_id" // 按照date,ad_id聚合
	AD_DATE_SITESET_TENCENT GroupByTypes = "date,ad_id,site_set" // 按照date、ad_id、site_set聚合

	// Oceans
	STAT_GROUP_BY_FIELD_ID GroupByTypes = "STAT_GROUP_BY_FIELD_ID" // ID 类型-按照 ID 分组
	// TODO
)
