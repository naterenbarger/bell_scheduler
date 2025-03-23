# Bell Scheduler Application

A web-based application for managing bell schedules with a Vue.js frontend and Go backend.

## Prerequisites

- Go 1.17 or later
- Node.js 14.x or later
- npm or yarn

## Project Structure

```
.
├── backend/           # Go backend application
│   ├── cmd/          # Main application entry point
│   ├── internal/     # Internal packages
│   │   ├── models/   # Data models
│   │   ├── handlers/ # HTTP handlers
│   │   └── middleware/ # Middleware components
│   └── pkg/          # Public packages
└── frontend/         # Vue.js frontend application
```

## Setup Instructions

### Backend Setup

1. Install Go from https://golang.org/dl/
2. Navigate to the backend directory:
   ```bash
   cd backend
   ```
3. Initialize the Go module and install dependencies:
   ```bash
   go mod init bell_scheduler
   go get -u github.com/gin-gonic/gin github.com/golang-jwt/jwt/v5 github.com/joho/godotenv
   ```
4. Create a `.env` file in the backend directory with the following content:
   ```
   PORT=8080
   JWT_SECRET=your-secret-key
   DB_CONNECTION=sqlite3.db
   ```
5. Run the backend:
   ```bash
   go run cmd/main.go
   ```

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Run the development server:
   ```bash
   npm run serve
   ```

## Development

- Backend API runs on http://localhost:8080
- Frontend development server runs on http://localhost:8081
- API documentation will be available at http://localhost:8080/api/docs

## Features

- User authentication with JWT
- Schedule management (CRUD operations)
- Time slot management for schedules
- Global bell ring duration configuration
- Default schedule setting
- Responsive design with Vuetify

## Security

- JWT-based authentication
- Password hashing
- Role-based access control
- HTTPS support (in production)

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 