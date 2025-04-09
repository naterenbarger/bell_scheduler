package store

import (
	"bell_scheduler/internal/models"
	"time"

	"gorm.io/gorm"
)

// LogRepository handles database operations for log entries
type LogRepository struct {
	db *gorm.DB
}

// NewLogRepository creates a new log repository instance
func NewLogRepository(db *gorm.DB) *LogRepository {
	return &LogRepository{db: db}
}

// Create adds a new log entry
func (r *LogRepository) Create(log *models.LogEntry) error {
	return r.db.Create(log).Error
}

// GetAll retrieves all log entries, ordered by timestamp descending
func (r *LogRepository) GetAll() ([]models.LogEntry, error) {
	var logs []models.LogEntry
	err := r.db.Order("timestamp DESC").Find(&logs).Error
	return logs, err
}

// GetByDateRange retrieves log entries within a date range
func (r *LogRepository) GetByDateRange(start, end time.Time) ([]models.LogEntry, error) {
	var logs []models.LogEntry
	err := r.db.Where("timestamp BETWEEN ? AND ?", start, end).
		Order("timestamp DESC").
		Find(&logs).Error
	return logs, err
}

// DeleteOldLogs removes log entries older than the specified duration
func (r *LogRepository) DeleteOldLogs(olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	return r.db.Where("timestamp < ?", cutoff).Delete(&models.LogEntry{}).Error
}
