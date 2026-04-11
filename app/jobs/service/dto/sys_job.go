package dto

import (
	"errors"
	"go-admin-api/app/jobs/models"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"

	"go-admin-api/common/dto"
	common "go-admin-api/common/models"
)

type SysJobSearch struct {
	dto.Pagination `search:"-"`
	JobId          int    `form:"jobId" search:"type:exact;column:job_id;table:sys_job"`
	JobName        string `form:"jobName" search:"type:icontains;column:job_name;table:sys_job"`
	JobGroup       string `form:"jobGroup" search:"type:exact;column:job_group;table:sys_job"`
	CronExpression string `form:"cronExpression" search:"type:exact;column:cron_expression;table:sys_job"`
	InvokeTarget   string `form:"invokeTarget" search:"type:exact;column:invoke_target;table:sys_job"`
	Status         int    `form:"status" search:"type:exact;column:status;table:sys_job"`
}

func (m *SysJobSearch) GetNeedSearch() interface{} {
	return *m
}

func (m *SysJobSearch) GetIndex() dto.Index {
	o := *m
	return &o
}

func (m *SysJobSearch) Bind(ctx *gin.Context) error {
	log := api.GetRequestLogger(ctx)
	err := ctx.ShouldBind(m)
	if err != nil {
		log.Errorf("Bind error: %s", err)
	}
	return err
}

func (m *SysJobSearch) Generate() *SysJobSearch {
	o := *m
	return &o
}

type SysJobControl struct {
	JobId          int    `json:"jobId"`
	JobName        string `json:"jobName" validate:"required"` // Name
	JobGroup       string `json:"jobGroup"`                    // Job group
	JobType        int    `json:"jobType"`                     // Job type
	CronExpression string `json:"cronExpression"`              // Cron expression
	InvokeTarget   string `json:"invokeTarget"`                // Invocation target
	Args           string `json:"args"`                        // Target arguments
	MisfirePolicy  int    `json:"misfirePolicy"`               // Execution policy
	Concurrent     int    `json:"concurrent"`                  // Whether concurrent execution is allowed
	Status         int    `json:"status"`                      // Status
	EntryId        int    `json:"entryId"`                     // ID returned when the job starts
}

func (s *SysJobControl) Bind(ctx *gin.Context) error {
	return ctx.ShouldBind(s)
}

func (s *SysJobControl) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysJobControl) GenerateM() (common.ActiveRecord, error) {
	return &models.SysJob{
		JobId:          s.JobId,
		JobName:        s.JobName,
		JobGroup:       s.JobGroup,
		JobType:        s.JobType,
		CronExpression: s.CronExpression,
		InvokeTarget:   s.InvokeTarget,
		Args:           s.Args,
		MisfirePolicy:  s.MisfirePolicy,
		Concurrent:     s.Concurrent,
		Status:         s.Status,
		EntryId:        s.EntryId,
	}, nil
}

func (s *SysJobControl) GetId() interface{} {
	return s.JobId
}

type SysJobById struct {
	dto.ObjectById
}

type SysJobDeleteReq struct {
	Id  int   `json:"id"`
	Ids []int `json:"ids"`
}

func (s *SysJobById) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysJobById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysJob{}, nil
}

func (s *SysJobDeleteReq) Bind(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(s); err != nil {
		return err
	}

	if s.Id > 0 {
		s.Ids = append(s.Ids, s.Id)
	}
	if len(s.Ids) == 0 {
		return errors.New("at least one job id is required")
	}
	return nil
}

func (s *SysJobDeleteReq) Generate() dto.Control {
	cp := *s
	return &cp
}

func (s *SysJobDeleteReq) GenerateM() (common.ActiveRecord, error) {
	return &models.SysJob{}, nil
}

func (s *SysJobDeleteReq) GetId() interface{} {
	return s.Ids
}

type SysJobItem struct {
	JobId          int    `json:"jobId"`
	JobName        string `json:"jobName" validate:"required"` // Name
	JobGroup       string `json:"jobGroup"`                    // Job group
	JobType        int    `json:"jobType"`                     // Job type
	CronExpression string `json:"cronExpression"`              // Cron expression
	InvokeTarget   string `json:"invokeTarget"`                // Invocation target
	Args           string `json:"args"`                        // Target arguments
	MisfirePolicy  int    `json:"misfirePolicy"`               // Execution policy
	Concurrent     int    `json:"concurrent"`                  // Whether concurrent execution is allowed
	Status         int    `json:"status"`                      // Status
	EntryId        int    `json:"entryId"`                     // ID returned when the job starts
}
