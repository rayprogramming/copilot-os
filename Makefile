# Version Variables
VERSION := dev
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Directories
DIST_DIR := dist
COVERAGE_DIR := coverage

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOFMT := gofmt
GOMOD := $(GOCMD) mod

# Build flags
LDFLAGS := -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.Commit=$(COMMIT)"

.PHONY: help
help:
	@echo "CopilotOS - Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  make build          - Build binary for current platform"
	@echo "  make build-all      - Build binaries for all supported platforms"
	@echo "  make test           - Run unit tests with coverage"
	@echo "  make test-coverage  - Run tests and generate HTML coverage report"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make fmt            - Format code with gofmt"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-push    - Push Docker image to registry"
	@echo "  make dev-setup      - Install development tools"
	@echo "  make ci             - Run full CI pipeline locally"
	@echo ""

.PHONY: build
build:
	@echo "Building copilot-os..."
	@mkdir -p $(DIST_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/copilot-os ./cmd/server
	@echo "✅ Build complete: $(DIST_DIR)/copilot-os"

.PHONY: build-all
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	@echo "  Building Linux x86_64..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/copilot-os-linux-amd64 ./cmd/server
	@echo "  Building Linux ARM64..."
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/copilot-os-linux-arm64 ./cmd/server
	@echo "  Building macOS Intel..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/copilot-os-darwin-amd64 ./cmd/server
	@echo "  Building macOS ARM64..."
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/copilot-os-darwin-arm64 ./cmd/server
	@echo "  Building Windows..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/copilot-os-windows-amd64.exe ./cmd/server
	@echo "✅ Build complete: $(DIST_DIR)/"
	@ls -lh $(DIST_DIR)/

.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@echo ""
	@echo "Coverage Summary:"
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -1

.PHONY: test-coverage
test-coverage: test
	@echo "Generating HTML coverage report..."
	@mkdir -p $(COVERAGE_DIR)
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "✅ Coverage report: $(COVERAGE_DIR)/coverage.html"

.PHONY: lint
lint:
	@echo "Running linters..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run ./... --timeout=5m
	@echo "✅ Linting complete"

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	go mod tidy
	@echo "✅ Code formatted"

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(DIST_DIR) $(COVERAGE_DIR)
	@echo "✅ Clean complete"

.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t copilot-os:$(VERSION) \
		-t copilot-os:latest \
		.
	@echo "✅ Docker build complete"
	@docker images | grep copilot-os

.PHONY: docker-push
docker-push: docker-build
	@echo "Pushing Docker image..."
	@echo "Note: Configure docker registry before pushing"
	# docker push copilot-os:$(VERSION)
	# docker push copilot-os:latest
	@echo "Docker push would target your configured registry"

.PHONY: dev-setup
dev-setup:
	@echo "Setting up development environment..."
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing gosec..."
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "✅ Development setup complete"

.PHONY: ci
ci: lint test build
	@echo "✅ CI pipeline successful"

.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
