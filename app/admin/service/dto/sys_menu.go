package dto

import (
	"go-admin-api/app/admin/models"
	common "go-admin-api/common/models"

	"go-admin-api/common/dto"
)

// SysMenuGetPageReq is used for listing or searching.
type SysMenuGetPageReq struct {
	dto.Pagination `search:"-"`
	Title          string `form:"title" search:"type:contains;column:title;table:sys_menu" comment:"Menu name"`          // Menu name
	Visible        int    `form:"visible" search:"type:exact;column:visible;table:sys_menu" comment:"Visibility status"` // Visibility status
}

func (m *SysMenuGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysMenuInsertReq struct {
	MenuId     int             `uri:"id" comment:"ID"`               // ID
	MenuName   string          `form:"menuName" comment:"Menu name"` // Menu name
	Title      string          `form:"title" comment:"Display name"` // Display name
	Icon       string          `form:"icon" comment:"Icon"`          // Icon
	Path       string          `form:"path" comment:"Path"`          // Path
	Paths      string          `form:"paths" comment:"ID path"`      // ID path
	MenuType   string          `form:"menuType" comment:"Menu type"` // Menu type
	SysApi     []models.SysApi `form:"sysApi"`
	Apis       []int           `form:"apis"`
	Action     string          `form:"action" comment:"Request method"`                 // Request method
	Permission string          `form:"permission" comment:"Permission code"`            // Permission code
	ParentId   int             `form:"parentId" comment:"Parent menu"`                  // Parent menu
	NoCache    bool            `form:"noCache" comment:"Whether to cache"`              // Whether to cache
	Breadcrumb string          `form:"breadcrumb" comment:"Whether to show breadcrumb"` // Whether to show breadcrumb
	Component  string          `form:"component" comment:"Component"`                   // Component
	Sort       int             `form:"sort" comment:"Sort order"`                       // Sort order
	Visible    string          `form:"visible" comment:"Whether visible"`               // Whether visible
	IsFrame    string          `form:"isFrame" comment:"Whether frame"`                 // Whether frame
	common.ControlBy
}

func (s *SysMenuInsertReq) Generate(model *models.SysMenu) {
	if s.MenuId != 0 {
		model.MenuId = s.MenuId
	}
	model.MenuName = s.MenuName
	model.Title = s.Title
	model.Icon = s.Icon
	model.Path = s.Path
	model.Paths = s.Paths
	model.MenuType = s.MenuType
	model.Action = s.Action
	model.SysApi = s.SysApi
	model.Permission = s.Permission
	model.ParentId = s.ParentId
	model.NoCache = s.NoCache
	model.Breadcrumb = s.Breadcrumb
	model.Component = s.Component
	model.Sort = s.Sort
	model.Visible = s.Visible
	model.IsFrame = s.IsFrame
	if s.CreateBy != 0 {
		model.CreateBy = s.CreateBy
	}
	if s.UpdateBy != 0 {
		model.UpdateBy = s.UpdateBy
	}
}

func (s *SysMenuInsertReq) GetId() interface{} {
	return s.MenuId
}

type SysMenuUpdateReq struct {
	MenuId     int             `uri:"id" comment:"ID"`               // ID
	MenuName   string          `form:"menuName" comment:"Menu name"` // Menu name
	Title      string          `form:"title" comment:"Display name"` // Display name
	Icon       string          `form:"icon" comment:"Icon"`          // Icon
	Path       string          `form:"path" comment:"Path"`          // Path
	Paths      string          `form:"paths" comment:"ID path"`      // ID path
	MenuType   string          `form:"menuType" comment:"Menu type"` // Menu type
	SysApi     []models.SysApi `form:"sysApi"`
	Apis       []int           `form:"apis"`
	Action     string          `form:"action" comment:"Request method"`                 // Request method
	Permission string          `form:"permission" comment:"Permission code"`            // Permission code
	ParentId   int             `form:"parentId" comment:"Parent menu"`                  // Parent menu
	NoCache    bool            `form:"noCache" comment:"Whether to cache"`              // Whether to cache
	Breadcrumb string          `form:"breadcrumb" comment:"Whether to show breadcrumb"` // Whether to show breadcrumb
	Component  string          `form:"component" comment:"Component"`                   // Component
	Sort       int             `form:"sort" comment:"Sort order"`                       // Sort order
	Visible    string          `form:"visible" comment:"Whether visible"`               // Whether visible
	IsFrame    string          `form:"isFrame" comment:"Whether frame"`                 // Whether frame
	common.ControlBy
}

func (s *SysMenuUpdateReq) Generate(model *models.SysMenu) {
	if s.MenuId != 0 {
		model.MenuId = s.MenuId
	}
	model.MenuName = s.MenuName
	model.Title = s.Title
	model.Icon = s.Icon
	model.Path = s.Path
	model.Paths = s.Paths
	model.MenuType = s.MenuType
	model.Action = s.Action
	model.SysApi = s.SysApi
	model.Permission = s.Permission
	model.ParentId = s.ParentId
	model.NoCache = s.NoCache
	model.Breadcrumb = s.Breadcrumb
	model.Component = s.Component
	model.Sort = s.Sort
	model.Visible = s.Visible
	model.IsFrame = s.IsFrame
	if s.CreateBy != 0 {
		model.CreateBy = s.CreateBy
	}
	if s.UpdateBy != 0 {
		model.UpdateBy = s.UpdateBy
	}
}

func (s *SysMenuUpdateReq) GetId() interface{} {
	return s.MenuId
}

type SysMenuGetReq struct {
	Id int `uri:"id"`
}

func (s *SysMenuGetReq) GetId() interface{} {
	return s.Id
}

type SysMenuDeleteReq struct {
	Ids []int `json:"ids"`
	common.ControlBy
}

func (s *SysMenuDeleteReq) GetId() interface{} {
	return s.Ids
}

type MenuLabel struct {
	Id       int         `json:"id,omitempty" gorm:"-"`
	Label    string      `json:"label,omitempty" gorm:"-"`
	Children []MenuLabel `json:"children,omitempty" gorm:"-"`
}

type MenuRole struct {
	models.SysMenu
	IsSelect bool `json:"is_select" gorm:"-"`
}

type SelectRole struct {
	RoleId int `uri:"roleId"`
}
