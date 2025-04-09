package models

import (
	"time"
)

// Migration represents a database migration
type Migration struct {
	ID        uint      `gorm:"primarykey"`
	Filename  string    `gorm:"uniqueIndex;not null"`
	AppliedAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
