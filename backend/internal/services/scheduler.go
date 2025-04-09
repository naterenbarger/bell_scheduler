package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"bell_scheduler/internal/models"
	"bell_scheduler/internal/store"
)

// SchedulerService manages the bell schedules and triggers
type SchedulerService struct {
	gpio         *GPIOService
	schedules    []models.Schedule
	logRepo      *store.LogRepository
	scheduleRepo *store.ScheduleRepository
	mu           sync.RWMutex
	stopChan     chan struct{}
}

// NewSchedulerService creates a new scheduler service instance
func NewSchedulerService(gpio *GPIOService, logRepo *store.LogRepository, scheduleRepo *store.ScheduleRepository) *SchedulerService {
	return &SchedulerService{
		gpio:         gpio,
		schedules:    make([]models.Schedule, 0),
		logRepo:      logRepo,
		scheduleRepo: scheduleRepo,
		stopChan:     make(chan struct{}),
	}
}

// Start begins the scheduler service
func (s *SchedulerService) Start() {
	go s.run()
}

// Stop gracefully stops the scheduler service
func (s *SchedulerService) Stop() {
	close(s.stopChan)
}

// run is the main scheduler loop
func (s *SchedulerService) run() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.checkSchedules()
			s.checkForScheduleReset()
		}
	}
}

// checkSchedules checks if any schedules should be triggered
func (s *SchedulerService) checkSchedules() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	currentTime := now.Format("15:04")
	currentDay := now.Weekday().String()

	// Find the active schedule
	var activeSchedule *models.Schedule
	for _, schedule := range s.schedules {
		if schedule.IsActive {
			schedule := schedule // Create a copy to avoid pointer issues
			activeSchedule = &schedule
			break
		}
	}

	// If no active schedule is found, check if there's a default schedule
	if activeSchedule == nil {
		for _, schedule := range s.schedules {
			if schedule.IsDefault {
				schedule := schedule // Create a copy to avoid pointer issues
				activeSchedule = &schedule
				break
			}
		}
	}

	// If still no schedule is found, return without triggering anything
	if activeSchedule == nil {
		return
	}

	// Only check timeslots for the active schedule
	for _, timeSlot := range activeSchedule.TimeSlots {
		if timeSlot.TriggerTime == currentTime {
			var days []string
			if err := json.Unmarshal([]byte(timeSlot.Days), &days); err != nil {
				fmt.Printf("Failed to parse days for schedule %d: %v\n", activeSchedule.ID, err)
				continue
			}
			for _, day := range days {
				if day == currentDay {
					if err := s.triggerSchedule(*activeSchedule, timeSlot); err != nil {
						fmt.Printf("Failed to trigger schedule %d: %v\n", activeSchedule.ID, err)
					}
					break
				}
			}
		}
	}
}

// triggerSchedule triggers a specific schedule and logs the event
func (s *SchedulerService) triggerSchedule(schedule models.Schedule, timeSlot models.TimeSlot) error {
	if err := s.gpio.Trigger(); err != nil {
		return fmt.Errorf("failed to trigger GPIO: %w", err)
	}

	logEntry := &models.LogEntry{
		Timestamp:    time.Now(),
		Trigger:      "schedule",
		ScheduleID:   schedule.ID,
		ScheduleName: schedule.Name,
		ScheduleTime: timeSlot.TriggerTime,
	}

	if err := s.logRepo.Create(logEntry); err != nil {
		return fmt.Errorf("failed to create log entry: %w", err)
	}

	return nil
}

// UpdateSchedules updates the list of active schedules
func (s *SchedulerService) UpdateSchedules(schedules []models.Schedule) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.schedules = schedules
}

// GetSchedules returns the current list of schedules
func (s *SchedulerService) GetSchedules() []models.Schedule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.schedules
}

// TriggerNow manually triggers the bell
func (s *SchedulerService) TriggerNow(userID int64, username string) error {
	if err := s.gpio.Trigger(); err != nil {
		return fmt.Errorf("failed to trigger GPIO: %w", err)
	}

	// Get the default schedule or first available schedule
	s.mu.RLock()
	var schedule models.Schedule
	var scheduleFound bool

	// First try to find the default schedule
	for _, s := range s.schedules {
		if s.IsDefault {
			schedule = s
			scheduleFound = true
			break
		}
	}

	// If no default schedule, use the first available schedule
	if !scheduleFound && len(s.schedules) > 0 {
		schedule = s.schedules[0]
		scheduleFound = true
	}
	s.mu.RUnlock()

	logEntry := &models.LogEntry{
		Timestamp: time.Now(),
		Trigger:   "manual",
		UserID:    userID,
		Username:  username,
	}

	// Include schedule information if a schedule was found
	if scheduleFound {
		logEntry.ScheduleID = schedule.ID
		logEntry.ScheduleName = schedule.Name
		// For manual triggers, we don't have a specific time slot, so we'll use the current time
		logEntry.ScheduleTime = time.Now().Format("15:04")
	}

	if err := s.logRepo.Create(logEntry); err != nil {
		return fmt.Errorf("failed to create log entry: %w", err)
	}

	return nil
}

// SetDuration updates the bell ring duration
func (s *SchedulerService) SetDuration(duration time.Duration) {
	s.gpio.SetDuration(duration)
}

// IsActive returns whether the bell is currently ringing
func (s *SchedulerService) IsActive() bool {
	return s.gpio.IsActive()
}

// checkForScheduleReset checks if it's midnight and resets any temporary schedules
func (s *SchedulerService) checkForScheduleReset() {
	now := time.Now()
	// Check if it's midnight (00:00)
	if now.Hour() == 0 && now.Minute() == 0 {
		// Get all schedules
		schedules, err := s.scheduleRepo.GetAll()
		if err != nil {
			fmt.Printf("Failed to get schedules for reset check: %v\n", err)
			return
		}

		// Find the default schedule
		var defaultSchedule *models.Schedule
		for _, s := range schedules {
			if s.IsDefault {
				defaultSchedule = &s
				break
			}
		}

		if defaultSchedule == nil {
			fmt.Println("No default schedule found for reset")
			return
		}

		// Check for temporary schedules that need to be reset
		for _, schedule := range schedules {
			if schedule.IsTemporary && schedule.IsDefault {
				// Reset this schedule and set the default schedule as active
				if err := s.scheduleRepo.SetDefault(defaultSchedule.ID); err != nil {
					fmt.Printf("Failed to reset temporary schedule %d: %v\n", schedule.ID, err)
					continue
				}

				// Update the schedule to no longer be temporary
				schedule.IsTemporary = false
				schedule.IsDefault = false
				if err := s.scheduleRepo.Update(&schedule); err != nil {
					fmt.Printf("Failed to update temporary schedule %d: %v\n", schedule.ID, err)
				}

				fmt.Printf("Reset temporary schedule %s (ID: %d) to default schedule %s (ID: %d)\n",
					schedule.Name, schedule.ID, defaultSchedule.Name, defaultSchedule.ID)
			}
		}

		// Update the scheduler with the updated schedules
		updatedSchedules, err := s.scheduleRepo.GetAll()
		if err != nil {
			fmt.Printf("Failed to get updated schedules after reset: %v\n", err)
			return
		}
		s.UpdateSchedules(updatedSchedules)
	}
}
