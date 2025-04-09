package handlers

import (
	"bell_scheduler/internal/models"
	"bell_scheduler/internal/store"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo store.UserRepository
}

func NewUserHandler(userRepo store.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetUsers retrieves all users with pagination, sorting, and filtering
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Parse sorting parameters
	sortBy := c.DefaultQuery("sort_by", "username")
	sortDesc := c.DefaultQuery("sort_desc", "false") == "true"

	// Parse search parameter
	search := c.DefaultQuery("search", "")

	// Get users with pagination
	users, total, err := h.userRepo.GetAllWithPagination(page, limit, sortBy, sortDesc, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Clear passwords before sending response
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
	})
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Create user from request
	user := models.User{
		Username:            req.Username,
		Email:               req.Email,
		Password:            req.Password,
		Role:                req.Role,
		ForcePasswordChange: true, // Set to true by default for new users
	}

	// Validate required fields
	if err := validateUserInput(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username or email already exists
	if existingUser, _ := h.userRepo.GetByUsername(user.Username); existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}
	if existingUser, _ := h.userRepo.GetByEmail(user.Email); existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Note: Password will be hashed in the repository's Create method
	// No need to hash it here

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	// Create user
	if err := h.userRepo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Clear password before sending response
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	fmt.Printf("[UpdateUser] Starting update for user ID: %s\n", id)

	// Log the raw request body for debugging
	var rawBody map[string]interface{}
	bodyData, _ := c.GetRawData()
	fmt.Printf("[UpdateUser] Raw request body: %s\n", string(bodyData))

	// Restore the body for later binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyData))

	// Try to unmarshal the raw body for debugging
	if err := json.Unmarshal(bodyData, &rawBody); err == nil {
		fmt.Printf("[UpdateUser] Parsed request body: %+v\n", rawBody)
		fmt.Printf("[UpdateUser] Password in raw request: %v\n", rawBody["password"])
	}

	// Convert string ID to int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get existing user
	existingUser, err := h.userRepo.GetByID(idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	fmt.Printf("[UpdateUser] Found existing user: %s (ID: %d)\n", existingUser.Username, existingUser.ID)
	fmt.Printf("[UpdateUser] Existing password hash length: %d\n", len(existingUser.Password))

	var updateData models.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	fmt.Printf("[UpdateUser] Received update data for user: %s\n", updateData.Username)
	fmt.Printf("[UpdateUser] Password provided in update: %v (length: %d)\n", updateData.Password != "", len(updateData.Password))

	// Set the ID from existing user to prevent password requirement validation
	updateData.ID = existingUser.ID

	// Validate input
	if err := validateUserInput(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check username uniqueness if changed
	if updateData.Username != existingUser.Username {
		if existingUser, _ := h.userRepo.GetByUsername(updateData.Username); existingUser != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
	}

	// Check email uniqueness if changed
	if updateData.Email != existingUser.Email {
		if existingUser, _ := h.userRepo.GetByEmail(updateData.Email); existingUser != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
	}

	// Update fields
	existingUser.Username = updateData.Username
	existingUser.Email = updateData.Email
	existingUser.Role = updateData.Role
	existingUser.IsActive = updateData.IsActive

	// Update password if provided
	if updateData.Password != "" {
		fmt.Printf("[UpdateUser] Setting plaintext password from update data (length: %d)\n", len(updateData.Password))
		// Store the plaintext password directly
		// The repository will handle the hashing
		existingUser.Password = updateData.Password
	} else {
		fmt.Printf("[UpdateUser] No password provided, setting empty password to prevent re-hashing\n")
		// If no new password provided, set to empty to prevent re-hashing
		// The repository will keep the existing hash
		existingUser.Password = ""
	}

	fmt.Printf("[UpdateUser] Before calling repository Update - Password state: %v (length: %d)\n",
		existingUser.Password != "", len(existingUser.Password))

	// Save updates
	if err := h.userRepo.Update(existingUser); err != nil {
		fmt.Printf("[UpdateUser] Error updating user: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	fmt.Printf("[UpdateUser] User updated successfully\n")

	// Clear password before sending response
	existingUser.Password = ""
	c.JSON(http.StatusOK, existingUser)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Convert string ID to int64
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user exists
	if _, err := h.userRepo.GetByID(idInt); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete user
	if err := h.userRepo.Delete(idInt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

// validateUserInput validates user input data
func validateUserInput(user *models.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if len(user.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(user.Email, "@") {
		return errors.New("invalid email format")
	}
	// For new users (ID == 0), password is required
	if user.ID == 0 {
		if user.Password == "" {
			return errors.New("password is required for new users")
		}
	}
	// When password is provided (for both new and existing users), validate length
	if user.Password != "" && len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if user.Role != "" && user.Role != "admin" && user.Role != "user" {
		return errors.New("invalid role")
	}
	return nil
}
