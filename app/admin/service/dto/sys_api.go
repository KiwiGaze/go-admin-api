package dto

import (
	"go-admin-api/app/admin/models"
	"go-admin-api/common/dto"
	common "go-admin-api/common/models"
)

// SysApiGetPageReq feature list request parameters
type SysApiGetPageReq struct {
	dto.Pagination `search:"-"`
	Title          string `form:"title"  search:"type:contains;column:title;table:sys_api" comment:"Title"`
	Path           string `form:"path"  search:"type:contains;column:path;table:sys_api" comment:"Path"`
	Action         string `form:"action"  search:"type:exact;column:action;table:sys_api" comment:"Request method"`
	ParentId       string `form:"parentId"  search:"type:exact;column:parent_id;table:sys_api" comment:"Button ID"`
	Type           string `form:"type" search:"-" comment:"Type"`
	SysApiOrder
}

type SysApiOrder struct {
	TitleOrder     string `search:"type:order;column:title;table:sys_api" form:"titleOrder"`
	PathOrder      string `search:"type:order;column:path;table:sys_api" form:"pathOrder"`
	CreatedAtOrder string `search:"type:order;column:created_at;table:sys_api" form:"createdAtOrder"`
}

func (m *SysApiGetPageReq) GetNeedSearch() interface{} {
	return *m
}

// SysApiInsertReq feature creation request parameters
type SysApiInsertReq struct {
	Id     int    `json:"-" comment:"Code"` // Code
	Handle string `json:"handle" comment:"handle"`
	Title  string `json:"title" comment:"Title"`
	Path   string `json:"path" comment:"Path"`
	Type   string `json:"type" comment:""`
	Action string `json:"action" comment:"Type"`
	common.ControlBy
}

func (s *SysApiInsertReq) Generate(model *models.SysApi) {
	model.Handle = s.Handle
	model.Title = s.Title
	model.Path = s.Path
	model.Type = s.Type
	model.Action = s.Action
}

func (s *SysApiInsertReq) GetId() interface{} {
	return s.Id
}

// SysApiUpdateReq feature update request parameters
type SysApiUpdateReq struct {
	Id     int    `uri:"id" comment:"Code"` // Code
	Handle string `json:"handle" comment:"handle"`
	Title  string `json:"title" comment:"Title"`
	Path   string `json:"path" comment:"Path"`
	Type   string `json:"type" comment:""`
	Action string `json:"action" comment:"Type"`
	common.ControlBy
}

func (s *SysApiUpdateReq) Generate(model *models.SysApi) {
	if s.Id != 0 {
		model.Id = s.Id
	}
	model.Handle = s.Handle
	model.Title = s.Title
	model.Path = s.Path
	model.Type = s.Type
	model.Action = s.Action
}

func (s *SysApiUpdateReq) GetId() interface{} {
	return s.Id
}

// SysApiGetReq feature retrieval request parameters
type SysApiGetReq struct {
	Id int `uri:"id"`
}

func (s *SysApiGetReq) GetId() interface{} {
	return s.Id
}

// SysApiDeleteReq feature deletion request parameters
type SysApiDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *SysApiDeleteReq) GetId() interface{} {
	return s.Ids
}
