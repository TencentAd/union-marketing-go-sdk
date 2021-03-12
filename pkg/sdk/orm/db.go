package orm

import (
	"sync"

	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	once      = &sync.Once{}
	singleton *gorm.DB
)

// GetDB 获取单例的数据库连接
func GetDB(options ...*Option) *gorm.DB {
	once.Do(func() {
		var option *Option = nil
		if len(options) > 0 {
			option = options[0]
		}

		db, err := New(option)
		if err != nil {
			log.Error(err)
			return
		}

		if err = Setup(db); err != nil {
			log.Error(err)
			return
		}

		singleton = db
	})
	return singleton
}

// Setup 配置数据库表结构
func Setup(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&sdk.AuthAccount{},
	); err != nil {
		return err
	}

	return nil
}
