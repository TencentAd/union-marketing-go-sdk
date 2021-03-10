package sdk

// CampaignGetInput 获取推广计划
type CampaignGetInput struct {
	BaseInput BaseInput          `json:"base_input,omitempty"` // 账户信息
	Filtering *CampaignFiltering `json:"filtering,omitempty"`  // 过滤信息
	Page      int64              `json:"page,omitempty"`       // 搜索页码，默认值：1 最小值 1，最大值 99999
	PageSize  int64              `json:"page_size,omitempty"`  // 一页显示的数据条数，默认值：10。最小值 1，最大值 500
}

// CampaignFiltering 过滤粒度
type CampaignFiltering struct {
	// 共有
	CampaignIDList []int64        `json:"campaign_ids,omitempty"`    // 计划id列表
	CampaignName   string         `json:"campaign_name,omitempty"`   // 广告计划name过滤,长度为1-30个字符
	LandingType    LandingType    `json:"landing_type,omitempty"`    // 广告计划推广目的过滤
	CreateTime     string         `json:"create_time,omitempty"`     // 广告计划创建时间，格式yyyy-mm-dd,表示过滤出当天创建的广告计划
	CampaignStatus CampaignStatus `json:"campaign_status,omitempty"` // 广告组状态过滤
	IsDeletedAMS   bool           `json:"is_deleted,omitempty"`      // 是否已删除，AMS特有参数
}

type LandingType string

const (
	// AMS
	PromotedObjectTypeAppAndroid              LandingType = "PROMOTED_OBJECT_TYPE_APP_ANDROID"                // Android 应用，创建广告前需通过 [promoted_objects 模块] 登记腾讯开放平台、腾讯广告上架的应用 id，创建广告时需填写之前登记的应用 id
	PromotedObjectTypeAppIos                  LandingType = "PROMOTED_OBJECT_TYPE_APP_IOS"                    // IOS 应用，创建广告前需通过 [promoted_objects 模块] 登记 App Store 的应用 id，创建广告时需填写之前登记的应用 id
	PromotedObjectTypeEcommerce               LandingType = "PROMOTED_OBJECT_TYPE_ECOMMERCE"                  // 电商推广，创建广告时无需创建和指定推广目标
	PromotedObjectTypeLinkWechat              LandingType = "PROMOTED_OBJECT_TYPE_LINK_WECHAT"                //品牌活动推广，创建广告时无需创建和指定推广目标
	PromotedObjectTypeAppAndroidMyapp         LandingType = "PROMOTED_OBJECT_TYPE_APP_ANDROID_MYAPP"          //应用宝推广，创建广告前需通过 [promoted_objects 模块] 登记腾讯应用宝的应用 id，创建广告时需填写之前登记的应用 id
	PromotedObjectTypeAppAndroidUnion         LandingType = "PROMOTED_OBJECT_TYPE_APP_ANDROID_UNION"          //Android 应用（广告包），仅可读
	PromotedObjectTypeLocalAdsWechat          LandingType = "PROMOTED_OBJECT_TYPE_LOCAL_ADS_WECHAT"           //本地广告（微信推广），创建广告前需在对应的微信公众号中注册登记门店信息，创建广告时需填写之前登记的门店 id，）门店信息的登记及获取可以通过微信公众平台提供的接口进行操作，具体方式可以参考 [本地门店的创建及获取]
	PromotedObjectTypeQqBrowserMiniProgram    LandingType = "PROMOTED_OBJECT_TYPE_QQ_BROWSER_MINI_PROGRAM"    // QQ 浏览器小程序，创建广告前需通过 [promoted_objects 模块] 登记 QQ 浏览器的小程序 id，创建广告时需填写之前登记的小程序 id
	PromotedObjectTypeLink                    LandingType = "PROMOTED_OBJECT_TYPE_LINK"                       //网页，创建广告时无需创建和指定推广目标
	PromotedObjectTypeQqMessage               LandingType = "PROMOTED_OBJECT_TYPE_QQ_MESSAGE"                 //QQ 消息，创建广告时无需创建和指定推广目标
	PromotedObjectTypeQzoneVideoPage          LandingType = "PROMOTED_OBJECT_TYPE_QZONE_VIDEO_PAGE"           //认证空间-视频说说，仅可读
	PromotedObjectTypeLocalAds                LandingType = "PROMOTED_OBJECT_TYPE_LOCAL_ADS"                  //本地广告，仅可读
	PromotedObjectTypeArticle                 LandingType = "PROMOTED_OBJECT_TYPE_ARTICLE"                    //好文广告，仅可读
	PromotedObjectTypeLeadAd                  LandingType = "PROMOTED_OBJECT_TYPE_LEAD_AD"                    //销售线索收集
	PromotedObjectTypeTencentKe               LandingType = "PROMOTED_OBJECT_TYPE_TENCENT_KE"                 //腾讯课堂，仅可读
	PromotedObjectTypeExchangeAppAndroidMyapp LandingType = "PROMOTED_OBJECT_TYPE_EXCHANGE_APP_ANDROID_MYAPP" //换量应用，仅可读
	PromotedObjectTypeQzonePageArticle        LandingType = "PROMOTED_OBJECT_TYPE_QZONE_PAGE_ARTICLE"         //QQ 空间日志页，仅可读
	PromotedObjectTypeQzonePageIframed        LandingType = "PROMOTED_OBJECT_TYPE_QZONE_PAGE_IFRAMED"         //QQ 空间嵌入页，仅可读
	PromotedObjectTypeQzonePage               LandingType = "PROMOTED_OBJECT_TYPE_QZONE_PAGE"                 //QQ 空间首页，仅可读
	PromotedObjectTypeAppPc                   LandingType = "PROMOTED_OBJECT_TYPE_APP_PC"                     //PC 应用，仅可读
	PromotedObjectTypeMiniGameWechat          LandingType = "PROMOTED_OBJECT_TYPE_MINI_GAME_WECHAT"           //微信小游戏，创建广告时需填写有效的微信小游戏 id
	PromotedObjectTypeMiniGameQq              LandingType = "PROMOTED_OBJECT_TYPE_MINI_GAME_QQ"               //QQ 小游戏
	PromotedObjectTypeAppPromotion            LandingType = "PROMOTED_OBJECT_TYPE_APP_PROMOTION"              //通用应用
	PromotedObjectTypeWechatChannels          LandingType = "PROMOTED_OBJECT_TYPE_WECHAT_CHANNELS"            //微信视频号
	// 头条
	LINK    LandingType = "LINK"    //销售线索收集
	APP     LandingType = "APP"     //应用推广
	DPA     LandingType = "DPA"     //商品目录推广
	GOODS   LandingType = "GOODS"   //商品推广（鲁班）
	STORE   LandingType = "STORE"   //门店推广
	AWEME   LandingType = "AWEME"   //抖音号推广
	SHOP    LandingType = "SHOP"    //电商店铺推广
	ARTICAL LandingType = "ARTICAL" //头条文章推广，目前API暂不支持该推广目的的设定，可在平台侧选择该推广目的类型
)

// CampaignGetOutput 推广计划列表
type CampaignGetOutput struct {
	List     []*CampaignGetInfo `json:"list,omitempty"`
	PageInfo *PageConf          `json:"page_info,omitempty"`
}
type CampaignGetInfo struct {
	CampaignId         int64           `json:"campaign_id,omitempty"`
	CampaignName       string          `json:"campaign_name,omitempty"`
	ConfiguredStatus   CampaignStatus  `json:"configured_status,omitempty"`
	CampaignType       CampaignTypeAMS `json:"campaign_type,omitempty"`
	PromotedObjectType LandingType     `json:"promoted_object_type,omitempty"`
	DailyBudget        float32           `json:"daily_budget,omitempty"`
	BudgetReachDate    int64           `json:"budget_reach_date,omitempty"`
	CreatedTime        string          `json:"created_time,omitempty"`
	LastModifiedTime   string          `json:"last_modified_time,omitempty"`
	SpeedMode          SpeedModeAMS    `json:"speed_mode,omitempty"`
	IsDeleted          bool           `json:"is_deleted,omitempty"`
}

type CampaignStatus string

// List of CampaignStatus
const (
	// AMS
	CampaignStatusNormal   CampaignStatus = "AD_STATUS_NORMAL"   // 有效
	CampaignStatusSuspend  CampaignStatus = "AD_STATUS_SUSPEND"  // 暂停
	CampaignStatusWithdraw CampaignStatus = "AD_STATUS_WITHDRAW" // 提现
	CampaignStatusPending  CampaignStatus = "AD_STATUS_PENDING"  // 审核中（广告提交后等待进入审核）
	CampaignStatusDenied   CampaignStatus = "AD_STATUS_DENIED"   // 审核不通过
	CampaignStatusFrozen   CampaignStatus = "AD_STATUS_FROZEN"   // 已冻结（广告因状态异常，被冻结并停止投放）
	CampaignStatusPrepare  CampaignStatus = "AD_STATUS_PREPARE"  // 准备中
	CampaignStatusDeleted  CampaignStatus = "AD_STATUS_DELETED"  // 删除
	// 头条
	CampaignStatusEnable                 CampaignStatus = "CAMPAIGN_STATUS_ENABLE"                   //启用
	CampaignStatusDisable                CampaignStatus = "CAMPAIGN_STATUS_DISABLE"                  //暂停
	CampaignStatusDelete                 CampaignStatus = "CAMPAIGN_STATUS_DELETE"                   //删除
	CampaignStatusAll                    CampaignStatus = "CAMPAIGN_STATUS_ALL"                      //所有包含已删除
	CampaignStatusNotDelete              CampaignStatus = "CAMPAIGN_STATUS_NOT_DELETE"               //所有不包含已删除（状态过滤默认值） \
	CampaignStatusAdvertiserBudgetExceed CampaignStatus = "CAMPAIGN_STATUS_ADVERTISER_BUDGET_EXCEED" //超出广告主日预算
)

// CampaignTypeAMS : 推广计划类型
type CampaignTypeAMS string

const (
	CampaignTypeSearch                 CampaignTypeAMS = "CAMPAIGN_TYPE_SEARCH"
	CampaignTypeNormal                 CampaignTypeAMS = "CAMPAIGN_TYPE_NORMAL"
	CampaignTypeContract               CampaignTypeAMS = "CAMPAIGN_TYPE_CONTRACT"
	CampaignTypeWechatOfficialAccounts CampaignTypeAMS = "CAMPAIGN_TYPE_WECHAT_OFFICIAL_ACCOUNTS"
	CampaignTypeWechatMoments          CampaignTypeAMS = "CAMPAIGN_TYPE_WECHAT_MOMENTS"
	CampaignTypeUnsupported            CampaignTypeAMS = "CAMPAIGN_TYPE_UNSUPPORTED"
)

// SpeedModeAMS : 投放速度模式
type SpeedModeAMS string

// List of SpeedModeAMS
const (
	SpeedModeFast            SpeedModeAMS = "SPEED_MODE_FAST"
	SpeedModeStandard        SpeedModeAMS = "SPEED_MODE_STANDARD"
	SpeedModeNone            SpeedModeAMS = "SPEED_MODE_NONE"
	SpeedModeAbsoluteUniform SpeedModeAMS = "SPEED_MODE_ABSOLUTE_UNIFORM"
)
