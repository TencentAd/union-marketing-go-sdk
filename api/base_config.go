package api

type BaseConfig struct {
	AccountId   int64 `json:"account_id"`
	Platform    PlatformType
	AccountType AccountType
	AccessToken string
}
type AccountType int

type PlatformType string

const (
//
)

const (
	AccountTypeInvalid       = 0
	AccountTypeTencent       = 1 // 腾讯账户
	AccountTypeTencentWechat = 2 // 腾讯微信账户
	AccountTypeMax           = 3
)
