package orm

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getTestDB(tb testing.TB) (*gorm.DB, error) {
	db, err := New(&Option{
		Type:  DBTypeSQLite,
		DSN:   filepath.Join(tb.TempDir(), "sqlite"),
		Debug: !omitGormDebug(),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func omitGormDebug() bool {
	args := os.Args
	for _, arg := range args {
		if arg == "omit_gorm_debug" {
			return true
		}
	}
	return false
}

func TestTables(t *testing.T) {
	db, err := getTestDB(t)
	assert.NoError(t, err)
	assert.NoError(t, Setup(db))
}

func TestLockDB(t *testing.T) {
	db, _ := getTestDB(t)
	if db.Dialector.Name() == string(DBTypeSQLite) {
		t.Skip()
	}
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		err := db.Transaction(func(tx *gorm.DB) error {
			assert.NoError(t, LockDB(tx, "tableName", 1))

			time.Sleep(2 * time.Second)

			return UnlockDB(tx, "tableName", 1)
		})
		assert.NoError(t, err)
		wg.Done()
	}()
	go func() {
		err := db.Transaction(func(tx *gorm.DB) error {
			time.Sleep(1 * time.Second)

			return LockDB(tx, "tableName", 1)
		})
		assert.Error(t, err)
		wg.Done()
	}()

	wg.Wait()
}

func TestGetLock(t *testing.T) {
	db, _ := getTestDB(t)
	if db.Dialector.Name() == string(DBTypeSQLite) {
		t.Skip()
	}
	wg := sync.WaitGroup{}
	wg.Add(2)

	out := make(chan error)

	run := func() {
		out <- db.Transaction(func(tx *gorm.DB) error {
			return LockDB(tx, "tableName", 2)
		})
		wg.Done()
	}
	for i := 0; i < 2; i++ {
		go run()
	}
	err := <-out
	if err == nil {
		assert.Error(t, <-out)
	} else {
		assert.Nil(t, <-out)
	}

	wg.Wait()
}
