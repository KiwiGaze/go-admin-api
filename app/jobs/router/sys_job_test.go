package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/config"

	common "go-admin-api/common/middleware"
)

func TestJobControlRoutesRequireAuth(t *testing.T) {
	originalSecret := config.JwtConfig.Secret
	config.JwtConfig.Secret = "test-secret"
	t.Cleanup(func() {
		config.JwtConfig.Secret = originalSecret
	})

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	authMiddleware, err := common.AuthInit()
	if err != nil {
		t.Fatalf("AuthInit() error = %v", err)
	}
	registerSysJobRouter(engine.Group("/api/v1"), authMiddleware)

	for _, path := range []string{"/api/v1/job/start/1", "/api/v1/job/remove/1"} {
		t.Run(path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, path, nil)

			engine.ServeHTTP(recorder, request)

			if recorder.Code == http.StatusNotFound {
				t.Fatalf("%s was not registered", path)
			}
			body := recorder.Body.String()
			if !strings.Contains(body, "401") {
				t.Fatalf("response body = %q, want unauthorized response", body)
			}
		})
	}
}
