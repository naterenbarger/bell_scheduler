package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetDefault sets a schedule as the default and active schedule
func (h *ScheduleHandler) SetDefault(c *gin.Context) {
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

	// Set as default
	if err := h.scheduleRepo.SetDefault(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set schedule as default"})
		return
	}

	// Update scheduler with updated schedules
	schedules, err := h.scheduleRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update scheduler"})
		return
	}
	h.scheduler.UpdateSchedules(schedules)

	c.JSON(http.StatusOK, gin.H{"message": "Schedule set as default successfully", "schedule": schedule})
}
