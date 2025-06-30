.PHONY: start stop start-backend stop-backend start-backend-no-air check-port clean-air

# Default target
all: start

# Start both servers
start: start-backend start-frontend

# Stop both servers
stop: stop-backend stop-frontend

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
	@cd backend && air

# Start backend server without hot reloading
start-backend-no-air: 
	@echo "Starting backend server without hot reloading..."
	@cd backend && go run main.go

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
