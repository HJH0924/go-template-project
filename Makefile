.PHONY: help install-deps install-hooks format lint build run test clean generate docs

# Variables
BINARY_NAME=server
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"
CONFIG_FILE=configs/config.yaml

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-deps: ## Check and install development dependencies
	@command -v golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest)
	@command -v gofumpt > /dev/null || (echo "Installing gofumpt..." && go install mvdan.cc/gofumpt@latest)
	@command -v goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	@command -v goimports-reviser > /dev/null || (echo "Installing goimports-reviser..." && go install github.com/incu6us/goimports-reviser/v3@latest)
	@command -v buf > /dev/null || (echo "Installing buf..." && go install github.com/bufbuild/buf/cmd/buf@latest)
	@$(MAKE) install-hooks
	@echo ""
	@echo "✅ All development dependencies installed successfully!"

install-hooks: ## Install git pre-commit hooks
	@if [ -d .git ]; then \
		echo "Installing pre-commit hook..."; \
		mkdir -p .git/hooks; \
		cp scripts/pre-commit .git/hooks/pre-commit; \
		chmod +x .git/hooks/pre-commit; \
		echo "✅ Pre-commit hook installed successfully!"; \
	else \
		echo "⚠️  Not a git repository, skipping hook installation"; \
	fi

format: ## Format code
	gofumpt -l -w .
	goimports-reviser -imports-order=std,project,company,general  -recursive ./

lint: ## Run linter
	golangci-lint run -c .golangci.toml ./...

build: format lint ## Build the binary
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/server

run: build ## Run the server
	./bin/$(BINARY_NAME) serve --config $(CONFIG_FILE)

dev: ## Run the server in development mode
	$(GO) run ./cmd/server serve --config $(CONFIG_FILE)

test: ## Run tests
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

test-short: ## Run tests without race detector
	$(GO) test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf gen/
	rm -f coverage.out coverage.html
	$(GO) clean

generate: ## Generate protobuf code
	buf generate

docs: ## Run docs
	cd docs && pnpm run docs:dev

.DEFAULT_GOAL := help
