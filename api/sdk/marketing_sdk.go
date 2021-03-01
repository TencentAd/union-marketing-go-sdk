package sdk

// MarketingSDK 对Marketing API的抽象
type MarketingSDK interface {
	Auth // 授权接口
	//ADDelivery
	//Account
	Report
}

// ADDelivery 广告投放接口
type ADDelivery interface {
	BudgetOperation
	CampaignOperation
	ADGroupOperation
	CreativeOperation
}

// BudgetOperation 预算相关接口
type BudgetOperation interface {
}

// CampaignOperation 推广计划相关接口
type CampaignOperation interface {
}

// ADGroupOperation 广告组相关接口
type ADGroupOperation interface {
}

// CreativeOperation 创意相关接口
type CreativeOperation interface {
}

// Account 账户管理接口
type Account interface {
}

// Material 物料管理接口
type Material interface {
	AddImage()
	GetImage()
	AddVideo()
	GetVideo()
}

// Report 报表相关接口
type Report interface {
	// GetReport 获取报表
	GetReport(reportInput *GetReportInput) (*GetReportOutput, error)
	//GetVideoReport(reportInput *GetReportInput) (*GetReportOutput, error)
	//GetImageReport(reportInput *GetReportInput) (*GetReportOutput, error)
}
