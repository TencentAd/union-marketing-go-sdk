package orm

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"gorm.io/gorm"
)

// AuthAccountUpsert 插入或者更新授权的账号信息
func AuthAccountUpsert(db *gorm.DB, authAccount *sdk.AuthAccount) error {
	var count int64
	if  err := db.Model(authAccount).Where("account_id = ?", authAccount.AccountID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return db.Create(authAccount).Error
	} else {
		return db.Updates(authAccount).Error
	}
}

// AuthAccountUpdate 更新授权的账号信息
func AuthAccountUpdate(db *gorm.DB, authAccount *sdk.AuthAccount) error {
	return db.Updates(authAccount).Error
}

// AuthAccountGetAll 获取所有的授权账号
func AuthAccountGetAll(db *gorm.DB) ([]*sdk.AuthAccount, error) {
	var account []*sdk.AuthAccount
	if err := db.Find(&account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

// AuthAccountTake 获取授权账号
func AuthAccountTake(db *gorm.DB, id string) (*sdk.AuthAccount, error) {
	var account sdk.AuthAccount
	if err := db.Where("ID = ?", id).Take(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}
