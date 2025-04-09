package store

import (
	"bell_scheduler/internal/models"

	"gorm.io/gorm"
)

type SettingsRepository struct {
	db *gorm.DB
}

func NewSettingsRepository(db *gorm.DB) *SettingsRepository {
	return &SettingsRepository{
		db: db,
	}
}

func (r *SettingsRepository) Get() (*models.Settings, error) {
	var settings models.Settings
	if err := r.db.First(&settings).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *SettingsRepository) Update(settings *models.Settings) error {
	return r.db.Save(settings).Error
}

func (r *SettingsRepository) UpdateBellRingDuration(duration int) error {
	return r.db.Model(&models.Settings{}).Update("bell_ring_duration", duration).Error
}
