package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseModel
	Username            string `json:"username" gorm:"type:varchar(255);uniqueIndex;not null"`
	Email               string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password            string `json:"password,omitempty" gorm:"type:varchar(255);not null"` // Changed from json:"-" to allow unmarshaling
	Role                string `json:"role" gorm:"type:varchar(50);not null;default:'user';index"`
	IsActive            bool   `json:"isActive" gorm:"not null;default:true"`
	ResetToken          string `gorm:"type:varchar(255)"`
	ResetTokenExpiry    *time.Time
	ForcePasswordChange bool `gorm:"type:boolean;default:false"`
}

// HashPassword hashes the user's password using bcrypt
func (u *User) HashPassword() error {
	// Check if the password is already hashed (bcrypt hashes start with $2a$ and are typically longer than 50 chars)
	if len(u.Password) > 50 && strings.HasPrefix(u.Password, "$2a$") {
		fmt.Printf("[HashPassword] Password appears to already be hashed, skipping hash operation\n")
		return nil
	}

	fmt.Printf("[HashPassword] Hashing password (length: %d)\n", len(u.Password))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("[HashPassword] Error hashing password: %v\n", err)
		return err
	}
	u.Password = string(hashedPassword)
	fmt.Printf("[HashPassword] Password hashed successfully (new length: %d)\n", len(u.Password))
	return nil
}

// CheckPassword compares the provided password with the hashed password
func (u *User) CheckPassword(password string) bool {
	fmt.Printf("[CheckPassword] Checking password for user: %s (ID: %d)\n", u.Username, u.ID)
	fmt.Printf("[CheckPassword] Stored hash length: %d\n", len(u.Password))
	fmt.Printf("[CheckPassword] Provided password length: %d\n", len(password))

	// Generate a hash from the provided password for comparison logging
	testHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err == nil {
		fmt.Printf("[CheckPassword] Test hash generated from provided password: %s\n", string(testHash))
	}

	// Actual password check
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		fmt.Printf("[CheckPassword] Password check failed: %v\n", err)
		fmt.Printf("[CheckPassword] Stored hash: %s\n", u.Password)
		return false
	}

	fmt.Printf("[CheckPassword] Password check successful\n")
	return true
}

// GenerateResetToken generates a new reset token
func (u *User) GenerateResetToken() error {
	// Generate random bytes
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}

	// Convert to base64 string
	token := base64.URLEncoding.EncodeToString(b)

	// Set token and expiry (1 hour from now)
	expiry := time.Now().Add(time.Hour)
	u.ResetToken = token
	u.ResetTokenExpiry = &expiry

	return nil
}

// IsResetTokenValid checks if the reset token is valid and not expired
func (u *User) IsResetTokenValid(token string) bool {
	if u.ResetToken == "" || u.ResetTokenExpiry == nil {
		return false
	}
	return u.ResetToken == token && time.Now().Before(*u.ResetTokenExpiry)
}

// ClearResetToken clears the reset token and expiry
func (u *User) ClearResetToken() {
	u.ResetToken = ""
	u.ResetTokenExpiry = nil
}

// CreateUserRequest represents the data needed to create a new user
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"omitempty,oneof=admin user"`
}
