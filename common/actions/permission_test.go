package actions

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"go-admin-api/internal/testutil"
)

type permissionUser struct {
	UserID int `gorm:"column:user_id;primaryKey"`
	DeptID int `gorm:"column:dept_id"`
	RoleID int `gorm:"column:role_id"`
}

func (permissionUser) TableName() string {
	return "sys_user"
}

type permissionRole struct {
	RoleID    int    `gorm:"column:role_id;primaryKey"`
	DataScope string `gorm:"column:data_scope"`
}

func (permissionRole) TableName() string {
	return "sys_role"
}

type permissionDept struct {
	DeptID   int    `gorm:"column:dept_id;primaryKey"`
	DeptPath string `gorm:"column:dept_path"`
}

func (permissionDept) TableName() string {
	return "sys_dept"
}

type permissionRoleDept struct {
	RoleID int `gorm:"column:role_id"`
	DeptID int `gorm:"column:dept_id"`
}

func (permissionRoleDept) TableName() string {
	return "sys_role_dept"
}

type permissionRecord struct {
	ID       int `gorm:"primaryKey"`
	CreateBy int
}

func (permissionRecord) TableName() string {
	return "permission_records"
}

func TestNewDataPermission(t *testing.T) {
	db := testutil.NewTestDB(t, &permissionUser{}, &permissionRole{})
	if err := db.Create(&permissionRole{RoleID: 7, DataScope: "4"}).Error; err != nil {
		t.Fatalf("seed role: %v", err)
	}
	if err := db.Create(&permissionUser{UserID: 11, DeptID: 22, RoleID: 7}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	permission, err := newDataPermission(db, 11)
	if err != nil {
		t.Fatalf("newDataPermission() error = %v", err)
	}

	want := &DataPermission{UserId: 11, DeptId: 22, RoleId: 7, DataScope: "4"}
	if !reflect.DeepEqual(permission, want) {
		t.Fatalf("newDataPermission() = %+v, want %+v", permission, want)
	}
}

func TestPermission(t *testing.T) {
	originalEnableDP := config.ApplicationConfig.EnableDP
	t.Cleanup(func() {
		config.ApplicationConfig.EnableDP = originalEnableDP
	})

	config.ApplicationConfig.EnableDP = true

	baseDB := testutil.NewTestDB(t, &permissionUser{}, &permissionRole{}, &permissionDept{}, &permissionRoleDept{}, &permissionRecord{})
	for _, user := range []permissionUser{
		{UserID: 1, DeptID: 10, RoleID: 100},
		{UserID: 2, DeptID: 20, RoleID: 200},
		{UserID: 3, DeptID: 11, RoleID: 300},
		{UserID: 4, DeptID: 10, RoleID: 400},
	} {
		if err := baseDB.Create(&user).Error; err != nil {
			t.Fatalf("seed user: %v", err)
		}
	}
	for _, dept := range []permissionDept{
		{DeptID: 10, DeptPath: "/0/10/"},
		{DeptID: 11, DeptPath: "/0/10/11/"},
		{DeptID: 20, DeptPath: "/0/20/"},
	} {
		if err := baseDB.Create(&dept).Error; err != nil {
			t.Fatalf("seed dept: %v", err)
		}
	}
	if err := baseDB.Create(&permissionRoleDept{RoleID: 100, DeptID: 20}).Error; err != nil {
		t.Fatalf("seed role dept: %v", err)
	}
	for _, record := range []permissionRecord{
		{ID: 1, CreateBy: 1},
		{ID: 2, CreateBy: 2},
		{ID: 3, CreateBy: 3},
		{ID: 4, CreateBy: 4},
	} {
		if err := baseDB.Create(&record).Error; err != nil {
			t.Fatalf("seed permission record: %v", err)
		}
	}

	tests := []struct {
		name string
		p    *DataPermission
		want []int
	}{
		{
			name: "default scope leaves records unfiltered",
			p:    &DataPermission{DataScope: "1", UserId: 1, DeptId: 10, RoleId: 100},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "custom dept scope includes assigned department users",
			p:    &DataPermission{DataScope: "2", UserId: 1, DeptId: 10, RoleId: 100},
			want: []int{2},
		},
		{
			name: "own dept scope includes users in same department",
			p:    &DataPermission{DataScope: "3", UserId: 1, DeptId: 10, RoleId: 100},
			want: []int{1, 4},
		},
		{
			name: "dept and children scope includes subtree users",
			p:    &DataPermission{DataScope: "4", UserId: 1, DeptId: 10, RoleId: 100},
			want: []int{1, 3, 4},
		},
		{
			name: "self scope includes current user only",
			p:    &DataPermission{DataScope: "5", UserId: 1, DeptId: 10, RoleId: 100},
			want: []int{1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var records []permissionRecord
			if err := baseDB.Model(&permissionRecord{}).
				Scopes(Permission("permission_records", tc.p)).
				Order("id").
				Find(&records).Error; err != nil {
				t.Fatalf("permission query: %v", err)
			}
			got := make([]int, 0, len(records))
			for _, record := range records {
				got = append(got, record.ID)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("filtered ids = %v, want %v", got, tc.want)
			}
		})
	}

	t.Run("disabled data permission bypasses filtering", func(t *testing.T) {
		config.ApplicationConfig.EnableDP = false
		defer func() {
			config.ApplicationConfig.EnableDP = true
		}()

		var records []permissionRecord
		if err := baseDB.Model(&permissionRecord{}).
			Scopes(Permission("permission_records", &DataPermission{DataScope: "5", UserId: 1})).
			Order("id").
			Find(&records).Error; err != nil {
			t.Fatalf("permission query: %v", err)
		}
		if len(records) != 4 {
			t.Fatalf("record count = %d, want 4", len(records))
		}
	})
}

func TestGetPermissionFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}

	if got := GetPermissionFromContext(ctx); got == nil || *got != (DataPermission{}) {
		t.Fatalf("missing permission = %+v, want zero-value permission", got)
	}

	want := &DataPermission{DataScope: "5", UserId: 7}
	ctx.Set(PermissionKey, want)
	got := GetPermissionFromContext(ctx)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetPermissionFromContext() = %+v, want %+v", got, want)
	}
}
