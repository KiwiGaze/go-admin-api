package handler

import (
	"go-admin-api/common/models"

	"gorm.io/gorm"
)

type SysUser struct {
	UserId   int    `gorm:"primaryKey;autoIncrement;comment:Code"  json:"userId"`
	Username string `json:"username" gorm:"size:64;comment:Username"`
	Password string `json:"-" gorm:"size:128;comment:Password"`
	NickName string `json:"nickName" gorm:"size:128;comment:Nickname"`
	Phone    string `json:"phone" gorm:"size:11;comment:Phone"`
	RoleId   int    `json:"roleId" gorm:"size:20;comment:Role ID"`
	Salt     string `json:"-" gorm:"size:255;comment:Salt"`
	Avatar   string `json:"avatar" gorm:"size:255;comment:Avatar"`
	Sex      string `json:"sex" gorm:"size:255;comment:Gender"`
	Email    string `json:"email" gorm:"size:128;comment:Email"`
	DeptId   int    `json:"deptId" gorm:"size:20;comment:Department"`
	PostId   int    `json:"postId" gorm:"size:20;comment:Position"`
	Remark   string `json:"remark" gorm:"size:255;comment:Remark"`
	Status   string `json:"status" gorm:"size:4;comment:Status"`
	DeptIds  []int  `json:"deptIds" gorm:"-"`
	PostIds  []int  `json:"postIds" gorm:"-"`
	RoleIds  []int  `json:"roleIds" gorm:"-"`
	//Dept     *SysDept `json:"dept"`
	models.ControlBy
	models.ModelTime
}

func (*SysUser) TableName() string {
	return "sys_user"
}

func (e *SysUser) AfterFind(_ *gorm.DB) error {
	e.DeptIds = []int{e.DeptId}
	e.PostIds = []int{e.PostId}
	e.RoleIds = []int{e.RoleId}
	return nil
}
