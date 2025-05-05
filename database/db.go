package database

import (
	"log"

	"github.com/vhybZApp/api.git/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize initializes the database connection and performs auto-migration
func Initialize() error {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.AppConfig.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return err
	}

	// Auto-migrate the schema
	if err := AutoMigrate(DB); err != nil {
		log.Fatal("Failed to migrate database:", err)
		return err
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
