package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/runtime"
	"github.com/go-admin-team/go-admin-core/storage"

	"go-admin-api/common/models"
)

type SysApi struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement;comment:primary key"`
	Handle string `json:"handle" gorm:"size:128;comment:handle"`
	Title  string `json:"title" gorm:"size:128;comment:title"`
	Path   string `json:"path" gorm:"size:128;comment:path"`
	Action string `json:"action" gorm:"size:16;comment:request method"`
	Type   string `json:"type" gorm:"size:16;comment:API type"`
	models.ModelTime
	models.ControlBy
}

func (*SysApi) TableName() string {
	return "sys_api"
}

func (e *SysApi) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysApi) GetId() interface{} {
	return e.Id
}

func SaveSysApi(message storage.Messager) (err error) {
	var rb []byte
	rb, err = json.Marshal(message.GetValues())
	if err != nil {
		err = fmt.Errorf("json Marshal error, %v", err.Error())
		return err
	}

	var l runtime.Routers
	err = json.Unmarshal(rb, &l)
	if err != nil {
		err = fmt.Errorf("json Unmarshal error, %s", err.Error())
		return err
	}
	dbList := sdk.Runtime.GetAllDb()
	for _, d := range dbList {
		for _, v := range l.List {
			if v.HttpMethod != "HEAD" ||
				strings.Contains(v.RelativePath, "/swagger/") ||
				strings.Contains(v.RelativePath, "/static/") ||
				strings.Contains(v.RelativePath, "/form-generator/") ||
				strings.Contains(v.RelativePath, "/sys/tables") {

				// Fill API title from @Summary in handler method comments; used by the code generator.
				// Can add path prefix checks here to selectively fill titles for code-generated apps only.
				jsonFile, _ := os.ReadFile("docs/swagger.json")
				jsonData, _ := simplejson.NewFromReader(bytes.NewReader(jsonFile))
				urlPath := v.RelativePath
				idPatten := "(.*)/:(\\w+)" // regex to replace :id with {id}
				reg, _ := regexp.Compile(idPatten)
				if reg.MatchString(urlPath) {
					urlPath = reg.ReplaceAllString(v.RelativePath, "${1}/{${2}}") // replace :id with {id}
				}
				apiTitle, _ := jsonData.Get("paths").Get(urlPath).Get(strings.ToLower(v.HttpMethod)).Get("summary").String()

				err := d.Debug().Where(SysApi{Path: v.RelativePath, Action: v.HttpMethod}).
					Attrs(SysApi{Handle: v.Handler, Title: apiTitle}).
					FirstOrCreate(&SysApi{}).
					//Update("handle", v.Handler).
					Error
				if err != nil {
					err := fmt.Errorf("Models SaveSysApi error: %s \r\n ", err.Error())
					return err
				}
			}
		}
	}
	return nil
}
