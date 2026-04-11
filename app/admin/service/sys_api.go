package service

import (
	"errors"
	"fmt"

	"go-admin-api/app/admin/models"
	"go-admin-api/app/admin/service/dto"
	"go-admin-api/common/actions"
	cDto "go-admin-api/common/dto"
	"go-admin-api/common/global"

	"github.com/go-admin-team/go-admin-core/sdk/runtime"
	"github.com/go-admin-team/go-admin-core/sdk/service"
)

type SysApi struct {
	service.Service
}

const legacyEmptyTypeLabel = "\u6682\u65e0"

// GetPage gets the SysApi list.
func (e *SysApi) GetPage(c *dto.SysApiGetPageReq, p *actions.DataPermission, list *[]models.SysApi, count *int64) error {
	var err error
	var data models.SysApi

	orm := e.Orm.Debug().Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		)
	if c.Type != "" {
		qType := c.Type
		if qType == "None" || qType == legacyEmptyTypeLabel {
			qType = ""
		}
		if global.Driver == "postgres" {
			orm = orm.Where("type = ?", qType)
		} else {
			orm = orm.Where("`type` = ?", qType)
		}

	}
	err = orm.Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("Service GetSysApiPage error:%s", err)
		return err
	}
	return nil
}

// Get retrieves the SysApi object by ID.
func (e *SysApi) Get(d *dto.SysApiGetReq, p *actions.DataPermission, model *models.SysApi) *SysApi {
	var data models.SysApi
	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		FirstOrInit(model, d.GetId()).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
		return e
	}
	if model.Id == 0 {
		err = errors.New("the requested object does not exist or cannot be viewed")
		e.Log.Errorf("Service GetSysApi error: %s", err)
		_ = e.AddError(err)
		return e
	}
	return e
}

// Update updates the SysApi object.
func (e *SysApi) Update(c *dto.SysApiUpdateReq, p *actions.DataPermission) error {
	var model = models.SysApi{}
	db := e.Orm.Debug().First(&model, c.GetId())
	if db.RowsAffected == 0 {
		return errors.New("no permission to update this data")
	}
	c.Generate(&model)
	db = e.Orm.Save(&model)
	if err := db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysApi error:%s", err)
		return err
	}

	return nil
}

// Remove deletes the SysApi object.
func (e *SysApi) Remove(d *dto.SysApiDeleteReq, p *actions.DataPermission) error {
	var data models.SysApi

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSysApi error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("no permission to delete this data")
	}
	return nil
}

// CheckStorageSysApi creates the SysApi object.
func (e *SysApi) CheckStorageSysApi(c *[]runtime.Router) error {
	for _, v := range *c {
		err := e.Orm.Debug().Where(models.SysApi{Path: v.RelativePath, Action: v.HttpMethod}).
			Attrs(models.SysApi{Handle: v.Handler}).
			FirstOrCreate(&models.SysApi{}).Error
		if err != nil {
			err := fmt.Errorf("Service CheckStorageSysApi error: %s \r\n ", err.Error())
			return err
		}
	}
	return nil
}
