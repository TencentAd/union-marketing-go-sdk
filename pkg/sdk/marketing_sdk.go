package api

import "git.code.oa.com/tme-server-component/kg_growth_open/api/io"

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
	GetReport(reqParam *io.GetReportInput) (*io.GetReportOutput, error);
}
