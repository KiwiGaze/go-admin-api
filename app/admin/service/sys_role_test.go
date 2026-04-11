package service

import (
	"reflect"
	"testing"

	"go-admin-api/app/admin/models"
	servicedto "go-admin-api/app/admin/service/dto"
	"go-admin-api/internal/testutil"
)

func TestSysRoleGetScenarios(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysRole{Service: testutil.NewTestService(db)}

	var missing models.SysRole
	if err := service.Get(&servicedto.SysRoleGetReq{Id: 404}, &missing); err == nil {
		t.Fatal("Get() error = nil, want missing role error")
	}

	menu1 := models.SysMenu{MenuId: 1, Title: "Root", MenuType: "M"}
	menu2 := models.SysMenu{MenuId: 2, Title: "Page", MenuType: "C"}
	role := models.SysRole{RoleId: 9, RoleName: "editor", RoleKey: "editor", Status: "2"}
	for _, seed := range []interface{}{&menu1, &menu2, &role} {
		if err := db.Create(seed).Error; err != nil {
			t.Fatalf("seed role get data: %v", err)
		}
	}
	for _, pair := range []map[string]int{
		{"role_id": 9, "menu_id": 1},
		{"role_id": 9, "menu_id": 2},
	} {
		if err := db.Exec("INSERT INTO sys_role_menu (role_id, menu_id) VALUES (?, ?)", pair["role_id"], pair["menu_id"]).Error; err != nil {
			t.Fatalf("seed sys_role_menu: %v", err)
		}
	}

	var loaded models.SysRole
	if err := service.Get(&servicedto.SysRoleGetReq{Id: 9}, &loaded); err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !reflect.DeepEqual(loaded.MenuIds, []int{1, 2}) {
		t.Fatalf("loaded menu ids = %v, want [1 2]", loaded.MenuIds)
	}
}

func TestSysRoleInsertRejectsDuplicateRoleKey(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysRole{Service: testutil.NewTestService(db)}
	enforcer := newRoleTestEnforcer(t)

	if err := db.Create(&models.SysRole{RoleId: 1, RoleName: "existing", RoleKey: "editor", Status: "2"}).Error; err != nil {
		t.Fatalf("seed role: %v", err)
	}

	err := service.Insert(&servicedto.SysRoleInsertReq{
		RoleName: "duplicate",
		RoleKey:  "editor",
		Status:   "2",
	}, enforcer)
	if err == nil || err.Error() != "roleKey already exists; change it before submitting" {
		t.Fatalf("Insert() error = %v, want duplicate roleKey error", err)
	}
}

func TestSysRoleInsertAssociatesMenusAndPolicySync(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysRole{Service: testutil.NewTestService(db)}
	enforcer := newRoleTestEnforcer(t)

	api1 := models.SysApi{Id: 10, Title: "List", Path: "/users", Action: "GET"}
	api2 := models.SysApi{Id: 11, Title: "Create", Path: "/users", Action: "POST"}
	menu1 := models.SysMenu{MenuId: 1, Title: "Users", MenuType: "C"}
	menu2 := models.SysMenu{MenuId: 2, Title: "Users 2", MenuType: "C"}
	for _, seed := range []interface{}{&api1, &api2, &menu1, &menu2} {
		if err := db.Create(seed).Error; err != nil {
			t.Fatalf("seed insert data: %v", err)
		}
	}
	if err := db.Model(&menu1).Association("SysApi").Append(&api1, &api2); err != nil {
		t.Fatalf("append menu1 apis: %v", err)
	}
	if err := db.Model(&menu2).Association("SysApi").Append(&api1); err != nil {
		t.Fatalf("append menu2 apis: %v", err)
	}

	err := service.Insert(&servicedto.SysRoleInsertReq{
		RoleId:   9,
		RoleName: "editor",
		RoleKey:  "editor",
		Status:   "2",
		MenuIds:  []int{1, 2},
	}, enforcer)
	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	menuIDs, err := service.GetRoleMenuId(9)
	if err != nil {
		t.Fatalf("GetRoleMenuId() error = %v", err)
	}
	if !reflect.DeepEqual(menuIDs, []int{1, 2}) {
		t.Fatalf("role menu ids = %v, want [1 2]", menuIDs)
	}

	policies, err := enforcer.GetFilteredNamedPolicy("p", 0, "editor")
	if err != nil {
		t.Fatalf("GetFilteredNamedPolicy() error = %v", err)
	}
	wantPolicies := [][]string{
		{"editor", "/users", "GET"},
		{"editor", "/users", "POST"},
	}
	if !reflect.DeepEqual(policies, wantPolicies) {
		t.Fatalf("policies = %v, want %v", policies, wantPolicies)
	}
}

func TestSysRoleUpdateDataScopeReplacesDeptAssociations(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysRole{Service: testutil.NewTestService(db)}

	role := models.SysRole{RoleId: 9, RoleName: "editor", RoleKey: "editor", Status: "2"}
	dept1 := models.SysDept{DeptId: 1, DeptName: "Root", DeptPath: "/0/1/", Status: 1}
	dept2 := models.SysDept{DeptId: 2, DeptName: "Child A", DeptPath: "/0/2/", Status: 1}
	dept3 := models.SysDept{DeptId: 3, DeptName: "Child B", DeptPath: "/0/3/", Status: 1}
	for _, seed := range []interface{}{&role, &dept1, &dept2, &dept3} {
		if err := db.Create(seed).Error; err != nil {
			t.Fatalf("seed role dept scope data: %v", err)
		}
	}
	if err := db.Model(&role).Association("SysDept").Append(&dept1); err != nil {
		t.Fatalf("append original dept: %v", err)
	}

	result := service.UpdateDataScope(&servicedto.RoleDataScopeReq{
		RoleId:    9,
		DataScope: "2",
		DeptIds:   []int{2, 3},
	})
	if result.Error != nil {
		t.Fatalf("UpdateDataScope() error = %v", result.Error)
	}

	var updated models.SysRole
	if err := db.Preload("SysDept").First(&updated, 9).Error; err != nil {
		t.Fatalf("reload role: %v", err)
	}
	got := make([]int, 0, len(updated.SysDept))
	for _, dept := range updated.SysDept {
		got = append(got, dept.DeptId)
	}
	if !reflect.DeepEqual(got, []int{2, 3}) {
		t.Fatalf("updated dept ids = %v, want [2 3]", got)
	}
}

func TestSysRoleRemoveDeletesAssociationsAndPolicies(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysRole{Service: testutil.NewTestService(db)}
	enforcer := newRoleTestEnforcer(t)

	role := models.SysRole{RoleId: 9, RoleName: "editor", RoleKey: "editor", Status: "2"}
	menu := models.SysMenu{MenuId: 1, Title: "Users", MenuType: "C"}
	dept := models.SysDept{DeptId: 1, DeptName: "Root", DeptPath: "/0/1/", Status: 1}
	for _, seed := range []interface{}{&role, &menu, &dept} {
		if err := db.Create(seed).Error; err != nil {
			t.Fatalf("seed remove data: %v", err)
		}
	}
	if err := db.Exec("INSERT INTO sys_role_menu (role_id, menu_id) VALUES (?, ?)", 9, 1).Error; err != nil {
		t.Fatalf("seed sys_role_menu: %v", err)
	}
	if err := db.Exec("INSERT INTO sys_role_dept (role_id, dept_id) VALUES (?, ?)", 9, 1).Error; err != nil {
		t.Fatalf("seed sys_role_dept: %v", err)
	}
	if _, err := enforcer.AddNamedPolicy("p", "editor", "/users", "GET"); err != nil {
		t.Fatalf("seed casbin policy: %v", err)
	}

	if err := service.Remove(&servicedto.SysRoleDeleteReq{Ids: []int{9}}, enforcer); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	var count int64
	if err := db.Model(&models.SysRole{}).Where("role_id = ?", 9).Count(&count).Error; err != nil {
		t.Fatalf("count role: %v", err)
	}
	if count != 0 {
		t.Fatalf("role count = %d, want 0", count)
	}
	if hasPolicy, err := enforcer.HasNamedPolicy("p", "editor", "/users", "GET"); err != nil || hasPolicy {
		t.Fatalf("policy remained after remove, hasPolicy=%v err=%v", hasPolicy, err)
	}
}

func TestSysRoleGetPermissionsAndUpdatePolicyBehavior(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysRole{Service: testutil.NewTestService(db)}
	enforcer := newRoleTestEnforcer(t)

	api := models.SysApi{Id: 10, Title: "List", Path: "/users", Action: "GET"}
	menu := models.SysMenu{MenuId: 1, Title: "Users", MenuType: "C", Permission: "sys:user:list"}
	role := models.SysRole{RoleId: 9, RoleName: "editor", RoleKey: "editor", Status: "2"}
	for _, seed := range []interface{}{&api, &menu, &role} {
		if err := db.Create(seed).Error; err != nil {
			t.Fatalf("seed update behavior data: %v", err)
		}
	}
	if err := db.Model(&menu).Association("SysApi").Append(&api); err != nil {
		t.Fatalf("append menu api: %v", err)
	}
	if err := db.Exec("INSERT INTO sys_role_menu (role_id, menu_id) VALUES (?, ?)", 9, 1).Error; err != nil {
		t.Fatalf("seed sys_role_menu: %v", err)
	}
	if _, err := enforcer.AddNamedPolicy("p", "editor", "/old", "GET"); err != nil {
		t.Fatalf("seed old policy: %v", err)
	}

	permissions, err := service.GetById(9)
	if err != nil {
		t.Fatalf("GetById() error = %v", err)
	}
	if !reflect.DeepEqual(permissions, []string{"sys:user:list"}) {
		t.Fatalf("permissions = %v, want [sys:user:list]", permissions)
	}

	err = service.Update(&servicedto.SysRoleUpdateReq{
		RoleId:   9,
		RoleName: "editor",
		RoleKey:  "editor",
		Status:   "2",
		MenuIds:  []int{1},
	}, enforcer)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	policies, err := enforcer.GetFilteredNamedPolicy("p", 0, "editor")
	if err != nil {
		t.Fatalf("GetFilteredNamedPolicy() error = %v", err)
	}
	wantPolicies := [][]string{{"editor", "/users", "GET"}}
	if !reflect.DeepEqual(policies, wantPolicies) {
		t.Fatalf("policies = %v, want %v", policies, wantPolicies)
	}
}
