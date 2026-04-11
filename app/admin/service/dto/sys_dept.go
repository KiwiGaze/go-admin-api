package dto

import (
	"go-admin-api/app/admin/models"
	common "go-admin-api/common/models"
)

// SysDeptGetPageReq Struct used for listing or searching
type SysDeptGetPageReq struct {
	DeptId         int    `form:"deptId" search:"type:exact;column:dept_id;table:sys_dept" comment:"id"`              // id
	ParentId       int    `form:"parentId" search:"type:exact;column:parent_id;table:sys_dept" comment:"Parent Dept"` // Parent Dept
	DeptPath       string `form:"deptPath" search:"type:exact;column:dept_path;table:sys_dept" comment:""`            // Path
	DeptName       string `form:"deptName" search:"type:exact;column:dept_name;table:sys_dept" comment:"Dept Name"`   // Dept Name
	Sort           int    `form:"sort" search:"type:exact;column:sort;table:sys_dept" comment:"Sort"`                 // Sort
	Leader         string `form:"leader" search:"type:exact;column:leader;table:sys_dept" comment:"Leader"`           // Leader
	Phone          string `form:"phone" search:"type:exact;column:phone;table:sys_dept" comment:"Phone"`              // Phone
	Email          string `form:"email" search:"type:exact;column:email;table:sys_dept" comment:"Email"`              // Email
	Status         string `form:"status" search:"type:exact;column:status;table:sys_dept" comment:"Status"`           // Status
}

func (m *SysDeptGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysDeptInsertReq struct {
	DeptId   int    `uri:"id" comment:"Code"`                                                 // Code
	ParentId int    `json:"parentId" comment:"Parent Dept" vd:"?"`                            // Parent Dept
	DeptPath string `json:"deptPath" comment:""`                                              // Path
	DeptName string `json:"deptName" comment:"Dept Name" vd:"len($)>0"`                       // Dept Name
	Sort     int    `json:"sort" comment:"Sort" vd:"?"`                                       // Sort
	Leader   string `json:"leader" comment:"Leader" vd:"@:len($)>0; msg:'leader cannot be empty'"` // Leader
	Phone    string `json:"phone" comment:"Phone" vd:"?"`                                     // Phone
	Email    string `json:"email" comment:"Email" vd:"?"`                                     // Email
	Status   int    `json:"status" comment:"Status" vd:"$>0"`                                 // Status
	common.ControlBy
}

func (s *SysDeptInsertReq) Generate(model *models.SysDept) {
	if s.DeptId != 0 {
		model.DeptId = s.DeptId
	}
	model.DeptName = s.DeptName
	model.ParentId = s.ParentId
	model.DeptPath = s.DeptPath
	model.Sort = s.Sort
	model.Leader = s.Leader
	model.Phone = s.Phone
	model.Email = s.Email
	model.Status = s.Status
}

// GetId Get the corresponding ID for the data
func (s *SysDeptInsertReq) GetId() interface{} {
	return s.DeptId
}

type SysDeptUpdateReq struct {
	DeptId   int    `uri:"id" comment:"Code"`                                                 // Code
	ParentId int    `json:"parentId" comment:"Parent Dept" vd:"?"`                            // Parent Dept
	DeptPath string `json:"deptPath" comment:""`                                              // Path
	DeptName string `json:"deptName" comment:"Dept Name" vd:"len($)>0"`                       // Dept Name
	Sort     int    `json:"sort" comment:"Sort" vd:"?"`                                       // Sort
	Leader   string `json:"leader" comment:"Leader" vd:"@:len($)>0; msg:'leader cannot be empty'"` // Leader
	Phone    string `json:"phone" comment:"Phone" vd:"?"`                                     // Phone
	Email    string `json:"email" comment:"Email" vd:"?"`                                     // Email
	Status   int    `json:"status" comment:"Status" vd:"$>0"`                                 // Status
	common.ControlBy
}

// Generate Convert struct data from SysDeptControl to the corresponding SysDept model
func (s *SysDeptUpdateReq) Generate(model *models.SysDept) {
	if s.DeptId != 0 {
		model.DeptId = s.DeptId
	}
	model.DeptName = s.DeptName
	model.ParentId = s.ParentId
	model.DeptPath = s.DeptPath
	model.Sort = s.Sort
	model.Leader = s.Leader
	model.Phone = s.Phone
	model.Email = s.Email
	model.Status = s.Status
}

// GetId Get the corresponding ID for the data
func (s *SysDeptUpdateReq) GetId() interface{} {
	return s.DeptId
}

type SysDeptGetReq struct {
	Id int `uri:"id"`
}

func (s *SysDeptGetReq) GetId() interface{} {
	return s.Id
}

type SysDeptDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *SysDeptDeleteReq) GetId() interface{} {
	return s.Ids
}

type DeptLabel struct {
	Id       int         `gorm:"-" json:"id"`
	Label    string      `gorm:"-" json:"label"`
	Children []DeptLabel `gorm:"-" json:"children"`
}
