package database

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	Password     string `gorm:"not null" json:"-"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	RefreshToken string `gorm:"-" json:"-"` // Not stored in DB, only used temporarily
}

// HashPassword hashes the password using bcrypt
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword checks if the provided password is correct
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
