package models

import (
	"gorm.io/gorm"
)

type Schedule struct {
	BaseModel
	Name        string     `json:"name"`
	Description string     `json:"description"`
	IsDefault   bool       `json:"isDefault"`
	IsTemporary bool       `json:"isTemporary"` // Flag to indicate if this is a temporary schedule that should be reset at end of day
	IsActive    bool       `json:"isActive"`    // Flag to indicate if this is the currently active schedule
	TimeSlots   []TimeSlot `json:"timeSlots" gorm:"foreignKey:ScheduleID;constraint:OnDelete:CASCADE"`
}

type TimeSlot struct {
	BaseModel
	ScheduleID  int64  `json:"scheduleId" gorm:"index"`
	TriggerTime string `json:"triggerTime"`           // HH:MM format
	Days        string `json:"days" gorm:"type:text"` // JSON array of days
	Description string `json:"description"`
}

// BeforeCreate converts the Days array to JSON string
func (ts *TimeSlot) BeforeCreate(tx *gorm.DB) error {
	if ts.Days == "" {
		ts.Days = "[]"
	}
	return nil
}

// AfterFind converts the JSON string back to array
func (ts *TimeSlot) AfterFind(tx *gorm.DB) error {
	if ts.Days == "" {
		ts.Days = "[]"
	}
	return nil
}
