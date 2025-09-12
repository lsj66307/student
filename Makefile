# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Binary names
BINARY_NAME=student-api
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe

# Directories
CMD_DIR=./cmd/student-api
BUILD_DIR=./build
COVERAGE_DIR=./coverage

# Default target
.PHONY: all
all: test build

# Build the application
.PHONY: build
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(CMD_DIR)

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v $(CMD_DIR)

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_WINDOWS) -v $(CMD_DIR)

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-windows build

# Run the application
.PHONY: run
run:
	@echo "Running application..."
	$(GOCMD) run $(CMD_DIR)/main.go

# Run with specific config
.PHONY: run-dev
run-dev:
	@echo "Running in development mode..."
	APP_CONFIG=./configs/config.yaml $(GOCMD) run $(CMD_DIR)/main.go

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"

# Run unit tests only
.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v ./test/unit/...

# Run integration tests only
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v ./test/integration/...

# Run benchmark tests
.PHONY: test-benchmark
test-benchmark:
	@echo "Running benchmark tests..."
	$(GOTEST) -v -bench=. ./test/benchmark/...

# Run tests with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -v -race ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOGET) -d -v ./...

# Update dependencies
.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	$(GOGET) -u -d -v ./...
	$(GOMOD) tidy

# Verify dependencies
.PHONY: deps-verify
deps-verify:
	@echo "Verifying dependencies..."
	$(GOMOD) verify

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Vet code
.PHONY: vet
vet:
	@echo "Vetting code..."
	$(GOVET) ./...

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Linting code..."
	$(GOLINT) run

# Install golangci-lint
.PHONY: install-lint
install-lint:
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

# Check code quality
.PHONY: check
check: fmt vet lint test

# Generate Swagger documentation
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g $(CMD_DIR)/main.go -o ./docs

# Install Swagger generator
.PHONY: install-swagger
install-swagger:
	@echo "Installing Swagger generator..."
	$(GOGET) github.com/swaggo/swag/cmd/swag

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t student-management-system .

# Docker run
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 student-management-system

# Database migration (placeholder)
.PHONY: migrate-up
migrate-up:
	@echo "Running database migrations..."
	# Add your migration commands here

# Database rollback (placeholder)
.PHONY: migrate-down
migrate-down:
	@echo "Rolling back database migrations..."
	# Add your rollback commands here

# Setup development environment
.PHONY: setup-dev
setup-dev: deps install-lint install-swagger
	@echo "Development environment setup complete!"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-windows  - Build for Windows"
	@echo "  build-all      - Build for all platforms"
	@echo "  run            - Run the application"
	@echo "  run-dev        - Run in development mode"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-benchmark - Run benchmark tests"
	@echo "  test-race      - Run tests with race detection"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Download dependencies"
	@echo "  deps-update    - Update dependencies"
	@echo "  deps-verify    - Verify dependencies"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  lint           - Lint code"
	@echo "  check          - Run all code quality checks"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  setup-dev      - Setup development environment"
	@echo "  help           - Show this help message"