package apis

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"gorm.io/gorm"

	"go-admin-api/app/jobs/models"
	"go-admin-api/common/actions"
	commonmodels "go-admin-api/common/models"
	"go-admin-api/internal/testutil"
)

type errorResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

func TestSysJobHandlersHideOutOfScopeJobs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previousEnableDP := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() {
		config.ApplicationConfig.EnableDP = previousEnableDP
	})

	t.Run("start returns not found for hidden job", func(t *testing.T) {
		db := testutil.NewTestDB(t, &models.SysJob{})
		seedHiddenJob(t, db)
		ctx, recorder := testutil.NewGinContext(t, http.MethodGet, "/api/v1/job/start/1", nil, db)
		ctx.Request.Host = "default"
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		ctx.Set(actions.PermissionKey, &actions.DataPermission{DataScope: "5", UserId: 1})

		SysJob{}.StartJobForService(ctx)

		assertResponseCode(t, recorder, http.StatusNotFound)
	})

	t.Run("remove returns not found for hidden job", func(t *testing.T) {
		db := testutil.NewTestDB(t, &models.SysJob{})
		seedHiddenJob(t, db)
		ctx, recorder := testutil.NewGinContext(t, http.MethodGet, "/api/v1/job/remove/1", nil, db)
		ctx.Request.Host = "default"
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}
		ctx.Set(actions.PermissionKey, &actions.DataPermission{DataScope: "5", UserId: 1})

		SysJob{}.RemoveJobForService(ctx)

		assertResponseCode(t, recorder, http.StatusNotFound)
	})
}

func seedHiddenJob(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Create(&models.SysJob{
		JobId:     1,
		JobName:   "hidden",
		Status:    2,
		JobType:   1,
		EntryId:   99,
		ControlBy: commonmodels.ControlBy{CreateBy: 2},
	}).Error; err != nil {
		t.Fatalf("seed hidden job: %v", err)
	}
}

func assertResponseCode(t *testing.T, recorder *httptest.ResponseRecorder, want int) {
	t.Helper()
	var payload errorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if int(payload.Code) != want {
		t.Fatalf("response code = %d, want %d; body = %s", payload.Code, want, recorder.Body.String())
	}
}
