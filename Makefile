.PHONY: help build run test docker-up docker-down clean install dev

# Variables
APP_NAME=phoenix-api
DOCKER_COMPOSE=docker-compose
GO=go

help: ## Show this help message
	@echo "PHOENIX SEO Platform - Makefile Commands"
	@echo "=========================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install Go dependencies
	$(GO) mod download
	$(GO) mod tidy

build: ## Build the application
	$(GO) build -o $(APP_NAME) ./cmd/api

run: ## Run the application locally
	$(GO) run ./cmd/api/main.go

dev: ## Run in development mode with auto-reload (requires air)
	air

test: ## Run tests
	$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Run tests with coverage report
	$(GO) tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

docker-build: ## Build Docker image
	docker build -t phoenix-seo-api:latest .

docker-up: ## Start all services with Docker Compose
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop all services
	$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker logs
	$(DOCKER_COMPOSE) logs -f

docker-restart: ## Restart all services
	$(DOCKER_COMPOSE) restart

docker-clean: ## Remove all containers and volumes
	$(DOCKER_COMPOSE) down -v

db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	@docker exec -i phoenix-postgres psql -U postgres -d phoenix_seo < internal/database/migrations/001_initial_schema.sql

db-shell: ## Open PostgreSQL shell
	docker exec -it phoenix-postgres psql -U postgres -d phoenix_seo

redis-cli: ## Open Redis CLI
	docker exec -it phoenix-redis redis-cli

clean: ## Clean build artifacts
	rm -f $(APP_NAME)
	rm -f coverage.txt coverage.html
	$(GO) clean

fmt: ## Format Go code
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

security: ## Run security checks (requires gosec)
	gosec ./...

all: clean install build test ## Run all build steps
