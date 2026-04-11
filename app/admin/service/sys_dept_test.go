package service

import (
	"reflect"
	"testing"

	"go-admin-api/app/admin/models"
	servicedto "go-admin-api/app/admin/service/dto"
	"go-admin-api/internal/testutil"
)

func TestSysDeptInsertBuildsPaths(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysDept{Service: testutil.NewTestService(db)}

	if err := service.Insert(&servicedto.SysDeptInsertReq{DeptId: 1, DeptName: "Root", Status: 1}); err != nil {
		t.Fatalf("Insert(root) error = %v", err)
	}
	if err := service.Insert(&servicedto.SysDeptInsertReq{DeptId: 2, ParentId: 1, DeptName: "Child", Status: 1}); err != nil {
		t.Fatalf("Insert(child) error = %v", err)
	}

	var root models.SysDept
	var child models.SysDept
	if err := db.First(&root, 1).Error; err != nil {
		t.Fatalf("reload root: %v", err)
	}
	if err := db.First(&child, 2).Error; err != nil {
		t.Fatalf("reload child: %v", err)
	}
	if root.DeptPath != "/0/1/" {
		t.Fatalf("root dept path = %q, want /0/1/", root.DeptPath)
	}
	if child.DeptPath != "/0/1/2/" {
		t.Fatalf("child dept path = %q, want /0/1/2/", child.DeptPath)
	}
}

func TestSysDeptUpdateRecomputesPath(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysDept{Service: testutil.NewTestService(db)}

	for _, dept := range []models.SysDept{
		{DeptId: 1, DeptName: "Root", DeptPath: "/0/1/", Status: 1},
		{DeptId: 2, DeptName: "Child", ParentId: 1, DeptPath: "/0/1/2/", Status: 1},
		{DeptId: 3, DeptName: "New Root", DeptPath: "/0/3/", Status: 1},
	} {
		if err := db.Create(&dept).Error; err != nil {
			t.Fatalf("seed dept: %v", err)
		}
	}

	if err := service.Update(&servicedto.SysDeptUpdateReq{DeptId: 2, ParentId: 3, DeptName: "Child", Status: 1}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	var updated models.SysDept
	if err := db.First(&updated, 2).Error; err != nil {
		t.Fatalf("reload dept: %v", err)
	}
	if updated.DeptPath != "/0/3/2/" {
		t.Fatalf("updated dept path = %q, want /0/3/2/", updated.DeptPath)
	}
}

func TestSysDeptTreeBuilders(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysDept{Service: testutil.NewTestService(db)}

	for _, dept := range []models.SysDept{
		{DeptId: 1, DeptName: "Root", DeptPath: "/0/1/", Status: 1},
		{DeptId: 2, DeptName: "Child A", ParentId: 1, DeptPath: "/0/1/2/", Status: 1},
		{DeptId: 3, DeptName: "Child B", ParentId: 1, DeptPath: "/0/1/3/", Status: 1},
	} {
		if err := db.Create(&dept).Error; err != nil {
			t.Fatalf("seed dept tree: %v", err)
		}
	}

	tree, err := service.SetDeptTree(&servicedto.SysDeptGetPageReq{})
	if err != nil {
		t.Fatalf("SetDeptTree() error = %v", err)
	}
	if len(tree) != 1 || len(tree[0].Children) != 2 {
		t.Fatalf("dept tree = %+v, want root with two children", tree)
	}

	page, err := service.SetDeptPage(&servicedto.SysDeptGetPageReq{})
	if err != nil {
		t.Fatalf("SetDeptPage() error = %v", err)
	}
	if len(page) != 1 || len(page[0].Children) != 2 {
		t.Fatalf("dept page tree = %+v, want root with two children", page)
	}

	labels, err := service.SetDeptLabel()
	if err != nil {
		t.Fatalf("SetDeptLabel() error = %v", err)
	}
	if len(labels) != 1 || len(labels[0].Children) != 2 {
		t.Fatalf("dept labels = %+v, want root with two children", labels)
	}
}

func TestSysDeptGetWithRoleIdReturnsLeafDepartments(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysDept{Service: testutil.NewTestService(db)}

	role := models.SysRole{RoleId: 9, RoleName: "editor", RoleKey: "editor", Status: "2"}
	root := models.SysDept{DeptId: 1, DeptName: "Root", DeptPath: "/0/1/", Status: 1}
	child := models.SysDept{DeptId: 2, ParentId: 1, DeptName: "Child", DeptPath: "/0/1/2/", Status: 1}
	leaf := models.SysDept{DeptId: 3, ParentId: 2, DeptName: "Leaf", DeptPath: "/0/1/2/3/", Status: 1}
	for _, seed := range []interface{}{&role, &root, &child, &leaf} {
		if err := db.Create(seed).Error; err != nil {
			t.Fatalf("seed role dept data: %v", err)
		}
	}
	if err := db.Model(&role).Association("SysDept").Append(&root, &child, &leaf); err != nil {
		t.Fatalf("append role depts: %v", err)
	}

	deptIDs, err := service.GetWithRoleId(9)
	if err != nil {
		t.Fatalf("GetWithRoleId() error = %v", err)
	}
	if !reflect.DeepEqual(deptIDs, []int{3}) {
		t.Fatalf("leaf dept ids = %v, want [3]", deptIDs)
	}
}
