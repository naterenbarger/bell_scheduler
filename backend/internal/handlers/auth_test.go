package handlers

import (
	"bell_scheduler/internal/models"
	"bell_scheduler/internal/services"
	"bell_scheduler/internal/store"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

var _ store.UserRepository = (*MockUserRepository)(nil) // Type assertion to ensure interface implementation

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id int64) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByResetToken(token string) (*models.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &MockUserRepository{}
	mockEmailService := services.NewEmailService("localhost", 25, "test", "test", "test@test.com")

	// Test successful login
	t.Run("Successful login", func(t *testing.T) {
		user := &models.User{
			BaseModel: models.BaseModel{
				ID: 1,
			},
			Username: "testuser",
			Password: "hashedpassword",
		}

		mockRepo.On("GetByUsername", "testuser").Return(user, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		loginData := map[string]string{
			"username": "testuser",
			"password": "password",
		}
		body, _ := json.Marshal(loginData)
		c.Request = httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))

		handler := NewAuthHandler(mockRepo, mockEmailService, "test_secret")
		handler.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test invalid credentials
	t.Run("Invalid credentials", func(t *testing.T) {
		user := &models.User{
			BaseModel: models.BaseModel{
				ID: 1,
			},
			Username: "testuser",
			Password: "wrongpassword",
		}

		mockRepo.On("GetByUsername", "testuser").Return(user, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		loginData := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(loginData)
		c.Request = httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))

		handler := NewAuthHandler(mockRepo, mockEmailService, "test_secret")
		handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &MockUserRepository{}
	mockEmailService := services.NewEmailService("localhost", 25, "test", "test", "test@test.com")

	// Test successful registration
	t.Run("Successful registration", func(t *testing.T) {
		mockRepo.On("GetByUsername", "newuser").Return(nil, nil)
		mockRepo.On("GetByEmail", "test@test.com").Return(nil, nil)
		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		registerData := map[string]string{
			"username": "newuser",
			"password": "password",
			"email":    "test@test.com",
		}
		body, _ := json.Marshal(registerData)
		c.Request = httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))

		handler := NewAuthHandler(mockRepo, mockEmailService, "test_secret")
		handler.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestValidateToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &MockUserRepository{}
	mockEmailService := services.NewEmailService("localhost", 25, "test", "test", "test@test.com")

	// Test invalid token
	t.Run("Invalid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid_token")

		handler := NewAuthHandler(mockRepo, mockEmailService, "test_secret")
		handler.ValidateToken("invalid_token")

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
