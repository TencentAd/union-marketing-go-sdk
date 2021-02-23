package sdk

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Auth interface {
	// 授权接口，在用户完成授权后，获取用户信息以及token相关信息
	ServeAuth(w http.ResponseWriter, req *http.Request)
}

// AuthResponse 授权输出
type AuthResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    *AuthAccount `json:"data"`
}

// AuthAccount 授权账号输出
type AuthAccount struct {
	AccountUin           int64           `gorm:"column:account_uin"              json:"account_uin,omitempty"`
	AccountId            int64           `gorm:"column:account_id;primaryKey"    json:"account_id,omitempty"`
	ScopeList            StringList      `gorm:"column:scope_list"               json:"scope_list,omitempty"`
	WechatAccountId      string          `gorm:"column:wechat_account_id"        json:"wechat_account_id,omitempty"`
	AccountRoleType      AccountRoleType `gorm:"column:account_role_type"        json:"account_role_type,omitempty"`
	AccountType          AccountType     `gorm:"column:account_type"             json:"account_type,omitempty"`
	RoleType             RoleType        `gorm:"column:role_type"                json:"role_type,omitempty"`
	AccessToken          string          `gorm:"column:access_token"             json:"access_token,omitempty"`
	RefreshToken         string          `gorm:"column:refresh_token"            json:"refresh_token,omitempty"`
	AccessTokenExpireAt  time.Time       `gorm:"column:access_token_expires_at"  json:"access_token_expires_at,omitempty"`
	RefreshTokenExpireAt time.Time       `gorm:"column:refresh_token_expires_at" json:"refresh_token_expires_at,omitempty"`

	CreatedAt time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
}

type StringList []string

func (s *StringList) Scan(value interface{}) error {
	return scan(value, s)
}

// Value return json value, implement driver.Valuer interface
func (s *StringList) Value() (driver.Value, error) {
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

func (aa *AuthAccount) Key() string {
	return strconv.FormatInt(aa.AccountId, 10)
}
