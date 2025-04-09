package handlers

import (
	"net/http"
	"strconv"

	"bell_scheduler/internal/models"
	"bell_scheduler/internal/services"
	"bell_scheduler/internal/store"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	scheduleRepo *store.ScheduleRepository
	scheduler    *services.SchedulerService
}

func NewScheduleHandler(scheduleRepo *store.ScheduleRepository, scheduler *services.SchedulerService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleRepo: scheduleRepo,
		scheduler:    scheduler,
	}
}

// GetAll returns all schedules
func (h *ScheduleHandler) GetAll(c *gin.Context) {
	schedules, err := h.scheduleRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get schedules"})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

// Get returns a specific schedule
func (h *ScheduleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	schedule, err := h.scheduleRepo.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// Create creates a new schedule
func (h *ScheduleHandler) Create(c *gin.Context) {
	var req models.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the request data
	for _, slot := range req.TimeSlots {
		println("Creating time slot with triggerTime:", slot.TriggerTime)
	}

	schedule := &models.Schedule{
		Name:        req.Name,
		Description: req.Description,
		IsDefault:   req.IsDefault,
		IsTemporary: req.IsTemporary,
		TimeSlots:   req.TimeSlots,
	}

	if err := h.scheduleRepo.Create(schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
		return
	}

	// Log the created schedule
	println("Created schedule with time slots:")
	for _, slot := range schedule.TimeSlots {
		println("Time slot triggerTime:", slot.TriggerTime)
	}

	// Update scheduler with new schedule
	schedules, err := h.scheduleRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update scheduler"})
		return
	}
	h.scheduler.UpdateSchedules(schedules)

	c.JSON(http.StatusCreated, schedule)
}

// Update updates an existing schedule
func (h *ScheduleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	var req models.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the request data
	for _, slot := range req.TimeSlots {
		println("Updating time slot with triggerTime:", slot.TriggerTime)
	}

	schedule, err := h.scheduleRepo.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	// Update basic schedule fields
	schedule.Name = req.Name
	schedule.Description = req.Description
	schedule.IsDefault = req.IsDefault
	schedule.IsTemporary = req.IsTemporary

	// Handle time slots update
	// First, delete time slots that are not in the request
	existingSlotIDs := make(map[int64]bool)
	for _, slot := range schedule.TimeSlots {
		existingSlotIDs[slot.ID] = true
	}

	// Keep track of which slots we've processed
	processedSlotIDs := make(map[int64]bool)

	// Update or create time slots
	for _, newSlot := range req.TimeSlots {
		if newSlot.ID > 0 && existingSlotIDs[newSlot.ID] {
			// Update existing slot
			for i, existingSlot := range schedule.TimeSlots {
				if existingSlot.ID == newSlot.ID {
					// Preserve the ID and ScheduleID
					newSlot.ScheduleID = schedule.ID
					schedule.TimeSlots[i] = newSlot
					processedSlotIDs[newSlot.ID] = true
					break
				}
			}
		} else {
			// Create new slot
			newSlot.ScheduleID = schedule.ID
			// Ensure ID is 0 for new slots
			if newSlot.ID <= 0 {
				newSlot.ID = 0
			}
			schedule.TimeSlots = append(schedule.TimeSlots, newSlot)
		}
	}

	// Remove slots that weren't in the request
	var updatedTimeSlots []models.TimeSlot
	for _, slot := range schedule.TimeSlots {
		if processedSlotIDs[slot.ID] || slot.ID == 0 {
			// Keep both processed existing slots and new slots (ID=0)
			updatedTimeSlots = append(updatedTimeSlots, slot)
		}
	}
	schedule.TimeSlots = updatedTimeSlots

	if err := h.scheduleRepo.Update(schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule"})
		return
	}

	// Log the updated schedule
	println("Updated schedule with time slots:")
	for _, slot := range schedule.TimeSlots {
		println("Time slot triggerTime:", slot.TriggerTime)
	}

	// Update scheduler with updated schedule
	schedules, err := h.scheduleRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update scheduler"})
		return
	}
	h.scheduler.UpdateSchedules(schedules)

	c.JSON(http.StatusOK, schedule)
}

// Delete deletes a schedule
func (h *ScheduleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	if err := h.scheduleRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete schedule"})
		return
	}

	// Update scheduler with remaining schedules
	schedules, err := h.scheduleRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update scheduler"})
		return
	}
	h.scheduler.UpdateSchedules(schedules)

	c.JSON(http.StatusOK, gin.H{"message": "Schedule deleted successfully"})
}

// TriggerNow manually triggers the bell
func (h *ScheduleHandler) TriggerNow(c *gin.Context) {
	userID := c.GetInt64("user_id")
	username := c.GetString("username")

	if err := h.scheduler.TriggerNow(userID, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to trigger bell"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bell triggered successfully"})
}
