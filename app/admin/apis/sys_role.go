package apis

import (
	"fmt"
	"go-admin-api/common/global"
	"net/http"

	"go-admin-api/app/admin/models"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin-api/app/admin/service"
	"go-admin-api/app/admin/service/dto"
)

type SysRole struct {
	api.Api
}

// GetPage
// @Summary Role list data
// @Description Get JSON
// @Tags Role
// @Param roleName query string false "roleName"
// @Param status query string false "status"
// @Param roleKey query string false "roleKey"
// @Param pageSize query int false "items per page"
// @Param pageIndex query int false "page number"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/role [get]
// @Security Bearer
func (e SysRole) GetPage(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleGetPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	list := make([]models.SysRole, 0)
	var count int64

	err = s.GetPage(&req, &list, &count)
	if err != nil {
		e.Error(500, err, "query failed")
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "query successful")
}

// Get
// @Summary Get role data
// @Description Get JSON
// @Tags Role
// @Param roleId path string false "roleId"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/role/{id} [get]
// @Security Bearer
func (e SysRole) Get(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleGetReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, fmt.Sprintf(" %s ", err.Error()))
		return
	}

	var object models.SysRole

	err = s.Get(&req, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "query failed")
		return
	}

	e.OK(object, "query successful")
}

// Insert
// @Summary Create role
// @Description Get JSON
// @Tags Role
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysRoleInsertReq true "data"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/role [post]
// @Security Bearer
func (e SysRole) Insert(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleInsertReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// Set creator
	req.CreateBy = user.GetUserId(c)
	if req.Status == "" {
		req.Status = "2"
	}
	cb := sdk.Runtime.GetCasbinByTenant(c.Request.Host)
	err = s.Insert(&req, cb)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, "creation failed, "+err.Error())
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, "creation failed, "+err.Error())
		return
	}
	e.OK(req.GetId(), "created successfully")
}

// Update user role
// @Summary Update user role
// @Description Get JSON
// @Tags Role
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysRoleUpdateReq true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/role/{id} [put]
// @Security Bearer
func (e SysRole) Update(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleUpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	cb := sdk.Runtime.GetCasbinByTenant(c.Request.Host)

	req.SetUpdateBy(user.GetUserId(c))

	err = s.Update(&req, cb)
	if err != nil {
		e.Logger.Error(err)
		return
	}

	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, "update failed, "+err.Error())
		return
	}

	e.OK(req.GetId(), "updated successfully")
}

// Delete
// @Summary Delete user role
// @Description Delete data
// @Tags Role
// @Param data body dto.SysRoleDeleteReq true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/role [delete]
// @Security Bearer
func (e SysRole) Delete(c *gin.Context) {
	s := new(service.SysRole)
	req := dto.SysRoleDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, fmt.Sprintf("failed to delete role %v,\r\nerror: %s", req.Ids, err.Error()))
		return
	}

	cb := sdk.Runtime.GetCasbinByTenant(c.Request.Host)
	err = s.Remove(&req, cb)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, "")
		return
	}

	e.OK(req.GetId(), fmt.Sprintf("role %v deleted successfully", req.GetId()))
}

// Update2Status update user role status
// @Summary Update user role
// @Description Get JSON
// @Tags Role
// @Accept  application/json
// @Product application/json
// @Param data body dto.UpdateStatusReq true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/role-status/{id} [put]
// @Security Bearer
func (e SysRole) Update2Status(c *gin.Context) {
	s := service.SysRole{}
	req := dto.UpdateStatusReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, fmt.Sprintf("failed to update role status, reason: %s ", err.Error()))
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	err = s.UpdateStatus(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("failed to update role status, reason: %s ", err.Error()))
		return
	}
	e.OK(req.GetId(), fmt.Sprintf("role %v status updated successfully", req.GetId()))
}

// Update2DataScope update role data permissions
// @Summary Update role data permissions
// @Description Get JSON
// @Tags Role
// @Accept  application/json
// @Product application/json
// @Param data body dto.RoleDataScopeReq true "body"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/role-status/{id} [put]
// @Security Bearer
func (e SysRole) Update2DataScope(c *gin.Context) {
	s := service.SysRole{}
	req := dto.RoleDataScopeReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	data := &models.SysRole{
		RoleId:    req.RoleId,
		DataScope: req.DataScope,
		DeptIds:   req.DeptIds,
	}
	data.UpdateBy = user.GetUserId(c)
	err = s.UpdateDataScope(&req).Error
	if err != nil {
		e.Error(500, err, fmt.Sprintf("failed to update role data permissions, error details: %s", err.Error()))
		return
	}
	e.OK(nil, "operation successful")
}
