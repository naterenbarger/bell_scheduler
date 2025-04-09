package config

import (
	"fmt"

	"bell_scheduler/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB creates a new database connection
func NewDB(dbPath string) (*gorm.DB, error) {
	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbPath), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Schedule{},
		&models.TimeSlot{},
		&models.Settings{},
		&models.LogEntry{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	// Create default admin user if it doesn't exist
	var adminUser models.User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		// Create default admin user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash admin password: %v", err)
		}

		adminUser = models.User{
			Username:            "admin",
			Email:               "admin@example.com",
			Password:            string(hashedPassword),
			Role:                "admin",
			ForcePasswordChange: true,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			return nil, fmt.Errorf("failed to create admin user: %v", err)
		}
	}

	// Create default settings if they don't exist
	var settings models.Settings
	if err := db.First(&settings).Error; err != nil {
		defaultSettings := models.DefaultSettings()
		if err := db.Create(defaultSettings).Error; err != nil {
			return nil, fmt.Errorf("failed to create default settings: %v", err)
		}
	}

	return db, nil
}
