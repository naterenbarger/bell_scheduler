package services

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// GPIOService handles GPIO operations for the bell system
type GPIOService struct {
	pin      rpio.Pin
	duration time.Duration
	isActive bool
	mock     bool
}

// NewGPIOService creates a new GPIO service instance
func NewGPIOService(pinNumber int, duration time.Duration) (*GPIOService, error) {
	service := &GPIOService{
		duration: duration,
		isActive: false,
	}

	// Try to open GPIO, if it fails, run in mock mode
	if err := rpio.Open(); err != nil {
		fmt.Println("Running in mock GPIO mode")
		service.mock = true
		return service, nil
	}

	service.mock = false
	service.pin = rpio.Pin(pinNumber)
	service.pin.Output()
	service.pin.Low() // Ensure pin starts in low state

	return service, nil
}

// Trigger activates the relay for the configured duration
func (s *GPIOService) Trigger() error {
	if s.isActive {
		return fmt.Errorf("relay is already active")
	}

	s.isActive = true
	if !s.mock {
		s.pin.High()
	} else {
		fmt.Println("\n🔔 BELL TRIGGERED 🔔")
		fmt.Println("┌─────────────────┐")
		fmt.Println("│   BELL RINGING  │")
		fmt.Println("│     ╔═══╗      │")
		fmt.Println("│     ║   ║      │")
		fmt.Println("│     ║   ║      │")
		fmt.Println("│     ║   ║      │")
		fmt.Println("│     ╚═══╝      │")
		fmt.Println("└─────────────────┘")
		fmt.Printf("Duration: %v\n", s.duration)
	}

	// Start a goroutine to handle the duration
	go func() {
		time.Sleep(s.duration)
		if !s.mock {
			s.pin.Low()
		} else {
			fmt.Println("\n🔕 BELL STOPPED 🔕")
			fmt.Println("┌─────────────────┐")
			fmt.Println("│   BELL STOPPED  │")
			fmt.Println("│     ╔═══╗      │")
			fmt.Println("│     ║   ║      │")
			fmt.Println("│     ║   ║      │")
			fmt.Println("│     ║   ║      │")
			fmt.Println("│     ╚═══╝      │")
			fmt.Println("└─────────────────┘")
		}
		s.isActive = false
	}()

	return nil
}

// SetDuration updates the trigger duration
func (s *GPIOService) SetDuration(duration time.Duration) {
	s.duration = duration
}

// IsActive returns whether the relay is currently active
func (s *GPIOService) IsActive() bool {
	return s.isActive
}

// Close cleans up GPIO resources
func (s *GPIOService) Close() {
	if !s.mock {
		s.pin.Low()
		rpio.Close()
	}
}
