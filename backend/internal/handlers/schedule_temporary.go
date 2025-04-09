package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetTemporary sets a schedule as temporary and active
func (h *ScheduleHandler) SetTemporary(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	// Check if schedule exists
	schedule, err := h.scheduleRepo.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	// Parse request body to get temporary flag
	var request struct {
		IsTemporary bool `json:"isTemporary"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		// If no body provided, default to not temporary
		request.IsTemporary = false
	}

	// Update the schedule to set temporary flag
	schedule.IsTemporary = request.IsTemporary
	schedule.IsActive = true // Also set as active

	// Save the temporary flag
	if err := h.scheduleRepo.Update(schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule"})
		return
	}

	// Set as active
	if err := h.scheduleRepo.SetActive(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set schedule as active"})
		return
	}

	// Update scheduler with updated schedules
	schedules, err := h.scheduleRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update scheduler"})
		return
	}
	h.scheduler.UpdateSchedules(schedules)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Schedule set as active successfully",
		"schedule": schedule,
	})
}
