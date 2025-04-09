.PHONY: test test-coverage build run clean lint format

# Build the application
build:
	go build -o bin/bell_scheduler cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Run linter
lint:
	golangci-lint run

# Format code
format:
	go fmt ./...

# Generate mocks
mock:
	mockgen -source=internal/store/user.go -destination=internal/store/mocks/user_mock.go
	mockgen -source=internal/store/schedule.go -destination=internal/store/mocks/schedule_mock.go
	mockgen -source=internal/store/settings.go -destination=internal/store/mocks/settings_mock.go

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run database migrations
migrate:
	go run cmd/migrate/main.go

# Create a new migration
migration:
	@read -p "Enter migration name: " name; \
	go run cmd/migrate/main.go create $$name

# Run the application with hot reload
dev:
	air

# Build for production
prod:
	GOOS=linux GOARCH=amd64 go build -o bin/bell_scheduler cmd/main.go

# Run security checks
security:
	gosec ./...
	golangci-lint run --enable gosec

# Run all checks (lint, test, build)
check: lint test build 