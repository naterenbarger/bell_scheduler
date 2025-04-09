package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"bell_scheduler/internal/config"
	"bell_scheduler/internal/handlers"
	"bell_scheduler/internal/middleware"
	"bell_scheduler/internal/models"
	"bell_scheduler/internal/services"
	"bell_scheduler/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := store.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	userRepo := store.NewUserRepository(db)
	scheduleRepo := store.NewScheduleRepository(db)
	settingsRepo := store.NewSettingsRepository(db)
	logRepo := store.NewLogRepository(db)

	// Load settings
	settings, err := settingsRepo.Get()
	if err != nil {
		log.Printf("Warning: Failed to load settings, using defaults: %v", err)
		settings = models.DefaultSettings()
	}

	// Initialize GPIO service
	gpioService, err := services.NewGPIOService(settings.GPIOPin, settings.RingDuration)
	if err != nil {
		log.Fatalf("Failed to initialize GPIO: %v", err)
	}
	defer gpioService.Close()

	// Initialize scheduler service
	scheduler := services.NewSchedulerService(gpioService, logRepo, scheduleRepo)
	scheduler.Start()
	defer scheduler.Stop()

	// Load active schedules
	schedules, err := scheduleRepo.GetAll()
	if err != nil {
		log.Printf("Warning: Failed to load schedules: %v", err)
	}
	scheduler.UpdateSchedules(schedules)

	// Initialize email service
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	emailService := services.NewEmailService(
		os.Getenv("SMTP_HOST"),
		smtpPort,
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_FROM"),
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, emailService, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userRepo)
	scheduleHandler := handlers.NewScheduleHandler(scheduleRepo, scheduler)
	settingsHandler := handlers.NewSettingsHandler(settingsRepo, scheduler)
	logHandler := handlers.NewLogHandler(logRepo)

	// Setup router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Public routes
	router.POST("/api/auth/login", authHandler.Login)
	router.POST("/api/auth/register", authHandler.Register)
	router.POST("/api/auth/forgot-password", authHandler.ForgotPassword)
	router.POST("/api/auth/reset-password", authHandler.ResetPassword)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.Auth(cfg.JWTSecret))
	{
		// User routes
		protected.GET("/users", userHandler.GetUsers)
		protected.POST("/users", userHandler.CreateUser)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
		protected.POST("/auth/change-password", authHandler.ChangePassword)

		// Schedule routes
		protected.GET("/schedules", scheduleHandler.GetAll)
		protected.POST("/schedules", scheduleHandler.Create)
		protected.GET("/schedules/:id", scheduleHandler.Get)
		protected.PUT("/schedules/:id", scheduleHandler.Update)
		protected.DELETE("/schedules/:id", scheduleHandler.Delete)
		protected.POST("/schedules/:id/trigger", scheduleHandler.TriggerNow)
		protected.PUT("/schedules/:id/default", scheduleHandler.SetDefault)
		protected.PUT("/schedules/:id/temporary", scheduleHandler.SetTemporary)
		protected.PUT("/schedules/:id/active", scheduleHandler.SetActive)

		// Settings routes
		protected.GET("/settings", settingsHandler.Get)
		protected.PUT("/settings", settingsHandler.Update)

		// Log routes
		protected.GET("/logs", logHandler.GetAll)
		protected.GET("/logs/range", logHandler.GetByDateRange)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		scheduler.Stop()
		gpioService.Close()
		os.Exit(0)
	}()

	// Start server
	if err := router.Run(cfg.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
