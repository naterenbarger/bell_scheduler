package handlers

import (
    "bell_scheduler/internal/models"
    "bell_scheduler/internal/store"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

// LogHandler handles HTTP requests for log entries
type LogHandler struct {
    logRepo *store.LogRepository
}

// NewLogHandler creates a new log handler instance
func NewLogHandler(logRepo *store.LogRepository) *LogHandler {
    return &LogHandler{logRepo: logRepo}
}

// GetAll retrieves all log entries
func (h *LogHandler) GetAll(c *gin.Context) {
    logs, err := h.logRepo.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
        return
    }
    c.JSON(http.StatusOK, logs)
}

// GetByDateRange retrieves log entries within a date range
func (h *LogHandler) GetByDateRange(c *gin.Context) {
    startStr := c.Query("start")
    endStr := c.Query("end")

    start, err := time.Parse(time.RFC3339, startStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
        return
    }

    end, err := time.Parse(time.RFC3339, endStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
        return
    }

    logs, err := h.logRepo.GetByDateRange(start, end)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
        return
    }
    c.JSON(http.StatusOK, logs)
}

// CreateLogEntry creates a new log entry
func (h *LogHandler) CreateLogEntry(c *gin.Context) {
    var log models.LogEntry
    if err := c.ShouldBindJSON(&log); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log entry"})
        return
    }

    log.Timestamp = time.Now()
    if err := h.logRepo.Create(&log); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create log entry"})
        return
    }

    c.JSON(http.StatusCreated, log)
} 