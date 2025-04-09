# Bell Scheduler Backend

This is the backend service for the Bell Scheduler application, built with Go and the Gin framework.

## Project Structure

```
backend/
├── cmd/            # Main applications
├── internal/       # Private application code
│   ├── config/     # Configuration
│   ├── handlers/   # HTTP handlers
│   ├── middleware/ # HTTP middleware
│   ├── models/     # Data models
│   ├── services/   # Business logic
│   └── store/      # Data access layer
├── pkg/           # Public packages
└── .env.example   # Example environment variables
```

## Prerequisites

- Go 1.17 or later
- SQLite3

## Setup

1. Copy `.env.example` to `.env` and configure your environment variables:
   ```bash
   cp .env.example .env
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   go run cmd/main.go
   ```

## Development

The application uses the following patterns:
- Repository pattern for data access
- Service layer for business logic
- Handler layer for HTTP endpoints
- Middleware for cross-cutting concerns

## API Endpoints

### Authentication
- POST `/api/auth/login` - User login
- POST `/api/auth/register` - User registration
- POST `/api/auth/forgot-password` - Request password reset
- POST `/api/auth/reset-password` - Reset password
- POST `/api/auth/change-password` - Change password

### Schedules
- GET `/api/schedules` - List all schedules
- POST `/api/schedules` - Create a schedule
- GET `/api/schedules/:id` - Get schedule details
- PUT `/api/schedules/:id` - Update schedule
- DELETE `/api/schedules/:id` - Delete schedule
- POST `/api/schedules/:id/set-default` - Set default schedule

### Schedule Times
- GET `/api/schedules/:scheduleId/times` - List schedule times
- POST `/api/schedules/:scheduleId/times` - Add time slot
- PUT `/api/schedules/:scheduleId/times/:id` - Update time slot
- DELETE `/api/schedules/:scheduleId/times/:id` - Delete time slot

### Settings
- GET `/api/settings` - Get global settings
- PUT `/api/settings` - Update global settings

### Admin
- GET `/api/admin/users` - List all users
- POST `/api/admin/users` - Create user
- PUT `/api/admin/users/:id` - Update user
- DELETE `/api/admin/users/:id` - Delete user 