package apis

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"gorm.io/gorm"

	"go-admin-api/common/service"
)

type Api struct {
	Context *gin.Context
	Logger  *logger.Helper
	Orm     *gorm.DB
	Errors  error
}

func (e *Api) AddError(err error) {
	if e.Errors == nil {
		e.Errors = err
	} else if err != nil {
		e.Logger.Error(err)
		e.Errors = fmt.Errorf("%v; %w", e.Errors, err)
	}
}

// MakeContext sets the HTTP context.
func (e *Api) MakeContext(c *gin.Context) *Api {
	e.Context = c
	e.Logger = api.GetRequestLogger(c)
	return e
}

// GetLogger returns the logger provided by the context.
func (e *Api) GetLogger() *logger.Helper {
	return api.GetRequestLogger(e.Context)
}

func (e *Api) Bind(d interface{}, bindings ...binding.Binding) *Api {
	var err error
	if len(bindings) == 0 {
		bindings = append(bindings, binding.JSON, nil)
	}
	for i := range bindings {
		switch bindings[i] {
		case binding.JSON:
			err = e.Context.ShouldBindWith(d, binding.JSON)
		case binding.XML:
			err = e.Context.ShouldBindWith(d, binding.XML)
		case binding.Form:
			err = e.Context.ShouldBindWith(d, binding.Form)
		case binding.Query:
			err = e.Context.ShouldBindWith(d, binding.Query)
		case binding.FormPost:
			err = e.Context.ShouldBindWith(d, binding.FormPost)
		case binding.FormMultipart:
			err = e.Context.ShouldBindWith(d, binding.FormMultipart)
		case binding.ProtoBuf:
			err = e.Context.ShouldBindWith(d, binding.ProtoBuf)
		case binding.MsgPack:
			err = e.Context.ShouldBindWith(d, binding.MsgPack)
		case binding.YAML:
			err = e.Context.ShouldBindWith(d, binding.YAML)
		case binding.Header:
			err = e.Context.ShouldBindWith(d, binding.Header)
		default:
			err = e.Context.ShouldBindUri(d)
		}
		if err != nil {
			e.AddError(err)
		}
	}
	return e
}

// GetOrm returns the ORM database handle.
func (e *Api) GetOrm() (*gorm.DB, error) {
	db, err := pkg.GetOrm(e.Context)
	if err != nil {
		e.Error(500, err, "failed to get database connection")
		return nil, err
	}
	return db, nil
}

// MakeOrm sets the ORM database handle.
func (e *Api) MakeOrm() *Api {
	var err error
	if e.Logger == nil {
		err = errors.New("at MakeOrm logger is nil")
		//e.Logger.Error(500, err, "at MakeOrm logger is nil")
		e.AddError(err)
		return e
	}
	db, err := pkg.GetOrm(e.Context)
	if err != nil {
		e.Logger.Error(500, err, "failed to get database connection")
		e.AddError(err)
	}
	e.Orm = db
	return e
}

func (e *Api) MakeService(c *service.Service) *Api {
	c.Log = e.Logger
	c.Orm = e.Orm
	return e
}

// Error handles standard error responses.
func (e *Api) Error(code int, err error, msg string) {
	response.Error(e.Context, code, err, msg)
}

// OK handles standard success responses.
func (e *Api) OK(data interface{}, msg string) {
	response.OK(e.Context, data, msg)
}

// PageOK handles paginated responses.
func (e *Api) PageOK(result interface{}, count int, pageIndex int, pageSize int, msg string) {
	response.PageOK(e.Context, result, count, pageIndex, pageSize, msg)
}

// Custom is a compatibility helper.
func (e *Api) Custom(data gin.H) {
	response.Custum(e.Context, data)
}
