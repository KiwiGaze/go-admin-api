package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/gin-gonic/gin"
	corelogger "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	sdkruntime "github.com/go-admin-team/go-admin-core/sdk/runtime"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-admin-api/common/middleware/handler"
)

const (
	testPassword  = "secret-123"
	testRequestID = "req-demo-123"
)

type authenticatedUser struct {
	User handler.SysUser
	Role handler.SysRole
}

type logEntry struct {
	level   corelogger.Level
	message string
	fields  map[string]string
}

type logStore struct {
	mu      sync.Mutex
	entries []logEntry
}

type captureLogger struct {
	opts   corelogger.Options
	store  *logStore
	fields map[string]interface{}
}

func TestMilestoneDemo(t *testing.T) {
	t.Run("info route responds", func(t *testing.T) {
		app := newMilestoneDemoApp(t)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/info", nil)
		request.Header.Set("X-Request-Id", testRequestID)

		app.router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected /info to return 200, got %d", recorder.Code)
		}

		var response map[string]string
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("decode /info response: %v", err)
		}
		if response["message"] != "ok" {
			t.Fatalf("expected /info message to be ok, got %q", response["message"])
		}
	})

	t.Run("request id is included in logs", func(t *testing.T) {
		app := newMilestoneDemoApp(t)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/info", nil)
		request.Header.Set("X-Request-Id", testRequestID)

		app.router.ServeHTTP(recorder, request)

		if !app.logs.containsFold(pkg.TrafficKey, testRequestID) {
			t.Fatalf("expected logs to contain %s=%q, entries=%v", pkg.TrafficKey, testRequestID, app.logs.snapshot())
		}
	})

	t.Run("login returns a token", func(t *testing.T) {
		app := newMilestoneDemoApp(t)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(
			`{"username":"demo","password":"`+testPassword+`","code":"1234","uuid":"captcha-1"}`,
		))
		request.Header.Set("Content-Type", "application/json")

		app.router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected login to return 200, got %d", recorder.Code)
		}

		var response struct {
			Code  int    `json:"code"`
			Token string `json:"token"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("decode login response: %v", err)
		}
		if response.Code != http.StatusOK {
			t.Fatalf("expected login code 200, got %d", response.Code)
		}
		if response.Token == "" {
			t.Fatal("expected login to return a token")
		}
	})

	t.Run("protected routes reject missing token", func(t *testing.T) {
		app := newMilestoneDemoApp(t)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)

		app.router.ServeHTTP(recorder, request)

		assertJWTRejection(t, recorder, jwtauth.ErrEmptyAuthHeader.Error())
	})

	t.Run("protected routes reject invalid token", func(t *testing.T) {
		app := newMilestoneDemoApp(t)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
		request.Header.Set("Authorization", "Bearer definitely-not-valid")

		app.router.ServeHTTP(recorder, request)

		assertJWTRejection(t, recorder, "token is malformed")
	})

	t.Run("login rejects bad credentials", func(t *testing.T) {
		app := newMilestoneDemoApp(t)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(
			`{"username":"demo","password":"wrong-password","code":"1234","uuid":"captcha-1"}`,
		))
		request.Header.Set("Content-Type", "application/json")

		app.router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected HTTP 200 envelope, got %d", recorder.Code)
		}
		var response struct {
			Code int `json:"code"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("decode login rejection response: %v", err)
		}
		if response.Code != http.StatusBadRequest {
			t.Fatalf("expected code 400, got %d (body: %s)", response.Code, recorder.Body.String())
		}
	})

	t.Run("role enforcement denies unauthorized role", func(t *testing.T) {
		app := newMilestoneDemoApp(t)

		db := sdk.Runtime.GetDbByTenant("*")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("hash password: %v", err)
		}

		viewerRole := handler.SysRole{RoleKey: "viewer", RoleName: "Viewer", Status: "2"}
		if err := db.Create(&viewerRole).Error; err != nil {
			t.Fatalf("insert viewer role: %v", err)
		}
		viewerUser := handler.SysUser{
			Username: "viewer",
			Password: string(hashedPassword),
			NickName: "Viewer User",
			RoleId:   viewerRole.RoleId,
			Status:   "2",
		}
		if err := db.Create(&viewerUser).Error; err != nil {
			t.Fatalf("insert viewer user: %v", err)
		}

		loginRecorder := httptest.NewRecorder()
		loginRequest := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(
			`{"username":"viewer","password":"`+testPassword+`","code":"1234","uuid":"captcha-1"}`,
		))
		loginRequest.Header.Set("Content-Type", "application/json")
		app.router.ServeHTTP(loginRecorder, loginRequest)

		var loginResponse struct {
			Code  int    `json:"code"`
			Token string `json:"token"`
		}
		if err := json.Unmarshal(loginRecorder.Body.Bytes(), &loginResponse); err != nil {
			t.Fatalf("decode viewer login: %v", err)
		}
		if loginResponse.Token == "" {
			t.Fatalf("expected viewer token, body: %s", loginRecorder.Body.String())
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
		request.Header.Set("Authorization", "Bearer "+loginResponse.Token)
		app.router.ServeHTTP(recorder, request)

		var response struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("decode role enforcement response: %v", err)
		}
		if response.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d (body: %s)", response.Code, recorder.Body.String())
		}
	})
}

type milestoneDemoApp struct {
	router *gin.Engine
	logs   *logStore
}

func newMilestoneDemoApp(t *testing.T) milestoneDemoApp {
	t.Helper()

	gin.SetMode(gin.TestMode)
	sdk.Runtime = sdkruntime.NewConfig()

	loggerStore := &logStore{}
	testLogger := newCaptureLogger(loggerStore)
	sdk.Runtime.SetLogger(testLogger)

	config.LoggerConfig.EnabledDB = false
	config.ApplicationConfig.Mode = ""

	db := newMilestoneDemoDB(t)
	sdk.Runtime.SetDbByTenant("*", db)
	sdk.Runtime.SetCasbinByTenant("*", newMilestoneDemoCasbin(t))

	authMiddleware := newMilestoneDemoJWT(t)

	router := gin.New()
	InitMiddleware(router)
	router.GET("/info", handler.PingHandler())
	router.POST("/api/v1/login", authMiddleware.LoginHandler)
	router.GET("/api/v1/protected", authMiddleware.MiddlewareFunc(), AuthCheckRole(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "protected"})
	})

	return milestoneDemoApp{
		router: router,
		logs:   loggerStore,
	}
}

func newMilestoneDemoDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	if err := db.AutoMigrate(&handler.SysRole{}, &handler.SysUser{}); err != nil {
		t.Fatalf("migrate milestone demo tables: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash test password: %v", err)
	}

	role := handler.SysRole{
		RoleKey:  "editor",
		RoleName: "Editor",
		Status:   "2",
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("insert test role: %v", err)
	}

	user := handler.SysUser{
		Username: "demo",
		Password: string(hashedPassword),
		NickName: "Demo User",
		RoleId:   role.RoleId,
		Status:   "2",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("insert test user: %v", err)
	}

	return db
}

func newMilestoneDemoCasbin(t *testing.T) *casbin.SyncedEnforcer {
	t.Helper()

	modelText := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && r.act == p.act
`

	model, err := casbinModel.NewModelFromString(modelText)
	if err != nil {
		t.Fatalf("create casbin model: %v", err)
	}

	enforcer, err := casbin.NewSyncedEnforcer(model)
	if err != nil {
		t.Fatalf("create casbin enforcer: %v", err)
	}
	if _, err := enforcer.AddPolicy("editor", "/api/v1/protected", http.MethodGet); err != nil {
		t.Fatalf("add casbin policy: %v", err)
	}

	return enforcer
}

func newMilestoneDemoJWT(t *testing.T) *jwtauth.GinJWTMiddleware {
	t.Helper()

	middleware, err := jwtauth.New(&jwtauth.GinJWTMiddleware{
		Key:         []byte("milestone-demo-secret"),
		Timeout:     time.Hour,
		TimeFunc:    time.Now,
		IdentityKey: jwtauth.IdentityKey,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginRequest handler.Login
			if err := c.ShouldBindJSON(&loginRequest); err != nil {
				return nil, err
			}

			db := sdk.Runtime.GetDbByTenant(c.Request.Host)
			if db == nil {
				db = sdk.Runtime.GetDbByTenant("*")
			}

			user, role, err := loginRequest.GetUser(db)
			if err != nil {
				return nil, err
			}

			return authenticatedUser{
				User: user,
				Role: role,
			}, nil
		},
		PayloadFunc: func(data interface{}) jwtauth.MapClaims {
			session := data.(authenticatedUser)
			return jwtauth.MapClaims{
				"identity": session.User.UserId,
				"roleid":   session.Role.RoleId,
				"rolekey":  session.Role.RoleKey,
				"rolename": session.Role.RoleName,
			}
		},
	})
	if err != nil {
		t.Fatalf("create jwt middleware: %v", err)
	}

	return middleware
}

func assertJWTRejection(t *testing.T, recorder *httptest.ResponseRecorder, expectedMessage string) {
	t.Helper()

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected JWT rejection response to use HTTP 200 envelope, got %d", recorder.Code)
	}

	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode JWT rejection response: %v", err)
	}
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected JWT rejection code 401, got %d", response.Code)
	}
	if !strings.Contains(response.Message, expectedMessage) {
		t.Fatalf("expected JWT rejection message to contain %q, got %q", expectedMessage, response.Message)
	}
}

func newCaptureLogger(store *logStore) *captureLogger {
	options := corelogger.DefaultOptions()
	options.Level = corelogger.TraceLevel

	return &captureLogger{
		opts:   options,
		store:  store,
		fields: map[string]interface{}{},
	}
}

func (l *captureLogger) Init(options ...corelogger.Option) error {
	for _, option := range options {
		option(&l.opts)
	}
	return nil
}

func (l *captureLogger) Options() corelogger.Options {
	options := l.opts
	options.Fields = cloneInterfaceMap(l.fields)
	return options
}

func (l *captureLogger) Fields(fields map[string]interface{}) corelogger.Logger {
	return &captureLogger{
		opts:   l.opts,
		store:  l.store,
		fields: cloneInterfaceMap(fields),
	}
}

func (l *captureLogger) Log(level corelogger.Level, values ...interface{}) {
	if !l.opts.Level.Enabled(level) {
		return
	}

	l.store.append(logEntry{
		level:   level,
		message: fmt.Sprint(values...),
		fields:  stringifyMap(l.fields),
	})
}

func (l *captureLogger) Logf(level corelogger.Level, format string, values ...interface{}) {
	if !l.opts.Level.Enabled(level) {
		return
	}

	l.store.append(logEntry{
		level:   level,
		message: fmt.Sprintf(format, values...),
		fields:  stringifyMap(l.fields),
	})
}

func (l *captureLogger) String() string {
	return "capture"
}

func (s *logStore) append(entry logEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, entry)
}

func (s *logStore) containsFold(key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range s.entries {
		for entryKey, entryValue := range entry.fields {
			if strings.EqualFold(entryKey, key) && entryValue == value {
				return true
			}
		}
	}

	return false
}

func (s *logStore) snapshot() []logEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	copied := make([]logEntry, 0, len(s.entries))
	for _, entry := range s.entries {
		copied = append(copied, logEntry{
			level:   entry.level,
			message: entry.message,
			fields:  cloneStringMap(entry.fields),
		})
	}

	return copied
}

func cloneInterfaceMap(values map[string]interface{}) map[string]interface{} {
	cloned := make(map[string]interface{}, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}

func stringifyMap(values map[string]interface{}) map[string]string {
	stringified := make(map[string]string, len(values))
	for key, value := range values {
		stringified[key] = fmt.Sprint(value)
	}
	return stringified
}

func cloneStringMap(values map[string]string) map[string]string {
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}
