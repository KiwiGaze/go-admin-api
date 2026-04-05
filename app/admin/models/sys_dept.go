package models

import "go-admin-api/common/models"

type SysDept struct {
	DeptId   int    `json:"deptId" gorm:"primaryKey;autoIncrement;"` // Department ID
	ParentId int    `json:"parentId" gorm:""`                        // Parent department
	DeptPath string `json:"deptPath" gorm:"size:255;"`               // Department path
	DeptName string `json:"deptName"  gorm:"size:128;"`              // Department name
	Sort     int    `json:"sort" gorm:"size:4;"`                     // Sort order
	Leader   string `json:"leader" gorm:"size:128;"`                 // Leader
	Phone    string `json:"phone" gorm:"size:11;"`                   // Phone
	Email    string `json:"email" gorm:"size:64;"`                   // Email
	Status   int    `json:"status" gorm:"size:4;"`                   // Status
	models.ControlBy
	models.ModelTime
	DataScope string    `json:"dataScope" gorm:"-"`
	Params    string    `json:"params" gorm:"-"`
	Children  []SysDept `json:"children" gorm:"-"`
}

func (*SysDept) TableName() string {
	return "sys_dept"
}

func (e *SysDept) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysDept) GetId() interface{} {
	return e.DeptId
}