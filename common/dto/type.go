package dto

import (
	"go-admin-api/common/models"

	"github.com/gin-gonic/gin"
)

type Index interface {
	GetIndex() Index
	Bind(ctx *gin.Context) error
	GetPageIndex() int
	GetPageSize() int
	GetNeedSearch() interface{}
}

type Control interface {
	Generate() Control
	Bind(ctx *gin.Context) error
	GenerateM() (models.ActiveRecord, error)
	GetId() interface{}
}