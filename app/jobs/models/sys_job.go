package models

import (
	"go-admin-api/common/models"

	"gorm.io/gorm"
)

type SysJob struct {
	JobId          int    `json:"jobId" gorm:"primaryKey;autoIncrement"` // Code
	JobName        string `json:"jobName" gorm:"size:255;"`              // Name
	JobGroup       string `json:"jobGroup" gorm:"size:255;"`             // Job group
	JobType        int    `json:"jobType" gorm:"size:1;"`                // Job type
	CronExpression string `json:"cronExpression" gorm:"size:255;"`       // Cron expression
	InvokeTarget   string `json:"invokeTarget" gorm:"size:255;"`         // Invocation target
	Args           string `json:"args" gorm:"size:255;"`                 // Target arguments
	MisfirePolicy  int    `json:"misfirePolicy" gorm:"size:255;"`        // Execution policy
	Concurrent     int    `json:"concurrent" gorm:"size:1;"`             // Whether concurrent
	Status         int    `json:"status" gorm:"size:1;"`                 // Status
	EntryId        int    `json:"entry_id" gorm:"size:11;"`              // ID returned when the job starts
	models.ControlBy
	models.ModelTime

	DataScope string `json:"dataScope" gorm:"-"`
}

func (*SysJob) TableName() string {
	return "sys_job"
}

func (e *SysJob) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysJob) GetId() interface{} {
	return e.JobId
}

func (e *SysJob) SetCreateBy(createBy int) {
	e.CreateBy = createBy
}

func (e *SysJob) SetUpdateBy(updateBy int) {
	e.UpdateBy = updateBy
}

func (e *SysJob) GetList(tx *gorm.DB, list interface{}) (err error) {
	return tx.Table(e.TableName()).Where("status = ?", 2).Find(list).Error
}

// Update updates SysJob.
func (e *SysJob) Update(tx *gorm.DB, id interface{}) (err error) {
	return tx.Table(e.TableName()).Where(id).Updates(&e).Error
}

func (e *SysJob) RemoveAllEntryID(tx *gorm.DB) (update SysJob, err error) {
	if err = tx.Table(e.TableName()).Where("entry_id > ?", 0).Update("entry_id", 0).Error; err != nil {
		return
	}
	return
}
