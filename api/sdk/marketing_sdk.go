package sdk

type MarketingSDK interface {
	Auth // fa
	ADDelivery
	Account
	Report
}

type Auth interface {
	GetToken(input interface{}) (interface{}, error)
}

type ADDelivery interface {
	BudgetOperation
	CampaignOperation
	ADGroupOperation
	CreativeOperation
}

type BudgetOperation interface {
	GetBudget(input interface{}) (interface{}, error)
	UpdateBudget(input interface{}) (interface{}, error)
}

type CampaignOperation interface {
}

type ADGroupOperation interface {
}

type CreativeOperation interface {
}

type Account interface {
}

type Report interface {
	GetReport(reportInput *GetReportInput) (*GetReportOutput, error)
	GetVideoReport(reportInput *GetReportInput) (*GetReportOutput, error)
	GetImageReport(reportInput *GetReportInput) (*GetReportOutput, error)
}
