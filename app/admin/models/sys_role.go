package models

import "go-admin-api/common/models"

type SysRole struct {
	RoleId    int        `json:"roleId" gorm:"primaryKey;autoIncrement"` // Role ID
	RoleName  string     `json:"roleName" gorm:"size:128;"`              // Role name
	Status    string     `json:"status" gorm:"size:4;"`                  // Status: 1=disabled 2=enabled
	RoleKey   string     `json:"roleKey" gorm:"size:128;"`               // Role code
	RoleSort  int        `json:"roleSort" gorm:""`                       // Sort order
	Flag      string     `json:"flag" gorm:"size:128;"`                  //
	Remark    string     `json:"remark" gorm:"size:255;"`                // Remark
	Admin     bool       `json:"admin" gorm:"size:4;"`
	DataScope string     `json:"dataScope" gorm:"size:128;"`
	Params    string     `json:"params" gorm:"-"`
	MenuIds   []int      `json:"menuIds" gorm:"-"`
	DeptIds   []int      `json:"deptIds" gorm:"-"`
	SysDept   []SysDept  `json:"sysDept" gorm:"many2many:sys_role_dept;foreignKey:RoleId;joinForeignKey:role_id;references:DeptId;joinReferences:dept_id;"`
	SysMenu   *[]SysMenu `json:"sysMenu" gorm:"many2many:sys_role_menu;foreignKey:RoleId;joinForeignKey:role_id;references:MenuId;joinReferences:menu_id;"`
	models.ControlBy
	models.ModelTime
}

func (*SysRole) TableName() string {
	return "sys_role"
}

func (e *SysRole) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysRole) GetId() interface{} {
	return e.RoleId
}
