package service

import (
	"testing"

	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"go-admin-api/app/admin/models"
	"go-admin-api/common/global"
	"go-admin-api/internal/testutil"
	"gorm.io/gorm"
)

const testCasbinModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`

func newServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	return testutil.NewTestDB(
		t,
		&models.SysApi{},
		&models.SysUser{},
		&models.SysRole{},
		&models.SysMenu{},
		&models.SysDept{},
		&models.SysPost{},
		&models.CasbinRule{},
	)
}

func setServiceTestGlobals(t *testing.T) {
	t.Helper()

	originalDriver := global.Driver
	originalEnableDP := config.ApplicationConfig.EnableDP
	originalDBDriver := config.DatabaseConfig.Driver

	global.Driver = "sqlite3"
	config.ApplicationConfig.EnableDP = true
	config.DatabaseConfig.Driver = "sqlite3"

	t.Cleanup(func() {
		global.Driver = originalDriver
		config.ApplicationConfig.EnableDP = originalEnableDP
		config.DatabaseConfig.Driver = originalDBDriver
	})
}

func newRoleTestEnforcer(t *testing.T) *casbin.SyncedEnforcer {
	t.Helper()

	model, err := casbinmodel.NewModelFromString(testCasbinModel)
	if err != nil {
		t.Fatalf("create casbin model: %v", err)
	}
	enforcer, err := casbin.NewSyncedEnforcer(model)
	if err != nil {
		t.Fatalf("create casbin enforcer: %v", err)
	}
	return enforcer
}
