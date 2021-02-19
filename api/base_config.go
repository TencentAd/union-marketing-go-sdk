package api

type BaseConfig struct {
	AccountId   int64
	AccountType AccountType
	AccessToken string
}
type AccountType int

const (
	AccountTypeInvalid       = 0
	AccountTypeTencent       = 1 // 腾讯账户
	AccountTypeTencentWechat = 2 // 腾讯微信账户
	AccountTypeMax           = 3
)
