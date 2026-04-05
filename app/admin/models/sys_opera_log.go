package models

import (
	"encoding/json"
	"errors"
	"time"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/storage"

	"go-admin-api/common/models"
)

type SysOperaLog struct {
	models.Model
	Title         string    `json:"title" gorm:"size:255;comment:Operation Module"`
	BusinessType  string    `json:"businessType" gorm:"size:128;comment:Operation Type"`
	BusinessTypes string    `json:"businessTypes" gorm:"size:128;comment:BusinessTypes"`
	Method        string    `json:"method" gorm:"size:128;comment:Function"`
	RequestMethod string    `json:"requestMethod" gorm:"size:128;comment:Request Method GET POST PUT DELETE"`
	OperatorType  string    `json:"operatorType" gorm:"size:128;comment:Operator Type"`
	OperName      string    `json:"operName" gorm:"size:128;comment:Operator"`
	DeptName      string    `json:"deptName" gorm:"size:128;comment:Department Name"`
	OperUrl       string    `json:"operUrl" gorm:"size:255;comment:Access URL"`
	OperIp        string    `json:"operIp" gorm:"size:128;comment:Client IP"`
	OperLocation  string    `json:"operLocation" gorm:"size:128;comment:Access Location"`
	OperParam     string    `json:"operParam" gorm:"text;comment:Request Parameters"`
	Status        string    `json:"status" gorm:"size:4;comment:Operation Status 1:Normal 2:Closed"`
	OperTime      time.Time `json:"operTime" gorm:"comment:Operation Time"`
	JsonResult    string    `json:"jsonResult" gorm:"size:255;comment:Response Data"`
	Remark        string    `json:"remark" gorm:"size:255;comment:Remark"`
	LatencyTime   string    `json:"latencyTime" gorm:"size:128;comment:Latency"`
	UserAgent     string    `json:"userAgent" gorm:"size:255;comment:User Agent"`
	CreatedAt     time.Time `json:"createdAt" gorm:"comment:Created At"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"comment:Last Updated At"`
	models.ControlBy
}

func (*SysOperaLog) TableName() string {
	return "sys_opera_log"
}

func (e *SysOperaLog) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysOperaLog) GetId() interface{} {
	return e.Id
}

// SaveOperaLog retrieves operation logs from the queue
func SaveOperaLog(message storage.Messager) (err error) {
	// Prepare db
	db := sdk.Runtime.GetDbByTenant(message.GetPrefix())
	if db == nil {
		err = errors.New("db not exist")
		log.Errorf("host[%s]'s %s", message.GetPrefix(), err.Error())
		// Log writing to the database ignores error
		return nil
	}
	var rb []byte
	rb, err = json.Marshal(message.GetValues())
	if err != nil {
		log.Errorf("json Marshal error, %s", err.Error())
		// Log writing to the database ignores error
		return nil
	}
	var l SysOperaLog
	err = json.Unmarshal(rb, &l)
	if err != nil {
		log.Errorf("json Unmarshal error, %s", err.Error())
		// Log writing to the database ignores error
		return nil
	}
	// Truncate return value if it exceeds 100 characters
	if len(l.JsonResult) > 100 {
		l.JsonResult = l.JsonResult[:100]
	}
	err = db.Create(&l).Error
	if err != nil {
		log.Errorf("db create error, %s", err.Error())
		// Log writing to the database ignores error
		return nil
	}
	return nil
}
