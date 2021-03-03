package sdk

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// Auth 授权相关接口
type Auth interface {
	// GenerateAuthURI 生成对应平台的授权链接
	GenerateAuthURI(*GenerateAuthURIInput) (*GenerateAuthURIOutput, error)

	// ProcessAuthCallback 在用户完成授权后，处理回调请求，获取用户信息以及token相关信息
	ProcessAuthCallback(input *ProcessAuthCallbackInput) ([]*ProcessAuthCallbackOutput, error)
}

// GenerateAuthURIInput
type GenerateAuthURIInput struct {
	RedirectURI string `json:"redirect_uri"` // 回调地址
}

// GenerateAuthURIOutput
type GenerateAuthURIOutput struct {
	AuthURI string `json:"auth_uri"` // 生成的授权链接
}

// ProcessAuthCallbackInput 授权回调处理输入
type ProcessAuthCallbackInput struct {
	AuthCallback *http.Request `json:"auth_callback"` // 回调请求，包含auth_code
	RedirectUri  string        `json:"redirect_uri"`  // 原始的回调地址
}

// ProcessAuthCallbackOutput 授权回调处理输出
type ProcessAuthCallbackOutput = AuthAccount

// AuthAccount 授权账号输出
type AuthAccount struct {
	Platform             MarketingPlatformType `gorm:"platform"                        json:"platform"`
	ID                   string                `gorm:"column:id;primaryKey"            json:"id,omitempty"`
	AccountUin           int64                 `gorm:"column:account_uin"              json:"account_uin,omitempty"`
	AccountID            string                `gorm:"column:account_id"               json:"account_id,omitempty"`
	ScopeList            StringList            `gorm:"column:scope_list"               json:"scope_list,omitempty"`
	WechatAccountID      string                `gorm:"column:wechat_account_id"        json:"wechat_account_id,omitempty"`
	AccountRoleType      AccountRoleType       `gorm:"column:account_role_type"        json:"account_role_type,omitempty"`
	AccountType          AccountType           `gorm:"column:account_type"             json:"account_type,omitempty"`
	AMSSystemType		AMSSystemType		   `gorm:"column:ams_system_type"          json:"ams_system_type,omitempty"`
	RoleType             RoleType              `gorm:"column:role_type"                json:"role_type,omitempty"`
	AccessToken          string                `gorm:"column:access_token"             json:"access_token,omitempty"`
	RefreshToken         string                `gorm:"column:refresh_token"            json:"refresh_token,omitempty"`
	AccessTokenExpireAt  time.Time             `gorm:"column:access_token_expires_at"  json:"access_token_expires_at,omitempty"`
	RefreshTokenExpireAt time.Time             `gorm:"column:refresh_token_expires_at" json:"refresh_token_expires_at,omitempty"`


	CreatedAt time.Time      `gorm:"column:created_at"      json:"-"`
	UpdatedAt time.Time      `gorm:"column:updated_at"      json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
}

// StringList string列表
type StringList []string

// Scan implement sql.Scanner
func (s *StringList) Scan(value interface{}) error {
	return scan(value, s)
}

// Value return json value, implement driver.Valuer interface
func (s StringList) Value() (driver.Value, error) {
	return value(s)
}

func scan(value interface{}, to interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, to)
}

func value(v interface{}) (driver.Value, error) {
	bytes, err := json.Marshal(v)
	return string(bytes), err
}

// NeedRefreshToken 判断是否需要刷新token
func (account *AuthAccount) NeedRefreshToken(f GetTokenRefreshTime) bool {
	return time.Now().After(f(account))
}

// GetTokenRefreshTimeDefault 获取账号更新时间，默认实现
func GetRefreshTimeDefault(account *AuthAccount) time.Time {
	return account.AccessTokenExpireAt
}

// GetTokenRefreshTime 获取账号的token更新时间
type GetTokenRefreshTime func(account *AuthAccount) time.Time

// RefreshToken
type RefreshToken func(account *AuthAccount) (*RefreshTokenOutput, error)

// RefreshTokenOutput 刷新token的输出
type RefreshTokenOutput struct {
	Platform             MarketingPlatformType `json:"platform"`
	ID                   string                `json:"id,omitempty"`
	AccessToken          string                `json:"access_token,omitempty"`
	RefreshToken         string                `json:"refresh_token,omitempty"`
	AccessTokenExpireAt  time.Time             `json:"access_token_expires_at,omitempty"`
	RefreshTokenExpireAt time.Time             `json:"refresh_token_expires_at,omitempty"`
}

func UpdateToken(original *AuthAccount, out *RefreshTokenOutput) {
	updateTime(&original.RefreshTokenExpireAt, out.RefreshTokenExpireAt)
	updateTime(&original.AccessTokenExpireAt, out.AccessTokenExpireAt)
	updateString(&original.RefreshToken, out.RefreshToken)
	updateString(&original.AccessToken, out.AccessToken)
}

func updateTime(original *time.Time, refreshed time.Time) {
	if refreshed.After(*original) {
		*original = refreshed
	}
}

func updateString(original *string, refreshed string) {
	if refreshed != "" {
		*original = refreshed
	}
}
