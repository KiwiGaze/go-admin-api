package models

import "go-admin-api/common/models"

type SysPost struct {
	PostId   int    `gorm:"primaryKey;autoIncrement" json:"postId"` // Post ID
	PostName string `gorm:"size:128;" json:"postName"`              // Post name
	PostCode string `gorm:"size:128;" json:"postCode"`              // Post code
	Sort     int    `gorm:"size:4;" json:"sort"`                    // Sort order
	Status   int    `gorm:"size:4;" json:"status"`                  // Status
	Remark   string `gorm:"size:255;" json:"remark"`                // Remark
	models.ControlBy
	models.ModelTime

	DataScope string `gorm:"-" json:"dataScope"`
	Params    string `gorm:"-" json:"params"`
}

func (*SysPost) TableName() string {
	return "sys_post"
}

func (e *SysPost) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysPost) GetId() interface{} {
	return e.PostId
}
