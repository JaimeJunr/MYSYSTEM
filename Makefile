.PHONY: build run install clean test test-verbose test-coverage test-short test-integration test-smoke help tidy

# Build configuration
BINARY_NAME=homestead
BUILD_DIR=.
CMD_DIR=./cmd/homestead
COVERAGE_FILE=coverage.out

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "✅ Build complete: ./$(BINARY_NAME)"

# Run the application
run:
	@echo "🚀 Running $(BINARY_NAME)..."
	go run $(CMD_DIR)

# Install to $GOPATH/bin
install:
	@echo "📦 Installing $(BINARY_NAME)..."
	go install $(CMD_DIR)
	@echo "✅ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	go clean
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "🧪 Running tests (verbose)..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)
	@echo ""
	@echo "📊 Coverage report saved to $(COVERAGE_FILE)"
	@echo "🌐 To view HTML report, run: go tool cover -html=$(COVERAGE_FILE)"

# Run only short/unit tests (skip integration)
test-short:
	@echo "🧪 Running short tests..."
	go test -short ./...

# Run only integration tests
test-integration:
	@echo "🧪 Running integration tests..."
	go test -v -run Integration ./...

# Black-box smoke: build binary and run -version (skipped with -short)
test-smoke:
	@echo "🧪 Running smoke tests..."
	go test -v -count=1 -run TestSmoke_ .

# Run tests and open coverage in browser
test-coverage-html: test-coverage
	@echo "🌐 Opening coverage report in browser..."
	go tool cover -html=$(COVERAGE_FILE)

# Benchmark tests
benchmark:
	@echo "⚡ Running benchmarks..."
	go test -bench=. -benchmem ./...

# Download and tidy dependencies
tidy:
	@echo "📚 Tidying dependencies..."
	go mod tidy
	go mod verify
	@echo "✅ Dependencies updated"

# Show help
help:
	@echo "Homestead - Makefile commands:"
	@echo ""
	@echo "Build & Run:"
	@echo "  make build              - Build the binary"
	@echo "  make run                - Run the application"
	@echo "  make install            - Install to GOPATH/bin"
	@echo "  make clean              - Remove build artifacts"
	@echo ""
	@echo "Testing:"
	@echo "  make test               - Run all tests"
	@echo "  make test-verbose       - Run tests with verbose output"
	@echo "  make test-coverage      - Run tests with coverage report"
	@echo "  make test-coverage-html - Run tests and open coverage in browser"
	@echo "  make test-short         - Run only unit tests (skip integration and smoke)"
	@echo "  make test-integration   - Run only integration tests"
	@echo "  make test-smoke         - Build CLI and verify -version (E2E-style smoke)"
	@echo "  make benchmark          - Run benchmark tests"
	@echo ""
	@echo "Dependencies:"
	@echo "  make tidy               - Update and verify dependencies"
	@echo ""
	@echo "Help:"
	@echo "  make help               - Show this help"
