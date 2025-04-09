package store

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"bell_scheduler/internal/models"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return db, mock, func() {
		sqlDB.Close()
	}
}

func TestUserRepository_Create(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		user    *models.User
		mock    func()
		wantErr bool
	}{
		{
			name: "success",
			user: &models.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "duplicate username",
			user: &models.User{
				Username: "existing",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("existing", "test@example.com", sqlmock.AnyArg(), "user").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Create(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	tests := []struct {
		name     string
		username string
		mock     func()
		want     *models.User
		wantErr  bool
	}{
		{
			name:     "success",
			username: "testuser",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "role", "created_at", "updated_at"}).
					AddRow(1, "testuser", "test@example.com", "hashedpassword", "user", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			want: &models.User{
				BaseModel: models.BaseModel{
					ID: 1,
				},
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpassword",
				Role:     "user",
			},
			wantErr: false,
		},
		{
			name:     "not found",
			username: "nonexistent",
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetByUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.Username, got.Username)
			assert.Equal(t, tt.want.Email, got.Email)
			assert.Equal(t, tt.want.Role, got.Role)
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		user    *models.User
		mock    func()
		wantErr bool
	}{
		{
			name: "success",
			user: &models.User{
				BaseModel: models.BaseModel{
					ID: 1,
				},
				Username: "updated",
				Email:    "updated@example.com",
				Role:     "admin",
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `users`").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "not found",
			user: &models.User{
				BaseModel: models.BaseModel{
					ID: 999,
				},
				Username: "nonexistent",
			},
			mock: func() {
				mock.ExpectExec("UPDATE users").
					WithArgs("nonexistent", "", "", 999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Update(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		id      int64
		mock    func()
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `users`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "not found",
			id:   999,
			mock: func() {
				mock.ExpectExec("DELETE FROM users").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Delete(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestUserRepository_GetAll(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		mock    func()
		want    []models.User
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "role", "created_at", "updated_at"}).
					AddRow(1, "user1", "user1@example.com", "hashedpassword", "user", time.Now(), time.Now()).
					AddRow(2, "user2", "user2@example.com", "hashedpassword", "admin", time.Now(), time.Now())
				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WillReturnRows(rows)
			},
			want: []models.User{
				{
					BaseModel: models.BaseModel{
						ID: 1,
					},
					Username: "user1",
					Email:    "user1@example.com",
					Password: "hashedpassword",
					Role:     "user",
				},
				{
					BaseModel: models.BaseModel{
						ID: 2,
					},
					Username: "user2",
					Email:    "user2@example.com",
					Password: "hashedpassword",
					Role:     "admin",
				},
			},
			wantErr: false,
		},
		{
			name: "empty",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "role", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rows)
			},
			want:    []models.User{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetAll()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got))
			for i := range tt.want {
				assert.Equal(t, tt.want[i].ID, got[i].ID)
				assert.Equal(t, tt.want[i].Username, got[i].Username)
				assert.Equal(t, tt.want[i].Email, got[i].Email)
				assert.Equal(t, tt.want[i].Role, got[i].Role)
			}
		})
	}
}
