package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin-api/app/admin/models"
	"go-admin-api/app/admin/service"
	"go-admin-api/app/admin/service/dto"
	"go-admin-api/common/actions"
)

type SysApi struct {
	api.Api
}

// GetPage gets the API management list
// @Summary Get API management list
// @Description Get API management list
// @Tags API Management
// @Param name query string false "Name"
// @Param title query string false "Title"
// @Param path query string false "Path"
// @Param action query string false "Type"
// @Param pageSize query int false "Items per page"
// @Param pageIndex query int false "Page number"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SysApi}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-api [get]
// @Security Bearer
func (e SysApi) GetPage(c *gin.Context) {
	s := service.SysApi{}
	req := dto.SysApiGetPageReq{}
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
	// Check data permissions
	p := actions.GetPermissionFromContext(c)
	list := make([]models.SysApi, 0)
	var count int64
	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, "Query failed")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "Query successful")
}

// Get gets API management details
// @Summary Get API management details
// @Description Get API management details
// @Tags API Management
// @Param id path string false "id"
// @Success 200 {object} response.Response{data=models.SysApi} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-api/{id} [get]
// @Security Bearer
func (e SysApi) Get(c *gin.Context) {
	req := dto.SysApiGetReq{}
	s := service.SysApi{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)
	var object models.SysApi
	err = s.Get(&req, p, &object).Error
	if err != nil {
		e.Error(500, err, "Query failed")
		return
	}
	e.OK(object, "Query successful")
}

// Update updates API management
// @Summary Update API management
// @Description Update API management
// @Tags API Management
// @Accept application/json
// @Product application/json
// @Param data body dto.SysApiUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "Update successful"}"
// @Router /api/v1/sys-api/{id} [put]
// @Security Bearer
func (e SysApi) Update(c *gin.Context) {
	req := dto.SysApiUpdateReq{}
	s := service.SysApi{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, "Update failed")
		return
	}
	e.OK(req.GetId(), "Update successful")
}

// DeleteSysApi deletes API management
// @Summary Delete API management
// @Description Delete API management
// @Tags API Management
// @Param data body dto.SysApiDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "Delete successful"}"
// @Router /api/v1/sys-api [delete]
// @Security Bearer
func (e SysApi) DeleteSysApi(c *gin.Context) {
	req := dto.SysApiDeleteReq{}
	s := service.SysApi{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		return
	}
	p := actions.GetPermissionFromContext(c)
	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, "Delete failed")
		return
	}
	e.OK(req.GetId(), "Delete successful")
}
