package store

import (
	"bell_scheduler/internal/config"

	"gorm.io/gorm"
)

// NewDB creates a new database connection
func NewDB(dbPath string) (*gorm.DB, error) {
	return config.NewDB(dbPath)
}
