package sdk

import sdkconfig "github.com/tencentad/union-marketing-go-sdk/pkg/sdk/config"

// MarketingSDK 对Marketing API的抽象
type MarketingSDK interface {
	GetConfig() *sdkconfig.Config
	Auth // 授权接口
	ADDelivery
	Account
	Report
	Material
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
	GetCampaignList(input *CampaignGetInput) (*CampaignGetOutput, error)
}

// ADGroupOperation 广告组相关接口
type ADGroupOperation interface {
	GetAdGroupList(input *AdGroupGetInput) (*AdGroupGetOutput, error)
}

// CreativeOperation 创意相关接口
type CreativeOperation interface {
}

// Account 账户管理接口
type Account interface {
	GetAuthAccount(input *BaseInput) (*AuthAccount, error)
}

// Material 物料管理接口
type Material interface {
	AddImage(input *ImageAddInput) (*ImagesAddOutput, error)
	GetImage(input *MaterialGetInput) (*ImageGetOutput, error)
	AddVideo(input *VideoAddInput) (*VideoAddOutput, error)
	GetVideo(input *MaterialGetInput) (*VideoGetOutput, error)
	BindMaterial(input *MaterialBindInput) (*MaterialBindOutput, error)
}

// Report 报表相关接口
type Report interface {
	// GetReport 获取报表
	GetReport(input *GetReportInput) (*GetReportOutput, error)
	GetVideoReport(input *GetReportInput) (*GetReportOutput, error)
	GetImageReport(input *GetReportInput) (*GetReportOutput, error)
}
