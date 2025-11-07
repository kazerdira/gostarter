.PHONY: help run build test clean migrate-up migrate-down migrate-create sqlc-generate docker-build docker-up docker-down lint fmt

# Default target
help:
	@echo "Available commands:"
	@echo "  make run              - Run the application"
	@echo "  make build            - Build the application"
	@echo "  make test             - Run tests"
	@echo "  make test-coverage    - Run tests with coverage"
	@echo "  make migrate-up       - Apply database migrations"
	@echo "  make migrate-down     - Rollback last migration"
	@echo "  make migrate-create   - Create a new migration (usage: make migrate-create name=create_table_name)"
	@echo "  make sqlc-generate    - Generate SQLC code"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-up        - Start Docker containers"
	@echo "  make docker-down      - Stop Docker containers"
	@echo "  make lint             - Run linters"
	@echo "  make fmt              - Format code"
	@echo "  make clean            - Clean build artifacts"

# Run the application
run:
	@echo "Starting application..."
	@go run cmd/api/main.go

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v -short ./...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -v -run Integration ./...

# Apply database migrations
migrate-up:
	@echo "Applying migrations..."
	@go run cmd/api/main.go migrate up

# Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	@go run cmd/api/main.go migrate down

# Create a new migration
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	@migrate create -ext sql -dir internal/db/migrations -seq $(name)

# Generate SQLC code
sqlc-generate:
	@echo "Generating SQLC code..."
	@sqlc generate

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t go-api:latest .

# Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

# Run golangci-lint
lint:
	@echo "Running linters..."
	@golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Install development dependencies
dev-deps:
	@echo "Installing development dependencies..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
