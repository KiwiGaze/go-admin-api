package models

import (
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"gorm.io/gorm"
)

// BaseUser defines the basic fields for a user, including username, salt, password hash, and password.
type BaseUser struct {
	Username     string `json:"username" gorm:"type:varchar(256);comment:Username"`
	Salt         string `json:"-" gorm:"type:varchar(256);comment:Salt;<-"`
	PasswordHash string `json:"-" gorm:"type:varchar(128);comment:Password hash;<-"`
	Password     string `json:"password" gorm:"-"`
}

// Set user's password
func (u *BaseUser) SetPassword(password string) {
	u.Password = password
	u.generateSalt() // Generate a random salt
	u.PasswordHash = u.GetPasswordHash() // Hash the password with the salt
}

// GetPasswordHash generates the password hash using the password and salt
func (u *BaseUser) GetPasswordHash() string {
	passwordHash, err := pkg.SetPassword(u.Password, u.Salt)
	if err != nil {
		return ""
	}
	return passwordHash
}

// generateSalt generates a random salt for the user
func (u *BaseUser) generateSalt() {
	u.Salt = pkg.GenerateRandomKey16()
}

// Verify user's password
func (u *BaseUser) VerifyPassword(db *gorm.DB, tableName string) bool {
	db.Table(tableName).Where("username = ?", u.Username).First(u)
	return u.GetPasswordHash() == u.PasswordHash
}