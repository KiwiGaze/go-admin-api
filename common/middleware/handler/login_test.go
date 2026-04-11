package handler

import (
	"testing"

	"go-admin-api/app/admin/models"
	"go-admin-api/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginGetUser(t *testing.T) {
	db := testutil.NewTestDB(t, &models.SysUser{}, &models.SysRole{})
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	enabledUser := map[string]interface{}{
		"user_id":   1,
		"username":  "enabled",
		"password":  string(passwordHash),
		"role_id":   9,
		"status":    "2",
		"dept_id":   1,
		"post_id":   1,
		"nick_name": "Enabled",
	}
	if err := db.Table("sys_user").Create(enabledUser).Error; err != nil {
		t.Fatalf("seed enabled user: %v", err)
	}
	if err := db.Table("sys_user").Create(map[string]interface{}{
		"user_id":  2,
		"username": "disabled",
		"password": string(passwordHash),
		"role_id":  9,
		"status":   "1",
	}).Error; err != nil {
		t.Fatalf("seed disabled user: %v", err)
	}
	if err := db.Table("sys_user").Create(map[string]interface{}{
		"user_id":  3,
		"username": "norole",
		"password": string(passwordHash),
		"role_id":  99,
		"status":   "2",
	}).Error; err != nil {
		t.Fatalf("seed missing-role user: %v", err)
	}
	if err := db.Table("sys_role").Create(map[string]interface{}{
		"role_id":   9,
		"role_name": "admin",
		"role_key":  "admin",
		"status":    "2",
	}).Error; err != nil {
		t.Fatalf("seed role: %v", err)
	}

	t.Run("returns user and role on valid credentials", func(t *testing.T) {
		login := &Login{Username: "enabled", Password: "secret"}
		user, role, err := login.GetUser(db)
		if err != nil {
			t.Fatalf("GetUser() error = %v", err)
		}
		if user.UserId != 1 {
			t.Fatalf("user id = %d, want 1", user.UserId)
		}
		if role.RoleId != 9 {
			t.Fatalf("role id = %d, want 9", role.RoleId)
		}
	})

	t.Run("rejects disabled users", func(t *testing.T) {
		login := &Login{Username: "disabled", Password: "secret"}
		if _, _, err := login.GetUser(db); err == nil {
			t.Fatal("GetUser() error = nil, want disabled user failure")
		}
	})

	t.Run("rejects bad password", func(t *testing.T) {
		login := &Login{Username: "enabled", Password: "wrong"}
		if _, _, err := login.GetUser(db); err == nil {
			t.Fatal("GetUser() error = nil, want password failure")
		}
	})

	t.Run("rejects users without a role record", func(t *testing.T) {
		login := &Login{Username: "norole", Password: "secret"}
		if _, _, err := login.GetUser(db); err == nil {
			t.Fatal("GetUser() error = nil, want role lookup failure")
		}
	})
}
