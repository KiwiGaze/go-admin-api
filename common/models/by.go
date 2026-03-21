package models

import (
	"time"

	"gorm.io/gorm"
)

type ControlBy struct {
	CreateBy int `json:"createdBy" gorm:"index;comment:Created by"`
	UpdateBy int `json:"updatedBy" gorm:"index;comment:Updated by"`
}

func (e *ControlBy) SetCreateBy(createBy int) {
	e.CreateBy = createBy 
}

func (e *ControlBy) SetUpdateBy(updateBy int) {
	e.UpdateBy = updateBy
}

type Model struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement;comment:Primary key ID"`
}

type ModelTime struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:Creation time"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:Update time"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:Deletion time"`
}