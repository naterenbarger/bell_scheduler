package handlers

import (
	"net/http"
	"time"

	"bell_scheduler/internal/models"
	"bell_scheduler/internal/services"
	"bell_scheduler/internal/store"

	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	settingsRepo *store.SettingsRepository
	scheduler    *services.SchedulerService
}

func NewSettingsHandler(settingsRepo *store.SettingsRepository, scheduler *services.SchedulerService) *SettingsHandler {
	return &SettingsHandler{
		settingsRepo: settingsRepo,
		scheduler:    scheduler,
	}
}

// Get returns the current settings
func (h *SettingsHandler) Get(c *gin.Context) {
	settings, err := h.settingsRepo.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get settings"})
		return
	}

	// Convert duration to seconds for frontend
	response := struct {
		RingDuration int    `json:"ringDuration"`
		GPIOPin      int    `json:"gpioPin"`
		Timezone     string `json:"timezone"`
	}{
		RingDuration: int(settings.RingDuration.Seconds()),
		GPIOPin:      settings.GPIOPin,
		Timezone:     settings.Timezone,
	}

	c.JSON(http.StatusOK, response)
}

// Update updates the settings
func (h *SettingsHandler) Update(c *gin.Context) {
	var req models.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings, err := h.settingsRepo.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get settings"})
		return
	}

	// Update settings
	settings.RingDuration = time.Duration(req.RingDuration) * time.Second
	settings.GPIOPin = req.GPIOPin
	settings.Timezone = req.Timezone

	if err := h.settingsRepo.Update(settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	// Update scheduler with new settings
	h.scheduler.SetDuration(settings.RingDuration)

	// Convert duration to seconds for frontend
	response := struct {
		RingDuration int    `json:"ringDuration"`
		GPIOPin      int    `json:"gpioPin"`
		Timezone     string `json:"timezone"`
	}{
		RingDuration: int(settings.RingDuration.Seconds()),
		GPIOPin:      settings.GPIOPin,
		Timezone:     settings.Timezone,
	}

	c.JSON(http.StatusOK, response)
}
