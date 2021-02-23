package account

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
)

type Storage interface {
	// Insert 插入一条授权的账号信息
	Insert(authAccount *sdk.AuthAccount) error

	// Upsert 更新一条授权的账号信息
	Update(authAccount *sdk.AuthAccount) error

	// List 列出已经授权的账号信息
	List() ([]*sdk.AuthAccount, error)
}

