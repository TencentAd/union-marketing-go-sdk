package orm

import (
	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"gorm.io/gorm"
)

func AuthAccountInsert(db *gorm.DB, authAccount *sdk.AuthAccount) error {
	return db.Create(authAccount).Error
}

func AuthAccountUpdate(db *gorm.DB, authAccount *sdk.AuthAccount) error {
	return db.Updates(authAccount).Error
}

func AuthAccountGetAll(db *gorm.DB) ([]*sdk.AuthAccount, error) {
	var account []*sdk.AuthAccount
	if err := db.Find(&account).Error; err != nil {
		return nil, err
	}

	return account, nil
}
