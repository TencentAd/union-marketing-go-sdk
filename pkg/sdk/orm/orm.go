package orm

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DBTypeSQLite  dbType = "sqlite"
	DBTypeMySQL   dbType = "mysql"
	DBTypePostgre dbType = "postgres"
)

var (
	DefaultOrmOption = &Option{
		Type: DBTypeSQLite,
		DSN:  "sqlite.db",
	}
)

type dbType string

// Option 数据库配置
type Option struct {
	Debug bool   `json:"debug"`
	DSN   string `json:"dsn"`
	Type  dbType `json:"type"`
}

// New 创建数据库实例
func New(option *Option) (*gorm.DB, error) {
	if option == nil {
		option = DefaultOrmOption
	}

	dialect, err := getDialect(option.Type, option.DSN)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if option.Debug {
		db = db.Debug()
	}

	return db, nil
}

// LockDB 锁数据库db
func LockDB(db *gorm.DB, resource int64) error {
	ty := dbType(db.Dialector.Name())
	switch ty {
	case DBTypeMySQL:
		// docs: https://dev.mysql.com/doc/refman/8.0/en/locking-functions.html
		var res int
		db.Raw("SELECT GET_LOCK(?, 0) WHERE (SELECT IS_FREE_LOCK(?))=1;", resource, resource).Scan(&res)
		if res == 0 {
			return fmt.Errorf("%v has been locked", resource)
		}

		return nil
	case DBTypePostgre:
		// docs: http://www.postgres.cn/docs/9.3/functions-admin.html
		var res bool
		db.Raw("SELECT pg_try_advisory_lock(?);", resource).Scan(&res)
		if !res {
			return fmt.Errorf("%v has been locked", resource)
		}

		return nil
	default:
		return fmt.Errorf("unsupported db type: %v", ty)
	}
}

// UnlockDB 解锁数据库db中的resource
func UnlockDB(db *gorm.DB, resource int64) error {
	ty := dbType(db.Dialector.Name())
	switch ty {
	case DBTypeMySQL:
		var res int
		if err := db.Raw("SELECT RELEASE_LOCK(?);", resource).Scan(&res).Error; err != nil {
			return err
		}
		if res != 1 {
			return fmt.Errorf("%v has been unlocked", resource)
		}
	case DBTypePostgre:
		return db.Raw("SELECT pg_advisory_unlock(?);", resource).Error
	default:
		return fmt.Errorf("unsupported db type: %v", ty)
	}

	return nil
}

func getDialect(ty dbType, dsn string) (gorm.Dialector, error) {
	switch ty {
	case DBTypeSQLite:
		return sqlite.Open(dsn), nil
	case DBTypeMySQL:
		return mysql.Open(dsn), nil
	case DBTypePostgre:
		return postgres.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported db type: %v", ty)
	}
}
