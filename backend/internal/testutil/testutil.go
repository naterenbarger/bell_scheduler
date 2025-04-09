package testutil

import (
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// NewTestDB creates a new test database connection
func NewTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, mock, cleanup
}

// LoadTestEnv loads test environment variables
func LoadTestEnv(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	os.Setenv("DB_PATH", ":memory:")
	os.Setenv("SMTP_HOST", "localhost")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test")
	os.Setenv("SMTP_PASSWORD", "test")
	os.Setenv("SMTP_FROM", "test@example.com")
}

// AssertError checks if an error matches the expected error
func AssertError(t *testing.T, got, want error) {
	if want == nil {
		assert.NoError(t, got)
		return
	}
	assert.Error(t, got)
	assert.Equal(t, want.Error(), got.Error())
} 