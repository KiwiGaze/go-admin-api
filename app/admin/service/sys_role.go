package service

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	"gorm.io/gorm/clause"

	"github.com/casbin/casbin/v2"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin-api/app/admin/models"
	"go-admin-api/app/admin/service/dto"
	cDto "go-admin-api/common/dto"
)

type SysRole struct {
	service.Service
}

// GetPage gets the SysRole list.
func (e *SysRole) GetPage(c *dto.SysRoleGetPageReq, list *[]models.SysRole, count *int64) error {
	var err error
	var data models.SysRole

	err = e.Orm.Model(&data).Preload("SysMenu").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Get gets the SysRole object.
func (e *SysRole) Get(d *dto.SysRoleGetReq, model *models.SysRole) error {
	var err error
	db := e.Orm.First(model, d.GetId())
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("the requested object does not exist or you do not have permission to view it")
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	model.MenuIds, err = e.GetRoleMenuId(model.RoleId)
	if err != nil {
		e.Log.Errorf("get menuIds error, %s", err.Error())
		return err
	}
	return nil
}

// Insert creates a SysRole object.
func (e *SysRole) Insert(c *dto.SysRoleInsertReq, cb *casbin.SyncedEnforcer) error {
	var err error
	var data models.SysRole
	var dataMenu []models.SysMenu
	err = e.Orm.Preload("SysApi").Where("menu_id in ?", c.MenuIds).Find(&dataMenu).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	c.SysMenu = dataMenu
	c.Generate(&data)
	tx := e.Orm
	if config.DatabaseConfig.Driver != "sqlite3" {
		tx = e.Orm.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	var count int64
	err = tx.Model(&data).Where("role_key = ?", c.RoleKey).Count(&count).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	if count > 0 {
		err = errors.New("roleKey already exists; change it before submitting")
		e.Log.Errorf("db error:%s", err)
		return err
	}

	err = tx.Create(&data).Error
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	mp := make(map[string]interface{}, 0)
	polices := make([][]string, 0)
	for _, menu := range dataMenu {
		for _, api := range menu.SysApi {
			if mp[data.RoleKey+"-"+api.Path+"-"+api.Action] != "" {
				mp[data.RoleKey+"-"+api.Path+"-"+api.Action] = ""
				polices = append(polices, []string{data.RoleKey, api.Path, api.Action})
			}
		}
	}

	if len(polices) <= 0 {
		return nil
	}

	// Write the current role records into the sys_casbin_rule permission table.
	_, err = cb.AddNamedPolicies("p", polices)
	if err != nil {
		return err
	}

	return nil
}

// Update updates the SysRole object.
func (e *SysRole) Update(c *dto.SysRoleUpdateReq, cb *casbin.SyncedEnforcer) error {
	var err error
	tx := e.Orm
	if config.DatabaseConfig.Driver != "sqlite3" {
		tx = e.Orm.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	var model = models.SysRole{}
	var mlist = make([]models.SysMenu, 0)
	tx.Preload("SysMenu").First(&model, c.GetId())
	tx.Preload("SysApi").Where("menu_id in ?", c.MenuIds).Find(&mlist)
	err = tx.Model(&model).Association("SysMenu").Delete(model.SysMenu)
	if err != nil {
		e.Log.Errorf("delete policy error:%s", err)
		return err
	}
	c.Generate(&model)
	model.SysMenu = &mlist
	// Update associated data using FullSaveAssociations mode.
	db := tx.Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model)

	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("you do not have permission to update this data")
	}

	// Clear all records for the current role from the sys_casbin_rule permission table.
	_, err = cb.RemoveFilteredPolicy(0, model.RoleKey)
	if err != nil {
		e.Log.Errorf("delete policy error:%s", err)
		return err
	}
	mp := make(map[string]interface{}, 0)
	polices := make([][]string, 0)
	for _, menu := range mlist {
		for _, api := range menu.SysApi {
			if mp[model.RoleKey+"-"+api.Path+"-"+api.Action] != "" {
				mp[model.RoleKey+"-"+api.Path+"-"+api.Action] = ""
				//_, err = cb.AddNamedPolicy("p", model.RoleKey, api.Path, api.Action)
				polices = append(polices, []string{model.RoleKey, api.Path, api.Action})
			}
		}
	}
	if len(polices) <= 0 {
		return nil
	}

	// Write the current role records into the sys_casbin_rule permission table.
	_, err = cb.AddNamedPolicies("p", polices)
	if err != nil {
		return err
	}
	return nil
}

// Remove deletes SysRole.
func (e *SysRole) Remove(c *dto.SysRoleDeleteReq, cb *casbin.SyncedEnforcer) error {
	var err error
	tx := e.Orm
	if config.DatabaseConfig.Driver != "sqlite3" {
		tx = e.Orm.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	var model = models.SysRole{}
	tx.Preload("SysMenu").Preload("SysDept").First(&model, c.GetId())
	// When deleting SysRole, also delete all role associations in other tables (SysMenu and SysDept).
	db := tx.Select(clause.Associations).Delete(&model)

	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("you do not have permission to update this data")
	}

	// Clear all records for the current role from the sys_casbin_rule permission table.
	_, _ = cb.RemoveFilteredPolicy(0, model.RoleKey)

	return nil
}

// GetRoleMenuId gets the menu IDs associated with the role.
func (e *SysRole) GetRoleMenuId(roleId int) ([]int, error) {
	menuIds := make([]int, 0)
	model := models.SysRole{}
	model.RoleId = roleId
	if err := e.Orm.Model(&model).Preload("SysMenu").First(&model).Error; err != nil {
		return nil, err
	}
	l := *model.SysMenu
	for i := 0; i < len(l); i++ {
		menuIds = append(menuIds, l[i].MenuId)
	}
	return menuIds, nil
}

func (e *SysRole) UpdateDataScope(c *dto.RoleDataScopeReq) *SysRole {
	var err error
	tx := e.Orm
	if config.DatabaseConfig.Driver != "sqlite3" {
		tx = e.Orm.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	var dlist = make([]models.SysDept, 0)
	var model = models.SysRole{}
	tx.Preload("SysDept").First(&model, c.RoleId)
	tx.Where("dept_id in ?", c.DeptIds).Find(&dlist)
	// Delete the association between SysRole and SysDept.
	err = tx.Model(&model).Association("SysDept").Delete(model.SysDept)
	if err != nil {
		e.Log.Errorf("delete SysDept error:%s", err)
		_ = e.AddError(err)
		return e
	}
	c.Generate(&model)
	model.SysDept = dlist
	// Update associated data using FullSaveAssociations mode.
	db := tx.Model(&model).Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model)
	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
		return e
	}
	if db.RowsAffected == 0 {
		_ = e.AddError(errors.New("you do not have permission to update this data"))
		return e
	}
	return e
}

// UpdateStatus updates the SysRole status.
func (e *SysRole) UpdateStatus(c *dto.UpdateStatusReq) error {
	var err error
	tx := e.Orm
	if config.DatabaseConfig.Driver != "sqlite3" {
		tx = e.Orm.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	var model = models.SysRole{}
	tx.First(&model, c.GetId())
	c.Generate(&model)
	// Update associated data using FullSaveAssociations mode.
	db := tx.Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model)
	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("you do not have permission to update this data")
	}
	return nil
}

// GetWithName gets the SysRole object.
func (e *SysRole) GetWithName(d *dto.SysRoleByName, model *models.SysRole) *SysRole {
	var err error
	db := e.Orm.Where("role_name = ?", d.RoleName).First(model)
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("the requested object does not exist or you do not have permission to view it")
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
		return e
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
		return e
	}
	model.MenuIds, err = e.GetRoleMenuId(model.RoleId)
	if err != nil {
		e.Log.Errorf("get menuIds error, %s", err.Error())
		_ = e.AddError(err)
		return e
	}
	return e
}

// GetById gets the SysRole object.
func (e *SysRole) GetById(roleId int) ([]string, error) {
	permissions := make([]string, 0)
	model := models.SysRole{}
	model.RoleId = roleId
	if err := e.Orm.Model(&model).Preload("SysMenu").First(&model).Error; err != nil {
		return nil, err
	}
	l := *model.SysMenu
	for i := 0; i < len(l); i++ {
		if l[i].Permission != "" {
			permissions = append(permissions, l[i].Permission)
		}
	}
	return permissions, nil
}
