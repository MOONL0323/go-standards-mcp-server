.PHONY: help build run test clean docker-build docker-run install lint

# Variables
BINARY_NAME=mcp-server
BUILD_DIR=bin
DOCKER_IMAGE=go-standards-mcp-server
DOCKER_TAG=latest

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build MCP server
	@echo "Building MCP server..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-cli: ## Build CLI tool
	@echo "Building CLI tool..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/go-standards ./cmd/cli
	@echo "Build complete: $(BUILD_DIR)/go-standards"

build-all: build build-cli ## Build both MCP server and CLI tool

run: ## Run the application in stdio mode
	@echo "Running $(BINARY_NAME) in stdio mode..."
	go run ./cmd/server

run-http: ## Run the application in HTTP mode
	@echo "Running $(BINARY_NAME) in HTTP mode..."
	go run ./cmd/server --mode http --port 8080

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -v -tags=integration ./tests/integration/...

lint: ## Run linters
	@echo "Running linters..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -rf tmp/
	rm -rf reports/
	rm -f coverage.out coverage.html
	@echo "Clean complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --name $(BINARY_NAME) $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-compose-up: ## Start services with docker-compose
	@echo "Starting services with docker-compose..."
	docker-compose up -d
	@echo "Services started"

docker-compose-down: ## Stop services with docker-compose
	@echo "Stopping services with docker-compose..."
	docker-compose down
	@echo "Services stopped"

docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

setup-dev: install ## Setup development environment
	@echo "Setting up development environment..."
	@command -v golangci-lint > /dev/null || (echo "Installing golangci-lint..." && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin)
	@echo "Development environment ready"

check: lint test ## Run all checks (lint + test)

all: clean install lint test build ## Run all tasks

.DEFAULT_GOAL := help
