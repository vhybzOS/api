package database

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// DBUser represents the user data as stored in the database
type DBUser struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Username     string         `gorm:"uniqueIndex;not null"`
	Password     string         `gorm:"not null"`
	Email        string         `gorm:"uniqueIndex;not null"`
	RefreshToken string         `gorm:"-"` // Not stored in DB, only used temporarily
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *DBUser) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

// DBTokenUsage represents the daily token usage for a user in the database
type DBTokenUsage struct {
	gorm.Model
	UserID uuid.UUID `gorm:"type:uuid;index;foreignKey:ID;references:ID;onDelete:CASCADE"`
	User   DBUser    `gorm:"foreignKey:UserID"`
	Date   time.Time `gorm:"index"`
	Tokens int       `gorm:"default:0"`
}

// DBTokenQuota represents the daily token quota for a user in the database
type DBTokenQuota struct {
	gorm.Model
	UserID     uuid.UUID `gorm:"type:uuid;uniqueIndex;foreignKey:ID;references:ID;onDelete:CASCADE"`
	User       DBUser    `gorm:"foreignKey:UserID"`
	DailyQuota int       `gorm:"default:100000"` // Default 100k tokens per day
}

// HashPassword hashes the password using bcrypt
func (u *DBUser) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword checks if the provided password is correct
func (u *DBUser) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// AutoMigrate performs auto-migration for all database models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&DBUser{},
		&DBTokenUsage{},
		&DBTokenQuota{},
	)
}
