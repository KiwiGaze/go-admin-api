package service

import (
	"testing"

	"go-admin-api/app/admin/models"
	servicedto "go-admin-api/app/admin/service/dto"
	"go-admin-api/common/actions"
	"go-admin-api/internal/testutil"
	"github.com/go-admin-team/go-admin-core/sdk/runtime"
)

func TestSysApiGetPageTypeNormalization(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysApi{Service: testutil.NewTestService(db)}

	for _, api := range []models.SysApi{
		{Id: 1, Title: "Blank", Path: "/blank", Action: "GET", Type: ""},
		{Id: 2, Title: "Menu", Path: "/menu", Action: "GET", Type: "menu"},
		{Id: 3, Title: "Blank 2", Path: "/blank-2", Action: "POST", Type: ""},
	} {
		if err := db.Create(&api).Error; err != nil {
			t.Fatalf("seed api: %v", err)
		}
	}

	tests := []struct {
		name string
		req  servicedto.SysApiGetPageReq
		want []int
	}{
		{
			name: "none maps to empty type",
			req:  servicedto.SysApiGetPageReq{Type: "None"},
			want: []int{1, 3},
		},
		{
			name: "legacy chinese empty label maps to empty type",
			req:  servicedto.SysApiGetPageReq{Type: "暂无"},
			want: []int{1, 3},
		},
		{
			name: "non-empty type filters directly",
			req:  servicedto.SysApiGetPageReq{Type: "menu"},
			want: []int{2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var list []models.SysApi
			var count int64
			if err := service.GetPage(&tc.req, &actions.DataPermission{}, &list, &count); err != nil {
				t.Fatalf("GetPage() error = %v", err)
			}
			if int(count) != len(tc.want) {
				t.Fatalf("count = %d, want %d", count, len(tc.want))
			}
			got := make([]int, 0, len(list))
			for _, item := range list {
				got = append(got, item.Id)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("ids = %v, want %v", got, tc.want)
			}
			for idx := range got {
				if got[idx] != tc.want[idx] {
					t.Fatalf("ids = %v, want %v", got, tc.want)
				}
			}
		})
	}
}

func TestSysApiGetMissingRecord(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysApi{Service: testutil.NewTestService(db)}

	var api models.SysApi
	result := service.Get(&servicedto.SysApiGetReq{Id: 404}, &actions.DataPermission{}, &api)
	if result.Error == nil || result.Error.Error() != "the requested object does not exist or cannot be viewed" {
		t.Fatalf("Get() error = %v, want missing object error", result.Error)
	}
}

func TestSysApiRemoveHonorsPermissionScope(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysApi{Service: testutil.NewTestService(db)}

	api := models.SysApi{Id: 1, Title: "Protected", Path: "/protected", Action: "DELETE"}
	api.CreateBy = 2
	if err := db.Create(&api).Error; err != nil {
		t.Fatalf("seed api: %v", err)
	}

	err := service.Remove(&servicedto.SysApiDeleteReq{Ids: []int{1}}, &actions.DataPermission{
		DataScope: "5",
		UserId:    1,
	})
	if err == nil || err.Error() != "no permission to delete this data" {
		t.Fatalf("Remove() error = %v, want permission error", err)
	}
}

func TestSysApiCheckStorageSysApiIsIdempotent(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysApi{Service: testutil.NewTestService(db)}

	routers := []runtime.Router{
		{RelativePath: "/api/v1/test", HttpMethod: "GET", Handler: "handler.One"},
		{RelativePath: "/api/v1/test", HttpMethod: "GET", Handler: "handler.Two"},
	}
	if err := service.CheckStorageSysApi(&routers); err != nil {
		t.Fatalf("CheckStorageSysApi() error = %v", err)
	}
	if err := service.CheckStorageSysApi(&routers); err != nil {
		t.Fatalf("CheckStorageSysApi() second call error = %v", err)
	}

	var count int64
	if err := db.Model(&models.SysApi{}).Count(&count).Error; err != nil {
		t.Fatalf("count apis: %v", err)
	}
	if count != 1 {
		t.Fatalf("api count = %d, want 1", count)
	}
}

func TestSysApiUpdateCurrentBehaviorIgnoresPermissionScopeOnLoad(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysApi{Service: testutil.NewTestService(db)}

	api := models.SysApi{Id: 1, Title: "Original", Path: "/resource", Action: "GET"}
	api.CreateBy = 2
	if err := db.Create(&api).Error; err != nil {
		t.Fatalf("seed api: %v", err)
	}

	err := service.Update(&servicedto.SysApiUpdateReq{
		Id:     1,
		Title:  "Updated",
		Path:   "/resource",
		Action: "GET",
	}, &actions.DataPermission{
		DataScope: "5",
		UserId:    1,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	var updated models.SysApi
	if err := db.First(&updated, 1).Error; err != nil {
		t.Fatalf("reload api: %v", err)
	}
	if updated.Title != "Updated" {
		t.Fatalf("updated title = %q, want Updated", updated.Title)
	}
}
