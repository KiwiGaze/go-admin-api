package dto

import (
	"reflect"
	"testing"

	"go-admin-api/common/global"
	"go-admin-api/internal/testutil"
)

type searchTestItem struct {
	ID     int    `gorm:"primaryKey"`
	Name   string
	Status string
	DeptID int
}

func (searchTestItem) TableName() string {
	return "search_test_items"
}

type searchTestDept struct {
	DeptID   int    `gorm:"primaryKey"`
	DeptPath string
}

func (searchTestDept) TableName() string {
	return "search_test_depts"
}

type SearchTestDeptJoin struct {
	DeptID string `search:"type:contains;column:dept_path;table:search_test_depts"`
}

type searchTestReq struct {
	Name      string `search:"type:contains;column:name;table:search_test_items"`
	Status    string `search:"type:exact;column:status;table:search_test_items"`
	NameOrder string `search:"type:order;column:name;table:search_test_items"`
	SearchTestDeptJoin `search:"type:left;on:dept_id:dept_id;table:search_test_items;join:search_test_depts"`
}

func TestGeneralDelDtoGetIds(t *testing.T) {
	t.Run("returns single id and positive batch ids", func(t *testing.T) {
		got := GeneralDelDto{Id: 1, Ids: []int{2, -1, 0, 3}}.GetIds()
		want := []int{1, 2, 3}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("GetIds() = %v, want %v", got, want)
		}
	})

	t.Run("returns sentinel zero when no valid ids exist", func(t *testing.T) {
		got := GeneralDelDto{Ids: []int{0, -1}}.GetIds()
		want := []int{0}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("GetIds() = %v, want %v", got, want)
		}
	})
}

func TestPaginate(t *testing.T) {
	db := testutil.NewTestDB(t, &searchTestItem{})
	for _, item := range []searchTestItem{
		{ID: 1, Name: "alpha"},
		{ID: 2, Name: "beta"},
		{ID: 3, Name: "gamma"},
	} {
		if err := db.Create(&item).Error; err != nil {
			t.Fatalf("seed item: %v", err)
		}
	}

	t.Run("applies limit and offset", func(t *testing.T) {
		var items []searchTestItem
		if err := db.Model(&searchTestItem{}).
			Order("id").
			Scopes(Paginate(2, 2)).
			Find(&items).Error; err != nil {
			t.Fatalf("paginate query: %v", err)
		}
		want := []int{3}
		got := make([]int, 0, len(items))
		for _, item := range items {
			got = append(got, item.ID)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("paginate ids = %v, want %v", got, want)
		}
	})

	t.Run("clamps negative offset to zero", func(t *testing.T) {
		var items []searchTestItem
		if err := db.Model(&searchTestItem{}).
			Order("id").
			Scopes(Paginate(2, 0)).
			Find(&items).Error; err != nil {
			t.Fatalf("paginate query: %v", err)
		}
		got := []int{items[0].ID, items[1].ID}
		want := []int{1, 2}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("paginate ids = %v, want %v", got, want)
		}
	})
}

func TestMakeCondition(t *testing.T) {
	db := testutil.NewTestDB(t, &searchTestItem{}, &searchTestDept{})
	for _, dept := range []searchTestDept{
		{DeptID: 10, DeptPath: "/0/10/"},
		{DeptID: 20, DeptPath: "/0/20/"},
	} {
		if err := db.Create(&dept).Error; err != nil {
			t.Fatalf("seed dept: %v", err)
		}
	}
	for _, item := range []searchTestItem{
		{ID: 1, Name: "alpha", Status: "enabled", DeptID: 10},
		{ID: 2, Name: "alphabet", Status: "enabled", DeptID: 20},
		{ID: 3, Name: "beta", Status: "disabled", DeptID: 10},
	} {
		if err := db.Create(&item).Error; err != nil {
			t.Fatalf("seed item: %v", err)
		}
	}

	originalDriver := global.Driver
	t.Cleanup(func() {
		global.Driver = originalDriver
	})

	t.Run("mysql driver applies exact, contains, join, and order filters", func(t *testing.T) {
		global.Driver = "mysql"
		req := searchTestReq{
			Name:      "alph",
			Status:    "enabled",
			NameOrder: "desc",
			SearchTestDeptJoin: SearchTestDeptJoin{
				DeptID: "/0/20/",
			},
		}

		var items []searchTestItem
		if err := db.Model(&searchTestItem{}).
			Scopes(MakeCondition(req)).
			Find(&items).Error; err != nil {
			t.Fatalf("query with make condition: %v", err)
		}
		if len(items) != 1 || items[0].ID != 2 {
			t.Fatalf("filtered items = %+v, want item 2 only", items)
		}
	})

	t.Run("postgres driver branch remains queryable", func(t *testing.T) {
		global.Driver = "postgres"
		req := searchTestReq{
			Status:    "enabled",
			NameOrder: "asc",
		}

		var items []searchTestItem
		if err := db.Model(&searchTestItem{}).
			Scopes(MakeCondition(req)).
			Find(&items).Error; err != nil {
			t.Fatalf("query with make condition: %v", err)
		}
		got := []int{items[0].ID, items[1].ID}
		want := []int{1, 2}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ordered ids = %v, want %v", got, want)
		}
	})
}
