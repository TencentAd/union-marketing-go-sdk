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
	AdGroupIDList []int64 `json:"adgroup_ids,omitempty"`  // 广告组id列表
	AdGroupName   string  `json:"adgroup_name,omitempty"` // 广告组name过滤,长度为1-30个字符
	CreateTime    string  `json:"create_time,omitempty"`  // 广告组创建时间，格式yyyy-mm-dd,表示过滤出当天创建的广告组
	IsDeletedAMS  bool    `json:"is_deleted,omitempty"`   // 是否已删除，AMS特有参数
}

// AdGroupGetOutput 推广计划列表
type AdGroupGetOutput struct {
	List     []*AdGroupGetInfo `json:"list,omitempty"`
	PageInfo *PageConf         `json:"page_info,omitempty"`
}
type AdGroupGetInfo struct {
	CampaignId         int64         `json:"campaign_id,omitempty"`
	AdGroupId          int64         `json:"adgroup_id,omitempty"`
	AdGroupName        string        `json:"adgroup_name,omitempty"`
	AdGroupStatus      AdGroupStatus `json:"status,omitempty"`
	PromotedObjectType LandingType   `json:"promoted_object_type,omitempty"`
	DailyBudget        float32       `json:"daily_budget,omitempty"`
	CreatedTime        string        `json:"created_time,omitempty"`
	LastModifiedTime   string        `json:"last_modified_time,omitempty"`
	IsDeleted          *bool         `json:"is_deleted,omitempty"`
}

type AdGroupStatus string

const (
	// AMS
	AdGroupStatusNormal   AdGroupStatus = "AD_STATUS_NORMAL"   // 有效
	AdGroupStatusSuspend  AdGroupStatus = "AD_STATUS_SUSPEND"  // 暂停
	AdGroupStatusWithdraw AdGroupStatus = "AD_STATUS_WITHDRAW" // 提现
	AdGroupStatusPending  AdGroupStatus = "AD_STATUS_PENDING"  // 审核中（广告提交后等待进入审核）
	AdGroupStatusDenied   AdGroupStatus = "AD_STATUS_DENIED"   // 审核不通过
	AdGroupStatusFrozen   AdGroupStatus = "AD_STATUS_FROZEN"   // 已冻结（广告因状态异常，被冻结并停止投放）
	AdGroupStatusPrepare  AdGroupStatus = "AD_STATUS_PREPARE"  // 准备中
	AdGroupStatusDeleted  AdGroupStatus = "AD_STATUS_DELETED"  // 删除

	// 头条
	AdStatusDeliveryOk             AdGroupStatus = "AD_STATUS_DELIVERY_OK"              //投放中
	AdStatusDisable                AdGroupStatus = "AD_STATUS_DISABLE"                  //计划暂停
	AdStatusAudit                  AdGroupStatus = "AD_STATUS_AUDIT"                    //新建审核中
	AdStatusReaudit                AdGroupStatus = "AD_STATUS_REAUDIT"                  //修改审核中
	AdStatusDone                   AdGroupStatus = "AD_STATUS_DONE"                     //已完成（投放达到结束时间）
	AdStatusCreate                 AdGroupStatus = "AD_STATUS_CREATE"                   //计划新建
	AdStatusAuditDeny              AdGroupStatus = "AD_STATUS_AUDIT_DENY"               //审核不通过
	AdStatusBalanceExceed          AdGroupStatus = "AD_STATUS_BALANCE_EXCEED"           //账户余额不足
	AdStatusBudgetExceed           AdGroupStatus = "AD_STATUS_BUDGET_EXCEED"            //超出预算
	AdStatusNotStart               AdGroupStatus = "AD_STATUS_NOT_START"                //未到达投放时间
	AdStatusNoSchedule             AdGroupStatus = "AD_STATUS_NO_SCHEDULE"              //不在投放时段
	AdStatusAdGroupDisable         AdGroupStatus = "AD_STATUS_AdGroup_DISABLE"          //已被广告组暂停
	AdStatusAdGroupExceed          AdGroupStatus = "AD_STATUS_AdGroup_EXCEED"           //广告组超出预算
	AdStatusDelete                 AdGroupStatus = "AD_STATUS_DELETE"                   //已删除
	AdStatusAll                    AdGroupStatus = "AD_STATUS_ALL"                      //所有包含已删除
	AdStatusNotDelete              AdGroupStatus = "AD_STATUS_NOT_DELETE"               //所有不包含已删除（状态过滤默认值）
	AdStatusAdvertiserBudgetExceed AdGroupStatus = "AD_STATUS_ADVERTISER_BUDGET_EXCEED" //超出广告主日预算
)
