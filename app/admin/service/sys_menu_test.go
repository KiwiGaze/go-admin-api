package service

import (
	"reflect"
	"testing"

	"go-admin-api/app/admin/models"
	servicedto "go-admin-api/app/admin/service/dto"
	commodels "go-admin-api/common/models"
	"go-admin-api/internal/testutil"
)

func TestSysMenuInitPaths(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysMenu{Service: testutil.NewTestService(db)}

	t.Run("root menu gets /0 prefixed path", func(t *testing.T) {
		menu := models.SysMenu{MenuId: 1, MenuName: "root", Title: "Root", MenuType: commodels.Directory}
		if err := db.Create(&menu).Error; err != nil {
			t.Fatalf("seed root menu: %v", err)
		}
		if err := service.initPaths(db, &menu); err != nil {
			t.Fatalf("initPaths() error = %v", err)
		}
		if menu.Paths != "/0/1" {
			t.Fatalf("menu paths = %q, want /0/1", menu.Paths)
		}
	})

	t.Run("child menu appends to parent path", func(t *testing.T) {
		parent := models.SysMenu{MenuId: 10, MenuName: "parent", Title: "Parent", MenuType: commodels.Directory, Paths: "/0/10"}
		child := models.SysMenu{MenuId: 11, MenuName: "child", Title: "Child", MenuType: commodels.Menu, ParentId: 10}
		if err := db.Create(&parent).Error; err != nil {
			t.Fatalf("seed parent: %v", err)
		}
		if err := db.Create(&child).Error; err != nil {
			t.Fatalf("seed child: %v", err)
		}
		if err := service.initPaths(db, &child); err != nil {
			t.Fatalf("initPaths() error = %v", err)
		}
		if child.Paths != "/0/10/11" {
			t.Fatalf("child paths = %q, want /0/10/11", child.Paths)
		}
	})

	t.Run("invalid parent path returns explicit error", func(t *testing.T) {
		parent := models.SysMenu{MenuId: 20, MenuName: "broken", Title: "Broken", MenuType: commodels.Directory}
		child := models.SysMenu{MenuId: 21, MenuName: "child", Title: "Child", MenuType: commodels.Menu, ParentId: 20}
		if err := db.Create(&parent).Error; err != nil {
			t.Fatalf("seed broken parent: %v", err)
		}
		if err := db.Create(&child).Error; err != nil {
			t.Fatalf("seed child: %v", err)
		}
		err := service.initPaths(db, &child)
		if err == nil || err.Error() != "parent paths are invalid, please try updating the parent menu of the current node" {
			t.Fatalf("initPaths() error = %v, want invalid parent paths error", err)
		}
	})
}

func TestSysMenuUpdateCascadesPaths(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysMenu{Service: testutil.NewTestService(db)}

	parent := models.SysMenu{MenuId: 1, MenuName: "root", Title: "Root", MenuType: commodels.Directory, Paths: "/0/1"}
	child := models.SysMenu{MenuId: 2, MenuName: "child", Title: "Child", MenuType: commodels.Menu, ParentId: 1, Paths: "/0/1/2"}
	grandchild := models.SysMenu{MenuId: 3, MenuName: "button", Title: "Button", MenuType: commodels.Button, ParentId: 2, Paths: "/0/1/2/3"}

	for _, menu := range []models.SysMenu{parent, child, grandchild} {
		if err := db.Create(&menu).Error; err != nil {
			t.Fatalf("seed menu: %v", err)
		}
	}

	result := service.Update(&servicedto.SysMenuUpdateReq{
		MenuId:   2,
		MenuName: "child",
		Title:    "Child",
		MenuType: commodels.Menu,
		ParentId: 1,
		Paths:    "/0/9/2",
	})
	if result.Error != nil {
		t.Fatalf("Update() error = %v", result.Error)
	}

	var updatedChild models.SysMenu
	var updatedGrandchild models.SysMenu
	if err := db.First(&updatedChild, 2).Error; err != nil {
		t.Fatalf("reload child: %v", err)
	}
	if err := db.First(&updatedGrandchild, 3).Error; err != nil {
		t.Fatalf("reload grandchild: %v", err)
	}
	if updatedChild.Paths != "/0/9/2" {
		t.Fatalf("child paths = %q, want /0/9/2", updatedChild.Paths)
	}
	if updatedGrandchild.Paths != "/0/9/2/3" {
		t.Fatalf("grandchild paths = %q, want /0/9/2/3", updatedGrandchild.Paths)
	}
}

func TestSysMenuGetCopiesApiIDs(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysMenu{Service: testutil.NewTestService(db)}

	menu := models.SysMenu{MenuId: 1, MenuName: "users", Title: "Users", MenuType: commodels.Menu}
	api1 := models.SysApi{Id: 10, Title: "List", Path: "/users", Action: "GET"}
	api2 := models.SysApi{Id: 11, Title: "Create", Path: "/users", Action: "POST"}
	for _, seed := range []interface{}{&menu, &api1, &api2} {
		if err := db.Create(seed).Error; err != nil {
			t.Fatalf("seed menu relation: %v", err)
		}
	}
	if err := db.Model(&menu).Association("SysApi").Append(&api1, &api2); err != nil {
		t.Fatalf("append menu apis: %v", err)
	}

	var loaded models.SysMenu
	result := service.Get(&servicedto.SysMenuGetReq{Id: 1}, &loaded)
	if result.Error != nil {
		t.Fatalf("Get() error = %v", result.Error)
	}
	if !reflect.DeepEqual(loaded.Apis, []int{10, 11}) {
		t.Fatalf("loaded api ids = %v, want [10 11]", loaded.Apis)
	}
}

func TestSysMenuTreeHelpers(t *testing.T) {
	menus := []models.SysMenu{
		{MenuId: 1, Title: "Root", MenuType: commodels.Directory},
		{MenuId: 2, ParentId: 1, Title: "Page", MenuType: commodels.Menu},
		{MenuId: 3, ParentId: 2, Title: "Button", MenuType: commodels.Button},
	}

	tree := menuCall(&menus, menus[0])
	if len(tree.Children) != 1 || tree.Children[0].MenuId != 2 {
		t.Fatalf("menuCall tree = %+v, want root child page", tree)
	}
	if len(tree.Children[0].Children) != 1 || tree.Children[0].Children[0].MenuId != 3 {
		t.Fatalf("menuCall children = %+v, want button leaf", tree.Children[0].Children)
	}

	labelTree := menuLabelCall(&menus, servicedto.MenuLabel{Id: 1, Label: "Root"})
	if len(labelTree.Children) != 1 || labelTree.Children[0].Id != 2 {
		t.Fatalf("menuLabelCall tree = %+v, want page child", labelTree)
	}
	if len(labelTree.Children[0].Children) != 1 || labelTree.Children[0].Children[0].Id != 3 {
		t.Fatalf("menuLabelCall children = %+v, want button leaf", labelTree.Children[0].Children)
	}

	leafMenus := []models.SysMenu{
		{MenuId: 10, Title: "Root", MenuType: commodels.Directory},
		{MenuId: 11, ParentId: 10, Title: "Leaf Page", MenuType: commodels.Menu},
	}
	leafLabel := menuLabelCall(&leafMenus, servicedto.MenuLabel{Id: 10, Label: "Root"})
	if leafLabel.Children[0].Children != nil {
		t.Fatalf("leaf label children = %+v, want nil for non-button leaf", leafLabel.Children[0].Children)
	}
}

func TestSysMenuGetSysMenuByRoleName(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysMenu{Service: testutil.NewTestService(db)}

	root := models.SysMenu{MenuId: 1, Title: "Root", MenuType: commodels.Directory, Sort: 1}
	page := models.SysMenu{MenuId: 2, Title: "Page", MenuType: commodels.Menu, Sort: 2}
	button := models.SysMenu{MenuId: 3, Title: "Button", MenuType: commodels.Button, Sort: 3}
	role := models.SysRole{RoleId: 9, RoleName: "editor", RoleKey: "editor", Status: "2"}
	for _, seed := range []interface{}{&root, &page, &button, &role} {
		if err := db.Create(seed).Error; err != nil {
			t.Fatalf("seed role menu data: %v", err)
		}
	}
	for _, pair := range []map[string]int{
		{"role_id": 9, "menu_id": 1},
		{"role_id": 9, "menu_id": 2},
		{"role_id": 9, "menu_id": 3},
	} {
		if err := db.Exec("INSERT INTO sys_role_menu (role_id, menu_id) VALUES (?, ?)", pair["role_id"], pair["menu_id"]).Error; err != nil {
			t.Fatalf("seed sys_role_menu: %v", err)
		}
	}

	adminMenus, err := service.GetSysMenuByRoleName("admin")
	if err != nil {
		t.Fatalf("GetSysMenuByRoleName(admin) error = %v", err)
	}
	if got := []int{adminMenus[0].MenuId, adminMenus[1].MenuId}; !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("admin menus = %v, want [1 2]", got)
	}

	editorMenus, err := service.GetSysMenuByRoleName("editor")
	if err != nil {
		t.Fatalf("GetSysMenuByRoleName(editor) error = %v", err)
	}
	if got := []int{editorMenus[0].MenuId, editorMenus[1].MenuId}; !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("editor menus = %v, want [1 2]", got)
	}
}
