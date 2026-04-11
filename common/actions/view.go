package actions

import (
	"errors"
	"net/http"

	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"gorm.io/gorm"

	"go-admin-api/common/dto"
	"go-admin-api/common/models"
)

// ViewAction returns a generic detail handler.
func ViewAction(control dto.Control, f func() interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := pkg.GenerateMsgIDFromContext(c)
		// View details.
		req := control.Generate()
		err = req.Bind(c)
		if err != nil {
			response.Error(c, http.StatusUnprocessableEntity, err, "Parameter validation failed")
			return
		}
		var object models.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			response.Error(c, 500, err, "Model generation failed")
			return
		}

		var rsp interface{}
		if f != nil {
			rsp = f()
		} else {
			rsp, _ = req.GenerateM()
		}

		// Data permission check.
		p := GetPermissionFromContext(c)

		err = db.Model(object).WithContext(c).Scopes(
			Permission(object.TableName(), p),
		).Where(req.GetId()).First(rsp).Error

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, nil, "The requested object does not exist or you do not have permission to view it")
			return
		}
		if err != nil {
			log.Errorf("MsgID[%s] View error: %s", msgID, err)
			response.Error(c, 500, err, "View failed")
			return
		}
		response.OK(c, rsp, "Query succeeded")
		c.Next()
	}
}
