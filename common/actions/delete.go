package actions

import (
	"go-admin-api/common/dto"
	"go-admin-api/common/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
)

func DeleteAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := pkg.GenerateMsgIDFromContext(c)
		req := control.Generate()

		err = req.Bind(c)
		if err != nil {
			log.Errorf("MsgID[%s] DeleteAction parameter binding error: %s", msgID, err)
			response.Error(c, http.StatusUnprocessableEntity, err, "Invalid parameters")
			return
		}
		var object models.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			log.Errorf("MsgID[%s] DeleteAction model generation error: %s", msgID, err)
			response.Error(c, 500, err, "Model generation failed")
			return
		}

		object.SetUpdateBy(user.GetUserId(c))
		
		// data authorization check
		p := GetPermissionFromContext(c)
		db = db.WithContext(c).Scopes(
			Permission(object.TableName(), p),
		).Where(req.GetId()).Delete(object)
		if err=db.Error; err != nil {
			log.Errorf("MsgID[%s] DeleteAction delete error: %s", msgID, err)
			response.Error(c, 500, err, "Delete failed")
			return
		}
		if db.RowsAffected == 0 {
			log.Warnf("MsgID[%s] DeleteAction delete failed: no rows affected", msgID)
			response.Error(c, http.StatusForbidden, nil, "Record not found or no permission to delete")
			return
		}
		response.OK(c, object.GetId(), "Delete succeeded.")
		c.Next()
	}
}