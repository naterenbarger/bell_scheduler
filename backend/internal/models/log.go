package models

import "time"

// LogEntry represents a bell ringing event log
type LogEntry struct {
    ID        int64     `json:"id" gorm:"primaryKey"`
    Timestamp time.Time `json:"timestamp"`
    Trigger   string    `json:"trigger"` // "schedule" or "manual"
    UserID    int64     `json:"userId,omitempty"`
    Username  string    `json:"username,omitempty"`
    ScheduleID int64    `json:"scheduleId,omitempty"`
    ScheduleName string `json:"scheduleName,omitempty"`
    ScheduleTime string `json:"scheduleTime,omitempty"`
    CreatedAt time.Time `json:"createdAt"`
}

// TableName specifies the table name for LogEntry
func (LogEntry) TableName() string {
    return "log_entries"
} 