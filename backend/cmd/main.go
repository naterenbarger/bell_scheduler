package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Set up Gin router
	r := gin.Default()

	// Serve static files from the frontend build directory
	r.Static("/", "./frontend/dist")

	// API routes group
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", handleLogin)
			auth.POST("/logout", handleLogout)
		}

		// Schedule routes
		schedules := api.Group("/schedules")
		{
			schedules.GET("/", handleGetSchedules)
			schedules.POST("/", handleCreateSchedule)
			schedules.GET("/:id", handleGetSchedule)
			schedules.PUT("/:id", handleUpdateSchedule)
			schedules.DELETE("/:id", handleDeleteSchedule)
		}

		// Schedule times routes
		times := api.Group("/schedules/:scheduleId/times")
		{
			times.GET("/", handleGetScheduleTimes)
			times.POST("/", handleCreateScheduleTime)
			times.PUT("/:id", handleUpdateScheduleTime)
			times.DELETE("/:id", handleDeleteScheduleTime)
		}

		// Settings routes
		settings := api.Group("/settings")
		{
			settings.GET("/", handleGetSettings)
			settings.PUT("/", handleUpdateSettings)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 