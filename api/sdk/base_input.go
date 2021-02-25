package sdk

// BaseInput 基本请求信息
type BaseInput struct {
	AccountId   int64           `json:"account_id"`
	AccountType AuthAccountType `json:"auth_account_type"`
	AccessToken string          `json:"account_token"`
}
