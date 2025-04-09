package models

import "time"

// BaseModel is a custom base model with int64 ID
type BaseModel struct {
    ID        int64     `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    DeletedAt time.Time `json:"-" gorm:"index"`
} 