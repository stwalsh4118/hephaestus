FRONTEND_DIR := frontend
BACKEND_DIR := backend

.DEFAULT_GOAL := help

.PHONY: help dev down dev-down build lint test clean

help: ## Show available make targets
	@printf "Available targets:\n"
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_-]+:.*## / {printf "  %-10s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Start all services with Docker Compose
	docker compose up --build

down: ## Stop all Docker Compose services
	docker compose down

dev-down: down ## Alias for down

build: ## Build frontend and backend
	pnpm --dir $(FRONTEND_DIR) build
	cd $(BACKEND_DIR) && go build ./...

lint: ## Run frontend and backend linters
	pnpm --dir $(FRONTEND_DIR) lint
	cd $(BACKEND_DIR) && golangci-lint run ./...

test: ## Run frontend and backend tests
	pnpm --dir $(FRONTEND_DIR) test
	cd $(BACKEND_DIR) && go test ./...

clean: ## Remove generated build artifacts and caches
	rm -rf $(FRONTEND_DIR)/.next $(BACKEND_DIR)/tmp
	go clean -cache
