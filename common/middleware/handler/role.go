package handler

import "go-admin-api/common/models"

type SysRole struct {
	RoleId    int    `json:"roleId" gorm:"primaryKey;autoIncrement"` // Role ID
	RoleName  string `json:"roleName" gorm:"size:128;"`              // Role name
	Status    string `json:"status" gorm:"size:4;"`                  //
	RoleKey   string `json:"roleKey" gorm:"size:128;"`               // Role key
	RoleSort  int    `json:"roleSort" gorm:""`                       // Role order
	Flag      string `json:"flag" gorm:"size:128;"`                  //
	Remark    string `json:"remark" gorm:"size:255;"`                // Remark
	Admin     bool   `json:"admin" gorm:"size:4;"`
	DataScope string `json:"dataScope" gorm:"size:128;"`
	Params    string `json:"params" gorm:"-"`
	MenuIds   []int  `json:"menuIds" gorm:"-"`
	DeptIds   []int  `json:"deptIds" gorm:"-"`
	models.ControlBy
	models.ModelTime
}

func (SysRole) TableName() string {
	return "sys_role"
}
