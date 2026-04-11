package apis

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin-api/app/jobs/service"
	"go-admin-api/common/actions"
	"go-admin-api/common/dto"
)

type SysJob struct {
	api.Api
}

// RemoveJobForService invokes the service implementation.
func (e SysJob) RemoveJobForService(c *gin.Context) {
	v := dto.GeneralDelDto{}
	s := service.SysJob{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&v).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}

	s.Cron = sdk.Runtime.GetCrontabByTenant(c.Request.Host)
	err = s.RemoveJob(&v, actions.GetPermissionFromContext(c))
	if errors.Is(err, service.ErrJobNotFoundOrNotVisible) {
		e.Error(http.StatusNotFound, nil, err.Error())
		return
	}
	if err != nil {
		e.Logger.Errorf("RemoveJob error, %s", err.Error())
		e.Error(500, err, "")
		return
	}
	e.OK(nil, s.Msg)
}

// StartJobForService starts the job through the service implementation.
func (e SysJob) StartJobForService(c *gin.Context) {
	e.MakeContext(c)
	log := e.GetLogger()
	db, err := e.GetOrm()
	if err != nil {
		log.Error(err)
		return
	}
	var v dto.GeneralGetDto
	err = c.BindUri(&v)
	if err != nil {
		log.Warnf("Parameter validation error: %s", err)
		e.Error(http.StatusUnprocessableEntity, err, "Parameter validation failed")
		return
	}
	s := service.SysJob{}
	s.Orm = db
	s.Log = log
	s.Cron = sdk.Runtime.GetCrontabByTenant(c.Request.Host)
	err = s.StartJob(&v, actions.GetPermissionFromContext(c))
	if errors.Is(err, service.ErrJobNotFoundOrNotVisible) {
		e.Error(http.StatusNotFound, nil, err.Error())
		return
	}
	if err != nil {
		log.Errorf("StartJob error, %s", err.Error())
		e.Error(500, err, err.Error())
		return
	}
	e.OK(nil, s.Msg)
}
