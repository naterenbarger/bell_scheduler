package models

import "time"

// Settings represents the application settings
type Settings struct {
	BaseModel
	RingDuration time.Duration `json:"ringDuration"`
	GPIOPin      int           `json:"gpioPin"`
	Timezone     string        `json:"timezone"`
}

// DefaultSettings returns the default application settings
func DefaultSettings() *Settings {
	return &Settings{
		RingDuration: 5 * time.Second,
		GPIOPin:      17, // Default to GPIO17
		Timezone:     "UTC",
	}
}
