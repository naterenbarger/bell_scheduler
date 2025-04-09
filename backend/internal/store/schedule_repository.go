package store

import (
	"bell_scheduler/internal/models"

	"gorm.io/gorm"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{
		db: db,
	}
}

func (r *ScheduleRepository) Create(schedule *models.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r *ScheduleRepository) Get(id int64) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := r.db.Preload("TimeSlots").First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepository) GetAll() ([]models.Schedule, error) {
	var schedules []models.Schedule
	if err := r.db.Preload("TimeSlots").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *ScheduleRepository) Update(schedule *models.Schedule) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		println("Starting update for schedule:", schedule.ID)
		println("Number of time slots:", len(schedule.TimeSlots))

		// First, get the existing schedule with its time slots
		var existingSchedule models.Schedule
		if err := tx.Preload("TimeSlots").First(&existingSchedule, schedule.ID).Error; err != nil {
			println("Error fetching existing schedule:", err.Error())
			return err
		}

		// Delete all existing time slots
		if err := tx.Where("schedule_id = ?", schedule.ID).Delete(&models.TimeSlot{}).Error; err != nil {
			println("Error deleting existing time slots:", err.Error())
			return err
		}

		// Create new time slots
		for i := range schedule.TimeSlots {
			// Create a copy to avoid modifying the original slice
			slot := schedule.TimeSlots[i]
			slot.ScheduleID = schedule.ID
			// Reset ID to 0 for new slots (those with ID <= 0)
			// This ensures the database will assign a new ID
			if slot.ID <= 0 {
				slot.ID = 0
			}
			println("Creating time slot with trigger time:", slot.TriggerTime)
			if err := tx.Create(&slot).Error; err != nil {
				println("Error creating time slot:", err.Error())
				return err
			}
			// Update the ID in the original slice
			schedule.TimeSlots[i].ID = slot.ID
			println("Created time slot with ID:", slot.ID)
		}

		// Update only the basic schedule fields
		updates := map[string]interface{}{
			"name":         schedule.Name,
			"description":  schedule.Description,
			"is_default":   schedule.IsDefault,
			"is_temporary": schedule.IsTemporary,
		}
		if err := tx.Model(&models.Schedule{}).Where("id = ?", schedule.ID).Updates(updates).Error; err != nil {
			println("Error updating schedule:", err.Error())
			return err
		}

		// Reload the schedule with updated time slots
		if err := tx.Preload("TimeSlots").First(schedule, schedule.ID).Error; err != nil {
			println("Error reloading schedule:", err.Error())
			return err
		}

		println("Update completed successfully")
		println("Final number of time slots:", len(schedule.TimeSlots))
		for _, slot := range schedule.TimeSlots {
			println("Time slot triggerTime:", slot.TriggerTime)
		}
		return nil
	})
}

func (r *ScheduleRepository) Delete(id int64) error {
	return r.db.Delete(&models.Schedule{}, id).Error
}

func (r *ScheduleRepository) FindByID(id int64) (*models.Schedule, error) {
	var schedule models.Schedule
	err := r.db.Preload("TimeSlots").First(&schedule, id).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepository) List() ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Preload("TimeSlots").Find(&schedules).Error
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *ScheduleRepository) GetDefault() (*models.Schedule, error) {
	var schedule models.Schedule
	err := r.db.Preload("TimeSlots").Where("is_default = ?", true).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepository) SetDefault(id int64) error {
	// First, unset any existing default schedule
	if err := r.db.Model(&models.Schedule{}).Where("id != ?", id).Update("is_default", false).Error; err != nil {
		return err
	}
	// Then set the new default schedule without affecting is_active status
	return r.db.Model(&models.Schedule{}).Where("id = ?", id).Update("is_default", true).Error
}

// SetActive sets a schedule as the active schedule
func (r *ScheduleRepository) SetActive(id int64) error {
	// First, unset any existing active schedule
	if err := r.db.Model(&models.Schedule{}).Where("id != ?", id).Update("is_active", false).Error; err != nil {
		return err
	}
	// Then set the new active schedule without affecting is_default status
	return r.db.Model(&models.Schedule{}).Where("id = ?", id).Update("is_active", true).Error
}

// TimeSlot operations
func (r *ScheduleRepository) AddTimeSlot(scheduleID int64, timeSlot *models.TimeSlot) error {
	timeSlot.ScheduleID = scheduleID
	return r.db.Create(timeSlot).Error
}

func (r *ScheduleRepository) UpdateTimeSlot(timeSlot *models.TimeSlot) error {
	return r.db.Save(timeSlot).Error
}

func (r *ScheduleRepository) DeleteTimeSlot(id int64) error {
	return r.db.Delete(&models.TimeSlot{}, id).Error
}

func (r *ScheduleRepository) GetTimeSlot(id int64) (*models.TimeSlot, error) {
	var timeSlot models.TimeSlot
	err := r.db.First(&timeSlot, id).Error
	if err != nil {
		return nil, err
	}
	return &timeSlot, nil
}
