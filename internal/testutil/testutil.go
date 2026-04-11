package testutil

import (
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	corelog "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	corepkg "github.com/go-admin-team/go-admin-core/sdk/pkg"
	coreservice "github.com/go-admin-team/go-admin-core/sdk/service"
	corecache "github.com/go-admin-team/go-admin-core/storage/cache"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewTestLogger() *corelog.Helper {
	return corelog.NewHelper(corelog.NewLogger(corelog.WithLevel(corelog.ErrorLevel)))
}

func NewTestDB(t *testing.T, models ...interface{}) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s-%d?mode=memory&cache=shared", sanitizeName(t.Name()), time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			t.Fatalf("auto migrate: %v", err)
		}
	}

	return db
}

func NewTestService(db *gorm.DB) coreservice.Service {
	return coreservice.Service{
		Orm: db,
		Log: NewTestLogger(),
	}
}

func NewGinContext(t *testing.T, method, target string, body io.Reader, db *gorm.DB) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(method, target, body)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	if db != nil {
		ctx.Set("db", db)
	}
	ctx.Set(corepkg.LoggerKey, NewTestLogger())
	return ctx, recorder
}

func UseRuntimeDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	previous := sdk.Runtime.GetDbByTenant("default")
	sdk.Runtime.SetDbByTenant("default", db)
	t.Cleanup(func() {
		sdk.Runtime.SetDbByTenant("default", previous)
	})
}

func UseRuntimeLogger(t *testing.T) {
	t.Helper()

	previous := sdk.Runtime.GetLogger()
	logger := corelog.NewLogger(corelog.WithLevel(corelog.ErrorLevel))
	sdk.Runtime.SetLogger(logger)
	t.Cleanup(func() {
		sdk.Runtime.SetLogger(previous)
	})
}

func UseRuntimeCache(t *testing.T) {
	t.Helper()

	previous := sdk.Runtime.GetCacheAdapter()
	sdk.Runtime.SetCacheAdapter(corecache.NewMemory())
	t.Cleanup(func() {
		sdk.Runtime.SetCacheAdapter(previous)
	})
}

func sanitizeName(name string) string {
	replacer := strings.NewReplacer("/", "_", " ", "_")
	return replacer.Replace(name)
}
