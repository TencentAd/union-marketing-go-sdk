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
		Debug: true,
	})

	assert.NotNil(t, db)

	assert.NoError(t, orm.AuthAccountUpsert(db, &sdk.AuthAccount{
		AccountId:   1,
		AccessToken: "a",
	}))

	storage := mysql.NewTokenStorage()
	assert.NoError(t, Init(storage))

	{
		all := GetAllAuthAccount()
		assert.Len(t, all, 1)
		assert.EqualValues(t, 1, all[0].AccountId)
		assert.EqualValues(t, "a", all[0].AccessToken)
	}

	{
		account := GetAuthAccount(1)
		assert.EqualValues(t, 1, account.AccountId)
		assert.EqualValues(t, "a", account.AccessToken)
	}

	{
		assert.NoError(t, Insert(&sdk.AuthAccount{
			AccountId:   1,
			AccessToken: "a1",
		}))

		all := GetAllAuthAccount()
		assert.Len(t, all, 1)

		account := GetAuthAccount(1)
		assert.EqualValues(t, 1, account.AccountId)
		assert.EqualValues(t, "a1", account.AccessToken)
	}

	{
		assert.NoError(t, Insert(&sdk.AuthAccount{
			AccountId:   2,
			AccessToken: "b",
		}))

		all := GetAllAuthAccount()
		assert.Len(t, all, 2)
	}

	{
		assert.NoError(t, RefreshToken(&sdk.AuthAccount{
			AccountId:   1,
			AccessToken: "a2",
		}))

		account := GetAuthAccount(1)
		assert.EqualValues(t, 1, account.AccountId)
		assert.EqualValues(t, "a2", account.AccessToken)
	}
}
