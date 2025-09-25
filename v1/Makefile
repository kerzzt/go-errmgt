.PHONY: help test test-race test-cover lint fmt vet clean build deps tidy

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint
	golangci-lint run

deps: ## Download dependencies
	go mod download

tidy: ## Tidy go modules
	go mod tidy

# Testing
test: ## Run tests
	go test -v ./...

test-race: ## Run tests with race detector
	go test -v -race ./...

test-cover: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

test-examples: ## Run example tests
	go test -v -run Example ./...

# Build
build: ## Build the library
	go build ./...

# Cleanup
clean: ## Clean build artifacts
	rm -f coverage.out coverage.html
	go clean -testcache
	go clean -modcache

# CI targets
ci-test: deps test-race ## Run CI tests
ci-lint: deps lint ## Run CI linting
ci-build: deps build ## Run CI build

# All checks (like CI)
check: fmt vet lint test-race ## Run all checks

# Development setup
dev-setup: ## Set up development environment
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go mod download