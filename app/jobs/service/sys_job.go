package service

import (
	"errors"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"go-admin-api/app/jobs"
	"go-admin-api/app/jobs/models"
	"go-admin-api/common/actions"
	"go-admin-api/common/dto"
)

type SysJob struct {
	service.Service
	Cron *cron.Cron
}

var (
	ErrSchedulerNotInitialized = errors.New("scheduler is not initialized")
	ErrJobAlreadyStarted       = errors.New("job is already started")
	ErrJobNotFoundOrNotVisible = errors.New("the requested object does not exist or you do not have permission to view it")
)

// RemoveJob removes a job.
func (e *SysJob) RemoveJob(c *dto.GeneralDelDto, p *actions.DataPermission) error {
	if p == nil {
		p = &actions.DataPermission{}
	}
	var err error
	var data models.SysJob
	err = e.Orm.Table(data.TableName()).
		Scopes(actions.Permission(data.TableName(), p)).
		First(&data, c.Id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrJobNotFoundOrNotVisible
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if data.EntryId == 0 {
		e.Msg = "Job is not running."
		return nil
	}
	if e.Cron == nil {
		return ErrSchedulerNotInitialized
	}
	cn := jobs.Remove(e.Cron, data.EntryId)

	select {
	case res := <-cn:
		if res {
			err = e.Orm.Table(data.TableName()).Where("entry_id = ?", data.EntryId).Update("entry_id", 0).Error
			if err != nil {
				e.Log.Errorf("db error: %s", err)
			}
			return err
		}
	case <-time.After(time.Second * 1):
		e.Msg = "Operation timed out!"
		return nil
	}
	return nil
}

// StartJob starts a job.
func (e *SysJob) StartJob(c *dto.GeneralGetDto, p *actions.DataPermission) error {
	if p == nil {
		p = &actions.DataPermission{}
	}
	var data models.SysJob
	var err error
	err = e.Orm.Table(data.TableName()).
		Scopes(actions.Permission(data.TableName(), p)).
		First(&data, c.Id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrJobNotFoundOrNotVisible
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}

	if data.Status == 1 {
		err = errors.New("The current job is disabled and cannot be started. Please enable it first.")
		return err
	}
	if data.EntryId > 0 {
		return ErrJobAlreadyStarted
	}
	if e.Cron == nil {
		return ErrSchedulerNotInitialized
	}

	if data.JobType == 1 {
		var j = &jobs.HttpJob{}
		j.InvokeTarget = data.InvokeTarget
		j.CronExpression = data.CronExpression
		j.JobId = data.JobId
		j.Name = data.JobName
		data.EntryId, err = jobs.AddJob(e.Cron, j)
		if err != nil {
			e.Log.Errorf("jobs AddJob[HttpJob] error: %s", err)
		}
	} else {
		var j = &jobs.ExecJob{}
		j.InvokeTarget = data.InvokeTarget
		j.CronExpression = data.CronExpression
		j.JobId = data.JobId
		j.Name = data.JobName
		j.Args = data.Args
		data.EntryId, err = jobs.AddJob(e.Cron, j)
		if err != nil {
			e.Log.Errorf("jobs AddJob[ExecJob] error: %s", err)
		}
	}
	if err != nil {
		return err
	}

	err = e.Orm.Table(data.TableName()).Where("job_id = ?", data.JobId).Updates(&data).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
	}
	return err
}
