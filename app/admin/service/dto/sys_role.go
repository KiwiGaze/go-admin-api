package dto

import (
	"go-admin-api/app/admin/models"
	common "go-admin-api/common/models"

	"go-admin-api/common/dto"
)

type SysRoleGetPageReq struct {
	dto.Pagination `search:"-"`

	RoleId    int    `form:"roleId" search:"type:exact;column:role_id;table:sys_role" comment:"Role ID"`       // Role ID
	RoleName  string `form:"roleName" search:"type:exact;column:role_name;table:sys_role" comment:"Role name"`   // Role name
	Status    string `form:"status" search:"type:exact;column:status;table:sys_role" comment:"Status"`          // Status
	RoleKey   string `form:"roleKey" search:"type:exact;column:role_key;table:sys_role" comment:"Role key"`      // Role key
	RoleSort  int    `form:"roleSort" search:"type:exact;column:role_sort;table:sys_role" comment:"Role sort"`   // Role sort
	Flag      string `form:"flag" search:"type:exact;column:flag;table:sys_role" comment:"Flag"`                // Flag
	Remark    string `form:"remark" search:"type:exact;column:remark;table:sys_role" comment:"Remark"`          // Remark
	Admin     bool   `form:"admin" search:"type:exact;column:admin;table:sys_role" comment:"Is admin"`
	DataScope string `form:"dataScope" search:"type:exact;column:data_scope;table:sys_role" comment:"Data scope"`
}

type SysRoleOrder struct {
	RoleIdOrder    string `search:"type:order;column:role_id;table:sys_role" form:"roleIdOrder"`
	RoleNameOrder  string `search:"type:order;column:role_name;table:sys_role" form:"roleNameOrder"`
	RoleSortOrder  string `search:"type:order;column:role_sort;table:sys_role" form:"roleSortOrder"`
	StatusOrder    string `search:"type:order;column:status;table:sys_role" form:"statusOrder"`
	CreatedAtOrder string `search:"type:order;column:created_at;table:sys_role" form:"createdAtOrder"`
}

func (m *SysRoleGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysRoleInsertReq struct {
	RoleId    int              `uri:"id" comment:"Role ID"`      // Role ID
	RoleName  string           `form:"roleName" comment:"Role name"` // Role name
	Status    string           `form:"status" comment:"Status"`  // Status: 1 disabled, 2 enabled
	RoleKey   string           `form:"roleKey" comment:"Role key"` // Role key
	RoleSort  int              `form:"roleSort" comment:"Role sort"` // Role sort
	Flag      string           `form:"flag" comment:"Flag"`      // Flag
	Remark    string           `form:"remark" comment:"Remark"`  // Remark
	Admin     bool             `form:"admin" comment:"Is admin"`
	DataScope string           `form:"dataScope"`
	SysMenu   []models.SysMenu `form:"sysMenu"`
	MenuIds   []int            `form:"menuIds"`
	SysDept   []models.SysDept `form:"sysDept"`
	DeptIds   []int            `form:"deptIds"`
	common.ControlBy
}

func (s *SysRoleInsertReq) Generate(model *models.SysRole) {
	if s.RoleId != 0 {
		model.RoleId = s.RoleId
	}
	model.RoleName = s.RoleName
	model.Status = s.Status
	model.RoleKey = s.RoleKey
	model.RoleSort = s.RoleSort
	model.Flag = s.Flag
	model.Remark = s.Remark
	model.Admin = s.Admin
	model.DataScope = s.DataScope
	model.SysMenu = &s.SysMenu
	model.SysDept = s.SysDept
}

func (s *SysRoleInsertReq) GetId() interface{} {
	return s.RoleId
}

type SysRoleUpdateReq struct {
	RoleId    int              `uri:"id" comment:"Role ID"`      // Role ID
	RoleName  string           `form:"roleName" comment:"Role name"` // Role name
	Status    string           `form:"status" comment:"Status"`  // Status
	RoleKey   string           `form:"roleKey" comment:"Role key"` // Role key
	RoleSort  int              `form:"roleSort" comment:"Role sort"` // Role sort
	Flag      string           `form:"flag" comment:"Flag"`      // Flag
	Remark    string           `form:"remark" comment:"Remark"`  // Remark
	Admin     bool             `form:"admin" comment:"Is admin"`
	DataScope string           `form:"dataScope"`
	SysMenu   []models.SysMenu `form:"sysMenu"`
	MenuIds   []int            `form:"menuIds"`
	SysDept   []models.SysDept `form:"sysDept"`
	DeptIds   []int            `form:"deptIds"`
	common.ControlBy
}

func (s *SysRoleUpdateReq) Generate(model *models.SysRole) {
	if s.RoleId != 0 {
		model.RoleId = s.RoleId
	}
	model.RoleName = s.RoleName
	model.Status = s.Status
	model.RoleKey = s.RoleKey
	model.RoleSort = s.RoleSort
	model.Flag = s.Flag
	model.Remark = s.Remark
	model.Admin = s.Admin
	model.DataScope = s.DataScope
	model.SysMenu = &s.SysMenu
	model.SysDept = s.SysDept
}

func (s *SysRoleUpdateReq) GetId() interface{} {
	return s.RoleId
}

type UpdateStatusReq struct {
	RoleId int    `form:"roleId" comment:"Role ID"` // Role ID
	Status string `form:"status" comment:"Status"`  // Status
	common.ControlBy
}

func (s *UpdateStatusReq) Generate(model *models.SysRole) {
	if s.RoleId != 0 {
		model.RoleId = s.RoleId
	}
	model.Status = s.Status
}

func (s *UpdateStatusReq) GetId() interface{} {
	return s.RoleId
}

type SysRoleByName struct {
	RoleName string `form:"role"` // Role name
}

type SysRoleGetReq struct {
	Id int `uri:"id"`
}

func (s *SysRoleGetReq) GetId() interface{} {
	return s.Id
}

type SysRoleDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *SysRoleDeleteReq) GetId() interface{} {
	return s.Ids
}

// RoleDataScopeReq role data permission update
type RoleDataScopeReq struct {
	RoleId    int    `json:"roleId" binding:"required"`
	DataScope string `json:"dataScope" binding:"required"`
	DeptIds   []int  `json:"deptIds"`
}

func (s *RoleDataScopeReq) Generate(model *models.SysRole) {
	if s.RoleId != 0 {
		model.RoleId = s.RoleId
	}
	model.DataScope = s.DataScope
	model.DeptIds = s.DeptIds
}

type DeptIdList struct {
	DeptId int `json:"DeptId"`
}
