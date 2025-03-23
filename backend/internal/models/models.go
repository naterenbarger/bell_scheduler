package models

import "time"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"` // Password is never sent to the client
	Role     string `json:"role"`
}

type Schedule struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsDefault   bool      `json:"isDefault"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Times       []TimeSlot `json:"times" gorm:"foreignKey:ScheduleID"`
}

type TimeSlot struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ScheduleID  uint      `json:"scheduleId"`
	TriggerTime string    `json:"triggerTime"` // HH:MM format
	Days        []string  `json:"days"`        // ["Monday", "Tuesday", etc.]
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Settings struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	BellRingDuration int   `json:"bellRingDuration"` // Duration in seconds
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
} 