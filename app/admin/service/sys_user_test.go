package service

import (
	"testing"

	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"go-admin-api/app/admin/models"
	servicedto "go-admin-api/app/admin/service/dto"
	"go-admin-api/common/actions"
	commondto "go-admin-api/common/dto"
	"go-admin-api/internal/testutil"
)

func TestSysUserInsertRejectsDuplicateUsername(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysUser{Service: testutil.NewTestService(db)}

	existing := models.SysUser{UserId: 1, Username: "alice", Password: "secret", Status: "2"}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	err := service.Insert(&servicedto.SysUserInsertReq{Username: "alice", Password: "new-secret"})
	if err == nil || err.Error() != "用户名已存在！" {
		t.Fatalf("Insert() error = %v, want duplicate username error", err)
	}
}

func TestSysUserInsertHashesPassword(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysUser{Service: testutil.NewTestService(db)}

	err := service.Insert(&servicedto.SysUserInsertReq{
		UserId:   1,
		Username: "alice",
		Password: "secret",
		NickName: "Alice",
		Phone:    "12345678901",
		RoleId:   9,
		Email:    "alice@example.com",
		DeptId:   1,
		PostId:   2,
		Status:   "2",
	})
	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	var user models.SysUser
	if err := db.Where("username = ?", "alice").First(&user).Error; err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if user.Password == "secret" {
		t.Fatal("stored password remained plaintext")
	}
	if ok, err := pkg.CompareHashAndPassword(user.Password, "secret"); err != nil || !ok {
		t.Fatalf("stored password hash mismatch, ok=%v err=%v", ok, err)
	}
}

func TestSysUserUpdateDoesNotOverwritePasswordOrSalt(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysUser{Service: testutil.NewTestService(db)}

	user := models.SysUser{
		UserId:   1,
		Username: "alice",
		Password: "secret",
		Salt:     "salt-1",
		NickName: "Alice",
		Phone:    "12345678901",
		RoleId:   9,
		Email:    "alice@example.com",
		DeptId:   1,
		PostId:   2,
		Status:   "2",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	var original models.SysUser
	if err := db.First(&original, 1).Error; err != nil {
		t.Fatalf("reload original: %v", err)
	}

	err := service.Update(&servicedto.SysUserUpdateReq{
		UserId:   1,
		Username: "alice",
		NickName: "Updated",
		Phone:    "10987654321",
		RoleId:   10,
		Email:    "updated@example.com",
		DeptId:   2,
		PostId:   3,
		Status:   "1",
	}, &actions.DataPermission{})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	var updated models.SysUser
	if err := db.First(&updated, 1).Error; err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if updated.Password != original.Password {
		t.Fatalf("password changed unexpectedly: got %q want %q", updated.Password, original.Password)
	}
	if updated.Salt != original.Salt {
		t.Fatalf("salt changed unexpectedly: got %q want %q", updated.Salt, original.Salt)
	}
	if updated.NickName != "Updated" {
		t.Fatalf("nick name = %q, want Updated", updated.NickName)
	}
}

func TestSysUserAvatarAndStatusHonorPermission(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysUser{Service: testutil.NewTestService(db)}

	user := models.SysUser{UserId: 1, Username: "alice", Password: "secret", Status: "2", Avatar: "old.png"}
	user.CreateBy = 1
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	permission := &actions.DataPermission{DataScope: "5", UserId: 1}
	if err := service.UpdateAvatar(&servicedto.UpdateSysUserAvatarReq{UserId: 1, Avatar: "new.png"}, permission); err != nil {
		t.Fatalf("UpdateAvatar() error = %v", err)
	}
	if err := service.UpdateStatus(&servicedto.UpdateSysUserStatusReq{UserId: 1, Status: "1"}, permission); err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}

	var updated models.SysUser
	if err := db.First(&updated, 1).Error; err != nil {
		t.Fatalf("reload user: %v", err)
	}
	if updated.Avatar != "new.png" || updated.Status != "1" {
		t.Fatalf("updated user = %+v, want avatar/status change", updated)
	}

	deniedUser := models.SysUser{UserId: 2, Username: "bob", Password: "secret", Status: "2"}
	deniedUser.CreateBy = 2
	if err := db.Create(&deniedUser).Error; err != nil {
		t.Fatalf("seed denied user: %v", err)
	}
	if err := service.UpdateStatus(&servicedto.UpdateSysUserStatusReq{UserId: 2, Status: "1"}, permission); err == nil {
		t.Fatal("UpdateStatus() error = nil, want permission failure")
	}
}

func TestSysUserResetPwdOnlyChangesPassword(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysUser{Service: testutil.NewTestService(db)}

	user := models.SysUser{
		UserId:   1,
		Username: "alice",
		Password: "old-secret",
		NickName: "Alice",
		Phone:    "12345678901",
		RoleId:   9,
		Avatar:   "avatar.png",
		Sex:      "F",
		Status:   "2",
	}
	user.CreateBy = 1
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	var original models.SysUser
	if err := db.First(&original, 1).Error; err != nil {
		t.Fatalf("reload original: %v", err)
	}

	err := service.ResetPwd(&servicedto.ResetSysUserPwdReq{UserId: 1, Password: "new-secret"}, &actions.DataPermission{DataScope: "5", UserId: 1})
	if err != nil {
		t.Fatalf("ResetPwd() error = %v", err)
	}

	var updated models.SysUser
	if err := db.First(&updated, 1).Error; err != nil {
		t.Fatalf("reload updated: %v", err)
	}
	if updated.Username != original.Username || updated.Avatar != original.Avatar || updated.Phone != original.Phone {
		t.Fatalf("non-password fields changed unexpectedly: before=%+v after=%+v", original, updated)
	}
	if updated.Password == original.Password {
		t.Fatal("password did not change")
	}
	if ok, err := pkg.CompareHashAndPassword(updated.Password, "new-secret"); err != nil || !ok {
		t.Fatalf("updated password hash mismatch, ok=%v err=%v", ok, err)
	}
}

func TestSysUserUpdatePwdScenarios(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysUser{Service: testutil.NewTestService(db)}

	user := models.SysUser{UserId: 1, Username: "alice", Password: "old-secret", Status: "2"}
	user.CreateBy = 1
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}

	t.Run("empty new password is a no-op", func(t *testing.T) {
		if err := service.UpdatePwd(1, "old-secret", "", &actions.DataPermission{DataScope: "5", UserId: 1}); err != nil {
			t.Fatalf("UpdatePwd() error = %v", err)
		}
	})

	t.Run("missing scoped record returns permission error", func(t *testing.T) {
		err := service.UpdatePwd(1, "old-secret", "new-secret", &actions.DataPermission{DataScope: "5", UserId: 2})
		if err == nil || err.Error() != "无权更新该数据" {
			t.Fatalf("UpdatePwd() error = %v, want permission error", err)
		}
	})

	t.Run("wrong old password fails", func(t *testing.T) {
		err := service.UpdatePwd(1, "wrong", "new-secret", &actions.DataPermission{DataScope: "5", UserId: 1})
		if err == nil || err.Error() != "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			t.Fatalf("UpdatePwd() error = %v, want bcrypt mismatch", err)
		}
	})

	t.Run("successful password change hashes new password", func(t *testing.T) {
		if err := service.UpdatePwd(1, "old-secret", "new-secret", &actions.DataPermission{DataScope: "5", UserId: 1}); err != nil {
			t.Fatalf("UpdatePwd() error = %v", err)
		}
		var updated models.SysUser
		if err := db.First(&updated, 1).Error; err != nil {
			t.Fatalf("reload user: %v", err)
		}
		if ok, err := pkg.CompareHashAndPassword(updated.Password, "new-secret"); err != nil || !ok {
			t.Fatalf("updated password hash mismatch, ok=%v err=%v", ok, err)
		}
	})
}

func TestSysUserGetProfileLoadsRelations(t *testing.T) {
	setServiceTestGlobals(t)
	db := newServiceTestDB(t)
	service := SysUser{Service: testutil.NewTestService(db)}

	dept := models.SysDept{DeptId: 1, DeptName: "Engineering", Status: 1}
	role := models.SysRole{RoleId: 9, RoleName: "admin", RoleKey: "admin", Status: "2"}
	post := models.SysPost{PostId: 2, PostName: "Developer", Status: 1}
	user := models.SysUser{
		UserId:   1,
		Username: "alice",
		Password: "secret",
		RoleId:   9,
		DeptId:   1,
		PostId:   2,
		Status:   "2",
	}

	for _, seed := range []interface{}{&dept, &role, &post, &user} {
		if err := db.Create(seed).Error; err != nil {
			t.Fatalf("seed relation data: %v", err)
		}
	}

	var loadedUser models.SysUser
	var roles []models.SysRole
	var posts []models.SysPost
	if err := service.GetProfile(&servicedto.SysUserById{ObjectById: commondto.ObjectById{Id: 1}}, &loadedUser, &roles, &posts); err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	if loadedUser.Dept == nil || loadedUser.Dept.DeptName != "Engineering" {
		t.Fatalf("loaded dept = %+v, want Engineering", loadedUser.Dept)
	}
	if len(roles) != 1 || roles[0].RoleKey != "admin" {
		t.Fatalf("loaded roles = %+v, want admin role", roles)
	}
	if len(posts) != 1 || posts[0].PostName != "Developer" {
		t.Fatalf("loaded posts = %+v, want Developer post", posts)
	}
}
