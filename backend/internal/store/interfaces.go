package store

import "bell_scheduler/internal/models"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetByID(id int64) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByResetToken(token string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int64) error
	GetAll() ([]models.User, error)
	GetAllWithPagination(page, limit int, sortBy string, sortDesc bool, search string) ([]models.User, int64, error)
}
