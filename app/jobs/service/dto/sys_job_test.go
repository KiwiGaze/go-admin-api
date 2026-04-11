package dto

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSysJobDeleteReqBind(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("accepts ids from delete body", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/sysjob", bytes.NewBufferString(`{"ids":[3,4]}`))
		ctx.Request.Header.Set("Content-Type", "application/json")

		req := SysJobDeleteReq{}
		if err := req.Bind(ctx); err != nil {
			t.Fatalf("Bind() error = %v", err)
		}
		if len(req.Ids) != 2 || req.Ids[0] != 3 || req.Ids[1] != 4 {
			t.Fatalf("ids = %v, want [3 4]", req.Ids)
		}
	})

	t.Run("normalizes single id into ids", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/sysjob", bytes.NewBufferString(`{"id":7}`))
		ctx.Request.Header.Set("Content-Type", "application/json")

		req := SysJobDeleteReq{}
		if err := req.Bind(ctx); err != nil {
			t.Fatalf("Bind() error = %v", err)
		}
		if len(req.Ids) != 1 || req.Ids[0] != 7 {
			t.Fatalf("ids = %v, want [7]", req.Ids)
		}
	})

	t.Run("rejects empty delete body", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/sysjob", bytes.NewBufferString(`{}`))
		ctx.Request.Header.Set("Content-Type", "application/json")

		req := SysJobDeleteReq{}
		if err := req.Bind(ctx); err == nil {
			t.Fatal("Bind() error = nil, want validation error")
		}
	})
}
