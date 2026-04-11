package apis

import (
	"go-admin-api/app/admin/models"
	"go-admin-api/app/admin/service"
	"go-admin-api/app/admin/service/dto"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
)

type SysDept struct {
	api.Api
}

// GetPage retrieves the paginated department list.
// @Summary Department list
// @Description Get the paginated department list
// @Tags Department
// @Param deptName query string false "deptName"
// @Param deptId query string false "deptId"
// @Param parentId query string false "parentId"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dept [get]
// @Security Bearer
func (e SysDept) GetPage(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptGetPageReq{}
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

	list, err := s.SetDeptPage(&req)
	if err != nil {
		e.Error(500, err, "Query failed")
		return
	}

	e.OK(list, "Query succeeded")
}

// Get retrieves department details.
// @Summary Department details
// @Description Get department details by ID
// @Tags Department
// @Param id path string true "id"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/dept/{id} [get]
// @Security Bearer
func (e SysDept) Get(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptGetReq{}
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

	var object models.SysDept
	err = s.Get(&req, &object)
	if err != nil {
		e.Error(500, err, "Query failed")
		return
	}

	e.OK(object, "Query succeeded")
}

// Insert creates a department.
// @Summary Create department
// @Description Create a department
// @Tags Department
// @Accept application/json
// @Product application/json
// @Param data body dto.SysDeptInsertReq true "data"
// @Success 200 {string} string "{"code": 200, "message": "Created successfully"}"
// @Success 200 {string} string "{"code": -1, "message": "Creation failed"}"
// @Router /api/v1/dept [post]
// @Security Bearer
func (e SysDept) Insert(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptInsertReq{}
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

	req.SetCreateBy(user.GetUserId(c))
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, "Creation failed")
		return
	}

	e.OK(req.GetId(), "Created successfully")
}

// Update updates a department.
// @Summary Update department
// @Description Update a department
// @Tags Department
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SysDeptUpdateReq true "body"
// @Success 200 {string} string "{"code": 200, "message": "Updated successfully"}"
// @Success 200 {string} string "{"code": -1, "message": "Update failed"}"
// @Router /api/v1/dept/{id} [put]
// @Security Bearer
func (e SysDept) Update(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptUpdateReq{}
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

	req.SetUpdateBy(user.GetUserId(c))
	err = s.Update(&req)
	if err != nil {
		e.Error(500, err, "Update failed")
		return
	}

	e.OK(req.GetId(), "Updated successfully")
}

// Delete deletes departments.
// @Summary Delete department
// @Description Delete departments by IDs
// @Tags Department
// @Param data body dto.SysDeptDeleteReq true "body"
// @Success 200 {string} string "{"code": 200, "message": "Deleted successfully"}"
// @Success 200 {string} string "{"code": -1, "message": "Delete failed"}"
// @Router /api/v1/dept [delete]
// @Security Bearer
func (e SysDept) Delete(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptDeleteReq{}
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

	err = s.Remove(&req)
	if err != nil {
		e.Error(500, err, "Delete failed")
		return
	}

	e.OK(req.GetId(), "Deleted successfully")
}

// Get2Tree retrieves the department tree used by the user management page.
func (e SysDept) Get2Tree(c *gin.Context) {
	s := service.SysDept{}
	req := dto.SysDeptGetPageReq{}
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

	list, err := s.SetDeptTree(&req)
	if err != nil {
		e.Error(500, err, "Query failed")
		return
	}

	e.OK(list, "")
}

// GetDeptTreeRoleSelect retrieves the department tree and selected keys for a role.
func (e SysDept) GetDeptTreeRoleSelect(c *gin.Context) {
	s := service.SysDept{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	roleID, err := pkg.StringToInt(c.Param("roleId"))
	if err != nil {
		e.Error(500, err, "Invalid role ID")
		return
	}

	depts, err := s.SetDeptLabel()
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	checkedDeptIDs := make([]int, 0)
	if roleID != 0 {
		checkedDeptIDs, err = s.GetWithRoleId(roleID)
		if err != nil {
			e.Error(500, err, err.Error())
			return
		}
	}

	e.OK(gin.H{
		"depts":       depts,
		"checkedKeys": checkedDeptIDs,
	}, "")
}
