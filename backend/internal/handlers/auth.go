package handlers

import (
	"bell_scheduler/internal/models"
	"bell_scheduler/internal/services"
	"bell_scheduler/internal/store"
	"bell_scheduler/pkg/ratelimiter"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userRepo     *store.GormUserRepository
	emailService *services.EmailService
	jwtSecret    []byte
	resetLimiter *ratelimiter.RateLimiter
}

func NewAuthHandler(userRepo *store.GormUserRepository, emailService *services.EmailService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:     userRepo,
		emailService: emailService,
		jwtSecret:    []byte(jwtSecret),
		resetLimiter: ratelimiter.NewRateLimiter(5*time.Minute, 3), // 3 attempts per 5 minutes
	}
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Debug: Log the user details (excluding password)
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("Login attempt for user: %+v\n", struct {
		ID       int64
		Username string
		Email    string
		Role     string
		IsActive bool
	}{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	})))

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is inactive"})
		return
	}

	// Debug: Log password check result
	passwordMatch := user.CheckPassword(req.Password)
	gin.DefaultWriter.Write([]byte(fmt.Sprintf("Password check result: %v\n", passwordMatch)))

	if !passwordMatch {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":               user.ID,
		"username":              user.Username,
		"role":                  user.Role,
		"force_password_change": user.ForcePasswordChange,
		"exp":                   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: tokenString,
		User:  *user,
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	existingUser, _ := h.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	// Note: Password will be hashed in the repository's Create method

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// ForgotPassword handles password reset requests
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check rate limit
	allowed, waitTime := h.resetLimiter.IsAllowed(req.Email)
	if !allowed {
		c.Header("Retry-After", waitTime.String())
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many password reset attempts. Please try again later.",
		})
		return
	}

	// Find user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		// Don't reveal if email exists or not
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link has been sent"})
		return
	}

	// Generate reset token
	if err := user.GenerateResetToken(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}

	// Save user with reset token
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset token"})
		return
	}

	// Send reset email
	if err := h.emailService.SendPasswordResetEmail(user.Email, user.ResetToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link has been sent"})
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by reset token
	user, err := h.userRepo.GetByResetToken(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	// Validate token
	if !user.IsResetTokenValid(req.Token) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	// Update password
	user.Password = req.Password
	// Removed HashPassword call as the repository will handle it

	// Clear reset token
	user.ClearResetToken()

	// Clear force password change flag
	user.ForcePasswordChange = false

	// Save user
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Clear rate limit for this email after successful password reset
	h.resetLimiter.Clear(user.Email)

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}

// ChangePassword handles password change requests
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context (set by auth middleware)
	userID := c.GetInt64("user_id")
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Verify current password
	if !user.CheckPassword(req.CurrentPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Update password
	user.Password = req.NewPassword
	// Removed HashPassword call as the repository will handle it

	// Clear force password change flag
	user.ForcePasswordChange = false

	// Save user
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ValidateToken validates the JWT token
func (h *AuthHandler) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return h.jwtSecret, nil
	})
}
