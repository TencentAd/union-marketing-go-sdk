package sdk

// BaseInput 基本请求信息
type BaseInput struct {
	AccountId   string           `json:"account_id,omitempty"`
	AMSSystemType AMSSystemType `json:"ams_system_type,omitempty"`
}
