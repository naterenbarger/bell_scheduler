package main

import (
	"fmt"
	"log"

	"bell_scheduler/internal/config"
	"bell_scheduler/internal/models"
	"bell_scheduler/internal/store"

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

	// Initialize repository
	scheduleRepo := store.NewScheduleRepository(db)

	// Create a test schedule with time slots
	testSchedule := &models.Schedule{
		Name:        "Test Schedule",
		Description: "Test schedule for timeslot handling",
		IsDefault:   false,
		TimeSlots: []models.TimeSlot{
			{
				TriggerTime: "09:00",
				Days:        "[\"Monday\",\"Wednesday\",\"Friday\"]",
				Description: "Morning Bell",
			},
			{
				TriggerTime: "12:00",
				Days:        "[\"Monday\",\"Tuesday\",\"Wednesday\",\"Thursday\",\"Friday\"]",
				Description: "Lunch Bell",
			},
		},
	}

	fmt.Println("Creating test schedule...")
	if err := scheduleRepo.Create(testSchedule); err != nil {
		log.Fatalf("Failed to create test schedule: %v", err)
	}
	fmt.Printf("Created schedule with ID: %d\n", testSchedule.ID)
	fmt.Println("Initial time slots:")
	for i, slot := range testSchedule.TimeSlots {
		fmt.Printf("  Slot %d: ID=%d, TriggerTime=%s, Description=%s\n",
			i+1, slot.ID, slot.TriggerTime, slot.Description)
	}

	// Retrieve the schedule to verify it was saved correctly
	fmt.Println("\nRetrieving schedule...")
	retrievedSchedule, err := scheduleRepo.Get(testSchedule.ID)
	if err != nil {
		log.Fatalf("Failed to retrieve schedule: %v", err)
	}
	fmt.Printf("Retrieved schedule with ID: %d\n", retrievedSchedule.ID)
	fmt.Println("Retrieved time slots:")
	for i, slot := range retrievedSchedule.TimeSlots {
		fmt.Printf("  Slot %d: ID=%d, TriggerTime=%s, Description=%s\n",
			i+1, slot.ID, slot.TriggerTime, slot.Description)
	}

	// Update the schedule by modifying one slot, adding a new one, and removing one
	fmt.Println("\nUpdating schedule...")
	// Modify the first slot
	retrievedSchedule.TimeSlots[0].TriggerTime = "08:30"
	retrievedSchedule.TimeSlots[0].Description = "Updated Morning Bell"

	// Add a new slot
	retrievedSchedule.TimeSlots = append(retrievedSchedule.TimeSlots, models.TimeSlot{
		TriggerTime: "15:00",
		Days:        "[\"Monday\",\"Tuesday\",\"Wednesday\",\"Thursday\",\"Friday\"]",
		Description: "Dismissal Bell",
	})

	// Update the schedule
	if err := scheduleRepo.Update(retrievedSchedule); err != nil {
		log.Fatalf("Failed to update schedule: %v", err)
	}

	// Retrieve the updated schedule
	fmt.Println("\nRetrieving updated schedule...")
	updatedSchedule, err := scheduleRepo.Get(testSchedule.ID)
	if err != nil {
		log.Fatalf("Failed to retrieve updated schedule: %v", err)
	}
	fmt.Printf("Retrieved updated schedule with ID: %d\n", updatedSchedule.ID)
	fmt.Println("Updated time slots:")
	for i, slot := range updatedSchedule.TimeSlots {
		fmt.Printf("  Slot %d: ID=%d, TriggerTime=%s, Description=%s\n",
			i+1, slot.ID, slot.TriggerTime, slot.Description)
	}

	// Clean up - delete the test schedule
	fmt.Println("\nCleaning up - deleting test schedule...")
	if err := scheduleRepo.Delete(testSchedule.ID); err != nil {
		log.Fatalf("Failed to delete test schedule: %v", err)
	}
	fmt.Println("Test schedule deleted successfully")

	fmt.Println("\nTest completed successfully!")
}
