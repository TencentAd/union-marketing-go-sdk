package sdk

import "net/http"

type Auth interface {
	// 授权接口，在用户完成授权后，获取用户信息以及token相关信息
	ServeAuth(w http.ResponseWriter, req *http.Request)
}

// AuthResponse 授权输出
type AuthResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    *AuthAccountOutput `json:"data"`
}

// AuthAccountOutput 授权账号输出
type AuthAccountOutput struct {
	AccountUin            int64           `json:"account_uin,omitempty"`
	AccountId             int64           `json:"account_id,omitempty"`
	ScopeList             *[]string       `json:"scope_list,omitempty"`
	WechatAccountId       string          `json:"wechat_account_id,omitempty"`
	AccountRoleType       AccountRoleType `json:"account_role_type,omitempty"`
	AccountType           AccountType     `json:"account_type,omitempty"`
	RoleType              RoleType        `json:"role_type,omitempty"`
	AccessToken           string          `json:"access_token,omitempty"`
	RefreshToken          string          `json:"refresh_token,omitempty"`
	AccessTokenExpiresIn  int64           `json:"access_token_expires_in,omitempty"`
	RefreshTokenExpiresIn int64           `json:"refresh_token_expires_in,omitempty"`
}
