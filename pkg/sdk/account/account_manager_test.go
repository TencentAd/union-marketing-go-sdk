package account

import (
	"path/filepath"
	"testing"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/account/mysql"
	"git.code.oa.com/tme-server-component/kg_growth_open/pkg/sdk/orm"
	"github.com/stretchr/testify/assert"
)

func Test_manager(t *testing.T) {
	db := orm.GetDB(&orm.Option{
		Type:  orm.DBTypeSQLite,
		DSN:   filepath.Join(t.TempDir(), "sqlite"),
		Debug: false,
	})

	assert.NotNil(t, db)

	assert.NoError(t, orm.AuthAccountUpsert(db, &sdk.AuthAccount{
		ID:          "1",
		AccessToken: "a",
	}))

	storage := mysql.NewTokenStorage()
	assert.NoError(t, Init(storage))

	{
		all := m.cache.getAll()
		assert.Len(t, all, 1)
		assert.EqualValues(t, "1", all[0].ID)
		assert.EqualValues(t, "a", all[0].AccessToken)
	}

	{
		account, err := GetAuthAccount("1")
		assert.NoError(t, err)
		assert.EqualValues(t, "1", account.ID)
		assert.EqualValues(t, "a", account.AccessToken)
	}

	{
		assert.NoError(t, Insert(&sdk.AuthAccount{
			ID:          "1",
			AccessToken: "a1",
		}))

		all := m.cache.getAll()
		assert.Len(t, all, 1)

		account, err := GetAuthAccount("1")
		assert.NoError(t, err)
		assert.EqualValues(t, "1", account.ID)
		assert.EqualValues(t, "a1", account.AccessToken)
	}
}
