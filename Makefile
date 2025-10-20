.PHONY: start stop start-backend stop-backend start-backend-no-air check-port clean-air \
	test test-handlers test-handlers-coverage test-unit test-integration \
	test-auth test-profile test-account test-plan test-billing test-legacy \
	test-setup test-clean test-watch \
	lint fmt vet mod-tidy build dev ci-test help

# Clean up any existing Air processes
clean-air:
	@echo "Cleaning up existing Air processes..."
	@pkill -f "air" 2>/dev/null || true
	@pkill -f "tmp/main" 2>/dev/null || true

# Check if port 8080 is in use
check-port:
	@if lsof -i :8080 > /dev/null; then \
		echo "Port 8080 is in use. Stopping existing process..."; \
		lsof -ti :8080 | xargs kill -9 2>/dev/null || true; \
	fi
#@if lsof -i :3000 > /dev/null; then \
#	echo "Port 3000 is in use. Stopping existing process..."; \
#	lsof -ti :3000 | xargs kill -9 2>/dev/null || true; \
#fi

# Start backend server with hot reloading (default)
start-backend: clean-air check-port
	@echo "Starting backend server with Air for hot reloading..."
	@air

# Start backend server without hot reloading
start-backend-no-air: 
	@echo "Starting backend server without hot reloading..."
	@go run main.go

# Stop backend server
stop-backend:
	@echo "Stopping backend server..."
	@pkill -f "air" 2>/dev/null || true
	@pkill -f "tmp/main" 2>/dev/null || true
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || true

db:
	@echo "Starting database..."
	@psql -U postgres -d goiter 

clean:
	@echo "Cleaning database..."
	@psql -U postgres -d postgres -c "DROP DATABASE goiter;"
	@psql -U postgres -d postgres -c "CREATE DATABASE goiter;"


# Test commands
test:
	@echo "Running all standard Go tests..."
	@go test ./... -v

test-handlers:
	@echo "Running handler tests..."
	@go test ./core/handlers -v

test-handlers-coverage:
	@echo "Running handler tests with coverage..."
	@go test ./core/handlers -v -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-unit:
	@echo "Running unit tests..."
	@go test ./core/handlers -short -v

test-integration:
	@echo "Running integration tests..."
	@go test ./core/handlers -run Integration -v

test-auth:
	@echo "Running authentication tests..."
	@go test ./core/handlers -run TestAuthentication -v

test-profile:
	@echo "Running profile handler tests..."
	@go test ./core/handlers -run TestProfileHandler -v

test-account:
	@echo "Running account handler tests..."
	@go test ./core/handlers -run TestAccountHandler -v

test-plan:
	@echo "Running plan handler tests..."
	@go test ./core/handlers -run TestPlanHandler -v

test-billing:
	@echo "Running billing handler tests..."
	@go test ./core/handlers -run TestBillingHandler -v

# Legacy test suite (for migration reference)
test-legacy:
	@echo "Running legacy test suite..."
	@go run testsuite/run/run.go

# Test environment setup
test-setup:
	@echo "Setting up test environment..."
	@export JWT_SECRET=test-secret-key-for-testing-only
	@export GOOGLE_CLIENT_ID=test-google-client-id
	@export GOOGLE_CLIENT_SECRET=test-google-client-secret
	@export GOOGLE_CALLBACK_URL=http://localhost:8080/auth/google/callback
	@export FRONTEND_URL=http://localhost:3000

# Clean test artifacts
test-clean:
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@rm -f *.db test_*.db

# Run tests in watch mode (requires entr: brew install entr)
test-watch:
	@echo "Running tests in watch mode..."
	@find . -name "*.go" | entr -c make test-handlers

# Development commands
fmt:
	@echo "Formatting Go code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

mod-tidy:
	@echo "Tidying Go modules..."
	@go mod tidy

# Build commands
build:
	@echo "Building application..."
	@go build -o bin/goiter main.go

build-prod:
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/goiter main.go

# Development workflow
dev: mod-tidy fmt vet test-handlers
	@echo "Development checks complete!"

# CI pipeline simulation
ci-test: mod-tidy fmt vet test test-handlers-coverage
	@echo "CI pipeline tests complete!"

# Help command
help:
	@echo "Available commands:"
	@echo "  start-backend       Start backend with hot reloading"
	@echo "  start-backend-no-air Start backend without hot reloading"
	@echo "  stop-backend        Stop backend server"
	@echo ""
	@echo "Testing:"
	@echo "  test               Run all tests"
	@echo "  test-handlers      Run handler tests"
	@echo "  test-handlers-coverage Run handler tests with coverage"
	@echo "  test-unit          Run unit tests"
	@echo "  test-integration   Run integration tests"
	@echo "  test-auth          Run authentication tests"
	@echo "  test-profile       Run profile handler tests"
	@echo "  test-account       Run account handler tests"
	@echo "  test-plan          Run plan handler tests"
	@echo "  test-billing       Run billing handler tests"
	@echo "  test-legacy        Run legacy test suite"
	@echo "  test-watch         Run tests in watch mode"
	@echo "  test-clean         Clean test artifacts"
	@echo ""
	@echo "Development:"
	@echo "  fmt                Format Go code"
	@echo "  vet                Run go vet"
	@echo "  lint               Run golangci-lint"
	@echo "  mod-tidy           Tidy Go modules"
	@echo "  build              Build application"
	@echo "  dev                Run development checks"
	@echo "  ci-test            Run CI pipeline simulation"
	@echo ""
	@echo "Database:"
	@echo "  db                 Connect to database"
	@echo "  clean              Clean database"
