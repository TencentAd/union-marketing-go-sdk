package sdk

type BaseInput struct {
	AccountId   int64       `json:"account_id"`
	AccountType AccountType `json:"account_type"`
	AccessToken string      `json:"account_token"`
}
