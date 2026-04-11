package handler

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	sdkcaptcha "github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"github.com/mojocn/base64Captcha"
	adminmodels "go-admin-api/app/admin/models"
	"go-admin-api/internal/testutil"
)

func TestPayloadFunc(t *testing.T) {
	claims := PayloadFunc(map[string]interface{}{
		"user": SysUser{UserId: 7, Username: "alice"},
		"role": SysRole{RoleId: 8, RoleKey: "admin", RoleName: "Administrator", DataScope: "4"},
	})

	if claims[jwt.IdentityKey] != 7 {
		t.Fatalf("identity claim = %v, want 7", claims[jwt.IdentityKey])
	}
	if claims[jwt.RoleKey] != "admin" {
		t.Fatalf("role key claim = %v, want admin", claims[jwt.RoleKey])
	}
	if claims[jwt.RoleNameKey] != "Administrator" {
		t.Fatalf("role name claim = %v, want Administrator", claims[jwt.RoleNameKey])
	}
}

func TestIdentityHandler(t *testing.T) {
	ctx, _ := testutil.NewGinContext(t, "GET", "/", nil, nil)
	ctx.Set(jwt.JwtPayloadKey, jwt.MapClaims{
		"identity":  1,
		"nice":      "alice",
		"rolekey":   "admin",
		"roleid":    2,
		"datascope": "4",
	})

	identity := IdentityHandler(ctx).(map[string]interface{})
	if identity["UserId"] != 1 {
		t.Fatalf("UserId = %v, want 1", identity["UserId"])
	}
	if identity["UserName"] != "alice" {
		t.Fatalf("UserName = %v, want alice", identity["UserName"])
	}
	if identity["DataScope"] != "4" {
		t.Fatalf("DataScope = %v, want 4", identity["DataScope"])
	}
}

func TestAuthorizator(t *testing.T) {
	ctx, _ := testutil.NewGinContext(t, "GET", "/", nil, nil)
	ok := Authorizator(map[string]interface{}{
		"user": adminmodels.SysUser{UserId: 3, Username: "alice"},
		"role": adminmodels.SysRole{RoleId: 4, RoleName: "admin", DataScope: "5"},
	}, ctx)
	if !ok {
		t.Fatal("Authorizator() = false, want true")
	}
	if got := ctx.GetInt("userId"); got != 3 {
		t.Fatalf("context userId = %d, want 3", got)
	}
	if got := ctx.GetString("role"); got != "admin" {
		t.Fatalf("context role = %q, want admin", got)
	}
	if got := ctx.GetString("dataScope"); got != "5" {
		t.Fatalf("context dataScope = %q, want 5", got)
	}
	if Authorizator("invalid", ctx) {
		t.Fatal("Authorizator() = true for invalid payload, want false")
	}
}

func TestUnauthorized(t *testing.T) {
	ctx, recorder := testutil.NewGinContext(t, "GET", "/", nil, nil)
	Unauthorized(ctx, 401, "denied")
	if recorder.Code != 200 {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["code"] != float64(401) || body["msg"] != "denied" {
		t.Fatalf("response = %+v, want code 401 and msg denied", body)
	}
}

func TestAuthenticator(t *testing.T) {
	originalMode := config.ApplicationConfig.Mode
	originalEnabledDB := config.LoggerConfig.EnabledDB
	t.Cleanup(func() {
		config.ApplicationConfig.Mode = originalMode
		config.LoggerConfig.EnabledDB = originalEnabledDB
	})
	config.LoggerConfig.EnabledDB = false

	db := testutil.NewTestDB(t, &adminmodels.SysUser{}, &adminmodels.SysRole{})
	if err := db.Create(&adminmodels.SysRole{
		RoleId:   9,
		RoleName: "admin",
		RoleKey:  "admin",
		Status:   "2",
	}).Error; err != nil {
		t.Fatalf("seed role: %v", err)
	}
	if err := db.Create(&adminmodels.SysUser{
		UserId:   1,
		Username: "alice",
		Password: "secret",
		RoleId:   9,
		Status:   "2",
		DeptId:   1,
		PostId:   1,
		NickName: "Alice",
	}).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	t.Run("returns missing values error when bind fails", func(t *testing.T) {
		config.ApplicationConfig.Mode = "dev"
		ctx, _ := testutil.NewGinContext(t, "POST", "/api/v1/login", strings.NewReader(`{}`), db)
		if _, err := Authenticator(ctx); err != jwt.ErrMissingLoginValues {
			t.Fatalf("Authenticator() error = %v, want %v", err, jwt.ErrMissingLoginValues)
		}
	})

	t.Run("rejects invalid captcha outside dev mode", func(t *testing.T) {
		config.ApplicationConfig.Mode = "prod"
		ctx, _ := testutil.NewGinContext(t, "POST", "/api/v1/login", strings.NewReader(`{"username":"alice","password":"secret","uuid":"captcha-1","code":"1234"}`), db)
		if _, err := Authenticator(ctx); err != jwt.ErrInvalidVerificationode {
			t.Fatalf("Authenticator() error = %v, want %v", err, jwt.ErrInvalidVerificationode)
		}
	})

	t.Run("bypasses captcha in dev mode and returns user payload", func(t *testing.T) {
		config.ApplicationConfig.Mode = "dev"
		ctx, _ := testutil.NewGinContext(t, "POST", "/api/v1/login", strings.NewReader(`{"username":"alice","password":"secret","uuid":"unused","code":"unused"}`), db)
		data, err := Authenticator(ctx)
		if err != nil {
			t.Fatalf("Authenticator() error = %v", err)
		}
		payload, ok := data.(map[string]interface{})
		if !ok {
			t.Fatalf("payload type = %T, want map[string]interface{}", data)
		}
		if payload["user"].(SysUser).Username != "alice" {
			t.Fatalf("payload user = %+v, want alice", payload["user"])
		}
	})

	t.Run("returns failed authentication for wrong password", func(t *testing.T) {
		config.ApplicationConfig.Mode = "dev"
		ctx, _ := testutil.NewGinContext(t, "POST", "/api/v1/login", strings.NewReader(`{"username":"alice","password":"wrong","uuid":"unused","code":"unused"}`), db)
		if _, err := Authenticator(ctx); err != jwt.ErrFailedAuthentication {
			t.Fatalf("Authenticator() error = %v, want %v", err, jwt.ErrFailedAuthentication)
		}
	})

	t.Run("accepts valid captcha outside dev mode", func(t *testing.T) {
		config.ApplicationConfig.Mode = "prod"
		sdkcaptcha.SetStore(base64Captcha.DefaultMemStore)
		if err := base64Captcha.DefaultMemStore.Set("captcha-2", "4321"); err != nil {
			t.Fatalf("seed captcha: %v", err)
		}
		ctx, _ := testutil.NewGinContext(t, "POST", "/api/v1/login", strings.NewReader(`{"username":"alice","password":"secret","uuid":"captcha-2","code":"4321"}`), db)
		data, err := Authenticator(ctx)
		if err != nil {
			t.Fatalf("Authenticator() error = %v", err)
		}
		if payload := data.(map[string]interface{}); payload["role"].(SysRole).RoleKey != "admin" {
			t.Fatalf("payload role = %+v, want admin role", payload["role"])
		}
	})
}
