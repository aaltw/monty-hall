# Monty Hall Terminal Application Makefile

.PHONY: build test clean run install dev lint fmt vet help

# Default target
all: build

# Build the application
build:
	@echo "Building Monty Hall application..."
	go build -o monty-hall ./cmd/monty-hall

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f monty-hall
	rm -f coverage.out coverage.html

# Run the application
run: build
	@echo "Starting Monty Hall application..."
	./monty-hall

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Development mode (build and run)
dev: build run

# Lint the code
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet the code
vet:
	@echo "Vetting code..."
	go vet ./...

# Check code quality (fmt, vet, test)
check: fmt vet test

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o dist/monty-hall-linux-amd64 ./cmd/monty-hall
	GOOS=darwin GOARCH=amd64 go build -o dist/monty-hall-darwin-amd64 ./cmd/monty-hall
	GOOS=darwin GOARCH=arm64 go build -o dist/monty-hall-darwin-arm64 ./cmd/monty-hall
	GOOS=windows GOARCH=amd64 go build -o dist/monty-hall-windows-amd64.exe ./cmd/monty-hall

# Create distribution directory
dist:
	@mkdir -p dist

# Package for distribution
package: dist build-all
	@echo "Creating distribution packages..."
	cd dist && tar -czf monty-hall-linux-amd64.tar.gz monty-hall-linux-amd64
	cd dist && tar -czf monty-hall-darwin-amd64.tar.gz monty-hall-darwin-amd64
	cd dist && tar -czf monty-hall-darwin-arm64.tar.gz monty-hall-darwin-arm64
	cd dist && zip monty-hall-windows-amd64.zip monty-hall-windows-amd64.exe

# Help target
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  run           - Build and run the application"
	@echo "  install       - Install dependencies"
	@echo "  dev           - Development mode (build and run)"
	@echo "  lint          - Run linter (requires golangci-lint)"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  check         - Run fmt, vet, and test"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  package       - Create distribution packages"
	@echo "  help          - Show this help message"