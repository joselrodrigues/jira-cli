# jira-cli Makefile

# Variables
BINARY_NAME=jira-cli
BINARY_DIR=bin
INSTALL_PATH=/usr/local/bin
GO=go
GOFLAGS=-ldflags="-s -w"
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Default target
.DEFAULT_GOAL := build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) .
	@echo "Binary created: $(BINARY_DIR)/$(BINARY_NAME)"

# Build with version info
.PHONY: build-version
build-version:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BINARY_DIR)/$(BINARY_NAME) .
	@echo "Binary created: $(BINARY_DIR)/$(BINARY_NAME)"

# Install to /usr/local/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@cp $(BINARY_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installed successfully!"

# Uninstall from /usr/local/bin
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstalled successfully!"

# Run tests
.PHONY: test
test:
	$(GO) test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Download dependencies
.PHONY: deps
deps:
	$(GO) mod download
	$(GO) mod tidy

# Format code
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# Lint code
.PHONY: lint
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BINARY_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-linux-arm64 .
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Binaries created in $(BINARY_DIR)/"

# Run the CLI (for development)
.PHONY: run
run:
	$(GO) run . $(ARGS)

# Show help
.PHONY: help
help:
	@echo "jira-cli Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build          Build the binary"
	@echo "  make install        Build and install to $(INSTALL_PATH)"
	@echo "  make uninstall      Remove from $(INSTALL_PATH)"
	@echo "  make test           Run tests"
	@echo "  make test-coverage  Run tests with coverage report"
	@echo "  make clean          Remove build artifacts"
	@echo "  make deps           Download and tidy dependencies"
	@echo "  make fmt            Format code"
	@echo "  make lint           Run linter"
	@echo "  make build-all      Build for macOS, Linux, Windows"
	@echo "  make run ARGS='...' Run CLI with arguments"
	@echo "  make help           Show this help"