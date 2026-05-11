# ============================================================
# Fish Game Management System - Makefile
# ============================================================
# Prerequisites: Docker, Docker Compose, Go 1.21+, Node.js 20+
# ============================================================

.PHONY: help up down logs build dev-backend dev-frontend \
        lint test clean db-shell redis-shell backend-shell

# ─── Variables ───────────────────────────────────────────────
DOCKER_COMPOSE  = docker compose
BACKEND_DIR     = backend-api
FRONTEND_DIR    = .
GO_BINARY       = fish-game-api
GO_MAIN         = ./cmd/api/main.go

# Colors
GREEN  = \033[0;32m
YELLOW = \033[0;33m
CYAN   = \033[0;36m
RESET  = \033[0m

# ─── Help ─────────────────────────────────────────────────────
help: ## Show this help message
	@echo ""
	@echo "$(CYAN)🐟 Fish Game Management System$(RESET)"
	@echo "$(YELLOW)────────────────────────────────────────$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)  %-20s$(RESET) %s\n", $$1, $$2}'
	@echo ""

# ─── Docker Compose ───────────────────────────────────────────
up: ## Start all services (detached)
	@echo "$(CYAN)▶ Starting Fish Game services...$(RESET)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)✓ All services started$(RESET)"
	@echo ""
	@echo "  Frontend Admin:  http://localhost:3000"
	@echo "  Backend API:     http://localhost:8080"
	@echo "  Health Check:    http://localhost:8080/api/v1/health"

down: ## Stop all services
	@echo "$(YELLOW)▶ Stopping Fish Game services...$(RESET)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)✓ All services stopped$(RESET)"

build: ## Build all Docker images
	@echo "$(CYAN)▶ Building Docker images...$(RESET)"
	$(DOCKER_COMPOSE) build --no-cache
	@echo "$(GREEN)✓ Build complete$(RESET)"

logs: ## Follow logs from all services
	$(DOCKER_COMPOSE) logs -f

restart: ## Restart all services
	$(DOCKER_COMPOSE) restart

# ─── Local Development ────────────────────────────────────────
dev-backend: ## Run backend in local dev mode (with hot reload via air)
	@echo "$(CYAN)▶ Starting backend dev server...$(RESET)"
	cd $(BACKEND_DIR) && go run $(GO_MAIN)

dev-frontend: ## Run frontend dev server
	@echo "$(CYAN)▶ Starting frontend dev server on :3000...$(RESET)"
	npm run dev --prefix $(FRONTEND_DIR)

dev: ## Start both frontend & backend locally (parallel)
	@echo "$(CYAN)▶ Starting local development environment...$(RESET)"
	$(MAKE) -j2 dev-backend dev-frontend

# ─── Go Backend ───────────────────────────────────────────────
build-backend: ## Build Go binary
	@echo "$(CYAN)▶ Building Go backend...$(RESET)"
	cd $(BACKEND_DIR) && CGO_ENABLED=0 go build -o bin/$(GO_BINARY) $(GO_MAIN)
	@echo "$(GREEN)✓ Binary: $(BACKEND_DIR)/bin/$(GO_BINARY)$(RESET)"

test: ## Run all Go tests
	@echo "$(CYAN)▶ Running tests...$(RESET)"
	cd $(BACKEND_DIR) && go test ./... -v -race -cover

lint: ## Run Go vet + staticcheck
	@echo "$(CYAN)▶ Running linter...$(RESET)"
	cd $(BACKEND_DIR) && go vet ./...

tidy: ## Tidy Go module dependencies
	cd $(BACKEND_DIR) && go mod tidy

# ─── Frontend ─────────────────────────────────────────────────
install: ## Install frontend dependencies
	npm install --prefix $(FRONTEND_DIR)

build-frontend: ## Build frontend for production
	npm run build --prefix $(FRONTEND_DIR)

# ─── Database ─────────────────────────────────────────────────
db-shell: ## Open a psql shell inside the postgres container
	$(DOCKER_COMPOSE) exec postgres psql -U fishgame -d fishgame_db

db-seed: ## Re-run the seed SQL manually
	$(DOCKER_COMPOSE) exec -T postgres psql -U fishgame -d fishgame_db < database/seed.sql

redis-shell: ## Open redis-cli in the redis container
	$(DOCKER_COMPOSE) exec redis redis-cli

# ─── Shell access ─────────────────────────────────────────────
backend-shell: ## Open a shell in the backend container
	$(DOCKER_COMPOSE) exec backend-api sh

# ─── Cleanup ──────────────────────────────────────────────────
clean: ## Remove build artifacts and stopped containers
	@echo "$(YELLOW)▶ Cleaning up...$(RESET)"
	cd $(BACKEND_DIR) && rm -rf bin/
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	@echo "$(GREEN)✓ Clean complete$(RESET)"

clean-all: clean ## Full clean including Docker images
	$(DOCKER_COMPOSE) down --rmi all --volumes --remove-orphans
