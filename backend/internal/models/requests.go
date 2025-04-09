package models

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=user admin"`
}

// ForgotPasswordRequest represents a password reset request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

// CreateScheduleRequest represents a schedule creation request
type CreateScheduleRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	IsDefault   bool       `json:"isDefault"`
	IsTemporary bool       `json:"isTemporary"`
	TimeSlots   []TimeSlot `json:"timeSlots" binding:"required,dive"`
}

// UpdateScheduleRequest represents a schedule update request
type UpdateScheduleRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	IsDefault   bool       `json:"isDefault"`
	IsTemporary bool       `json:"isTemporary"`
	TimeSlots   []TimeSlot `json:"timeSlots" binding:"required,dive"`
}

// UpdateSettingsRequest represents a settings update request
type UpdateSettingsRequest struct {
	RingDuration int    `json:"ringDuration" binding:"required,min=1,max=60"`
	GPIOPin      int    `json:"gpioPin" binding:"required,min=1,max=40"`
	Timezone     string `json:"timezone" binding:"required"`
}
