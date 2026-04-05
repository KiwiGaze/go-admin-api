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

type SysLoginLog struct {
	models.Model
	Username      string    `json:"username" gorm:"size:128;comment:username"`
	Status        string    `json:"status" gorm:"size:4;comment:status"`
	Ipaddr        string    `json:"ipaddr" gorm:"size:255;comment:ip address"`
	LoginLocation string    `json:"loginLocation" gorm:"size:255;comment:login location"`
	Browser       string    `json:"browser" gorm:"size:255;comment:browser"`
	Os            string    `json:"os" gorm:"size:255;comment:operating system"`
	Platform      string    `json:"platform" gorm:"size:255;comment:platform"`
	LoginTime     time.Time `json:"loginTime" gorm:"comment:login time"`
	Remark        string    `json:"remark" gorm:"size:255;comment:remark"`
	Msg           string    `json:"msg" gorm:"size:255;comment:message"`
	CreatedAt     time.Time `json:"createdAt" gorm:"comment:created time"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"comment:last updated time"`
	models.ControlBy
}

func (*SysLoginLog) TableName() string {
	return "sys_login_log"
}

func (e *SysLoginLog) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysLoginLog) GetId() interface{} {
	return e.Id
}

// SaveLoginLog retrieves login log from the queue
func SaveLoginLog(message storage.Messager) (err error) {
	// prepare db
	db := sdk.Runtime.GetDbByTenant(message.GetPrefix())
	if db == nil {
		err = errors.New("db not exist")
		log.Errorf("host[%s]'s %s", message.GetPrefix(), err.Error())
		return err
	}
	var rb []byte
	rb, err = json.Marshal(message.GetValues())
	if err != nil {
		log.Errorf("json Marshal error, %s", err.Error())
		return err
	}
	var l SysLoginLog
	err = json.Unmarshal(rb, &l)
	if err != nil {
		log.Errorf("json Unmarshal error, %s", err.Error())
		return err
	}
	err = db.Create(&l).Error
	if err != nil {
		log.Errorf("db create error, %s", err.Error())
		return err
	}
	return nil
}
