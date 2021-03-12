package account

import (
	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
)

// Storage 账号存储接口
type Storage interface {
	// Upsert 插入或者更新一条授权的账号信息
	Upsert(authAccount *sdk.AuthAccount) error

	// UpdateToken 更新一条授权的账号Token
	UpdateToken(authAccount *sdk.RefreshTokenOutput) error

	// List 列出已经授权的账号信息
	List() ([]*sdk.AuthAccount, error)

	// Take 根据id获取授权账号
	Take(id string) (*sdk.AuthAccount, error)
}
