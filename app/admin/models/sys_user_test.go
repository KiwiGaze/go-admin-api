package models

import (
	"testing"

	"github.com/go-admin-team/go-admin-core/sdk/pkg"
)

func TestSysUserEncrypt(t *testing.T) {
	user := &SysUser{Password: "secret"}
	if err := user.Encrypt(); err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	if user.Password == "secret" {
		t.Fatal("Encrypt() left plaintext password unchanged")
	}
	if ok, err := pkg.CompareHashAndPassword(user.Password, "secret"); err != nil || !ok {
		t.Fatalf("hashed password did not match original, ok=%v err=%v", ok, err)
	}
}

func TestSysUserBeforeCreateAndBeforeUpdate(t *testing.T) {
	t.Run("before create hashes non-empty password", func(t *testing.T) {
		user := &SysUser{Password: "secret"}
		if err := user.BeforeCreate(nil); err != nil {
			t.Fatalf("BeforeCreate() error = %v", err)
		}
		if ok, err := pkg.CompareHashAndPassword(user.Password, "secret"); err != nil || !ok {
			t.Fatalf("BeforeCreate() hash mismatch, ok=%v err=%v", ok, err)
		}
	})

	t.Run("before update skips hashing when password is empty", func(t *testing.T) {
		user := &SysUser{}
		if err := user.BeforeUpdate(nil); err != nil {
			t.Fatalf("BeforeUpdate() error = %v", err)
		}
		if user.Password != "" {
			t.Fatalf("BeforeUpdate() changed empty password to %q", user.Password)
		}
	})
}

func TestSysUserAfterFind(t *testing.T) {
	user := &SysUser{DeptId: 2, PostId: 3, RoleId: 4}
	if err := user.AfterFind(nil); err != nil {
		t.Fatalf("AfterFind() error = %v", err)
	}
	if len(user.DeptIds) != 1 || user.DeptIds[0] != 2 {
		t.Fatalf("DeptIds = %v, want [2]", user.DeptIds)
	}
	if len(user.PostIds) != 1 || user.PostIds[0] != 3 {
		t.Fatalf("PostIds = %v, want [3]", user.PostIds)
	}
	if len(user.RoleIds) != 1 || user.RoleIds[0] != 4 {
		t.Fatalf("RoleIds = %v, want [4]", user.RoleIds)
	}
}
