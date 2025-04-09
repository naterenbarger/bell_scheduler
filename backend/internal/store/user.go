package store

import (
	"bell_scheduler/internal/models"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct {
	db *gorm.DB
	// Removed cache field to prevent stale data issues
}

func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{
		db: db,
	}
}

// GetByID retrieves a user by ID
func (r *GormUserRepository) GetByID(id int64) (*models.User, error) {
	// For authentication purposes, always fetch from database to ensure fresh password hash
	// This prevents issues with stale password hashes in the cache
	fmt.Printf("[GetByID] Looking up user with ID: %d\n", id)

	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		fmt.Printf("[GetByID] Error fetching user from database: %v\n", err)
		return nil, err
	}

	fmt.Printf("[GetByID] User fetched from database, password hash length: %d\n", len(user.Password))
	return &user, nil
}

// Create creates a new user
func (r *GormUserRepository) Create(user *models.User) error {
	// Hash the password before saving
	if err := user.HashPassword(); err != nil {
		return err
	}

	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	// No cache to invalidate
	return nil
}

// Update updates an existing user
func (r *GormUserRepository) Update(user *models.User) error {
	// Get the current user from the database
	var currentUser models.User
	if err := r.db.First(&currentUser, user.ID).Error; err != nil {
		fmt.Printf("[Update] Error fetching current user (ID: %d): %v\n", user.ID, err)
		return err
	}

	fmt.Printf("[Update] ===== PASSWORD UPDATE DEBUGGING =====\n")
	fmt.Printf("[Update] User ID: %d, Username: %s\n", user.ID, user.Username)
	fmt.Printf("[Update] Incoming password - Is empty: %v, Length: %d\n", user.Password == "", len(user.Password))
	fmt.Printf("[Update] Current DB password hash length: %d\n", len(currentUser.Password))

	// Check if the incoming password looks like a bcrypt hash
	isLikelyHash := false
	if len(user.Password) > 0 {
		isLikelyHash = len(user.Password) > 50 && strings.HasPrefix(user.Password, "$2a$")
		fmt.Printf("[Update] Is incoming password likely already a hash? %v\n", isLikelyHash)

		// For debugging, try to compare with existing hash
		if isLikelyHash {
			fmt.Printf("[Update] WARNING: Incoming password appears to be a hash already!\n")
			fmt.Printf("[Update] This suggests double hashing may be occurring\n")
		}
	}

	// Only hash the password if it's provided (not empty) and doesn't look like a hash already
	if user.Password != "" && !isLikelyHash {
		fmt.Printf("[Update] Password provided and appears to be plaintext, about to hash password\n")
		plaintextPassword := user.Password // Store for logging

		// Log the plaintext password length for debugging
		fmt.Printf("[Update] PLAINTEXT PASSWORD LENGTH: %d\n", len(plaintextPassword))

		if err := user.HashPassword(); err != nil {
			fmt.Printf("[Update] Error hashing password: %v\n", err)
			return err
		}
		fmt.Printf("[Update] Password hashed successfully - Original length: %d, Hashed length: %d\n",
			len(plaintextPassword), len(user.Password))
		fmt.Printf("[Update] New password hash (first 10 chars): %s...\n", user.Password[:10])
	} else if user.Password != "" && isLikelyHash {
		fmt.Printf("[Update] Password appears to already be hashed, skipping hash operation\n")
		fmt.Printf("[Update] Using provided hash (first 10 chars): %s...\n", user.Password[:10])
	} else {
		// If password is empty, keep the current hashed password
		fmt.Printf("[Update] No password provided, keeping current password hash\n")
		user.Password = currentUser.Password
		fmt.Printf("[Update] Using existing hash (first 10 chars): %s...\n", user.Password[:10])
	}

	fmt.Printf("[Update] Final password state before save - Length: %d\n", len(user.Password))
	fmt.Printf("[Update] ===== END PASSWORD UPDATE DEBUGGING =====\n")

	if err := r.db.Save(user).Error; err != nil {
		fmt.Printf("[Update] Error saving user: %v\n", err)
		return err
	}
	fmt.Printf("[Update] User saved successfully to database\n")

	// Verify the user was updated correctly by fetching it again
	var updatedUser models.User
	if err := r.db.First(&updatedUser, user.ID).Error; err != nil {
		fmt.Printf("[Update] Error fetching updated user: %v\n", err)
	} else {
		fmt.Printf("[Update] Verification - Updated password hash length: %d\n", len(updatedUser.Password))
		fmt.Printf("[Update] Verification - Updated password hash (first 10 chars): %s...\n", updatedUser.Password[:10])
	}

	return nil
}

// GetByUsername retrieves a user by username
func (r *GormUserRepository) GetByUsername(username string) (*models.User, error) {
	// For authentication purposes, always fetch from database to ensure fresh password hash
	// This prevents issues with stale password hashes in the cache
	fmt.Printf("[GetByUsername] Looking up user with username: %s\n", username)

	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Printf("[GetByUsername] Error fetching user from database: %v\n", err)
		return nil, err
	}

	fmt.Printf("[GetByUsername] User fetched from database, password hash length: %d\n", len(user.Password))
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *GormUserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByResetToken retrieves a user by reset token
func (r *GormUserRepository) GetByResetToken(token string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("reset_token = ? AND reset_token_expiry > ?", token, time.Now()).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all users
func (r *GormUserRepository) GetAll() ([]models.User, error) {
	// Always fetch from database to ensure fresh data
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// Delete deletes a user
func (r *GormUserRepository) Delete(id int64) error {
	if err := r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}

	// No cache to invalidate
	fmt.Printf("[Delete] User deleted directly from database\n")
	return nil
}

// GetAllWithPagination retrieves users with pagination, sorting, and filtering
func (r *GormUserRepository) GetAllWithPagination(page, limit int, sortBy string, sortDesc bool, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Build query
	query := r.db.Model(&models.User{})

	// Apply search filter if provided
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sortDesc {
		query = query.Order(sortBy + " DESC")
	} else {
		query = query.Order(sortBy + " ASC")
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// No caching to ensure fresh data

	return users, total, nil
}
