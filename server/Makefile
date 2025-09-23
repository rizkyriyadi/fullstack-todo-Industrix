# Todo Backend Makefile

# Variables
APP_NAME=todo-backend
BINARY_NAME=todo-api
MAIN_PATH=./cmd/api

# Go commands
.PHONY: build run clean test deps help

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Starting development server with hot reload..."
	@air

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run

# Database commands
db-migrate-up:
	@echo "Running database migrations..."
	# This would run migrations using migrate tool or custom migration runner

db-migrate-down:
	@echo "Rolling back database migrations..."
	# This would rollback migrations

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

docker-compose-up:
	@echo "Starting services with docker-compose..."
	@docker-compose up -d

docker-compose-down:
	@echo "Stopping services with docker-compose..."
	@docker-compose down

docker-compose-logs:
	@echo "Viewing docker-compose logs..."
	@docker-compose logs -f

# Help
help:
	@echo "Available commands:"
	@echo "  build              Build the application"
	@echo "  run                Run the application"
	@echo "  dev                Run with hot reload (requires air)"
	@echo "  clean              Clean build artifacts"
	@echo "  deps               Download and tidy dependencies"
	@echo "  test               Run tests"
	@echo "  test-coverage      Run tests with coverage report"
	@echo "  fmt                Format code"
	@echo "  lint               Lint code (requires golangci-lint)"
	@echo "  docker-build       Build Docker image"
	@echo "  docker-run         Run Docker container"
	@echo "  docker-compose-up  Start all services with docker-compose"
	@echo "  docker-compose-down Stop all services with docker-compose"
	@echo "  docker-compose-logs View docker-compose logs"
	@echo "  help               Show this help message"