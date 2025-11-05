.PHONY: help build run test clean docker-build docker-up docker-down

# Variables
BINARY_NAME=todo-app
GO=go
DOCKER_IMAGE=todo-app-fiber:latest

help:
	@echo "Todo API - Fiber Framework"
	@echo ""
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run in development mode with auto-reload"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-up     - Start containers with docker-compose"
	@echo "  docker-down   - Stop containers with docker-compose"
	@echo "  migrate       - Run database migrations"

deps:
	$(GO) mod download
	$(GO) mod tidy

build: deps
	$(GO) build -o $(BINARY_NAME) ./cmd/main.go
	@echo "Build complete: ./$(BINARY_NAME)"

run: build
	./$(BINARY_NAME)

dev: deps
	@command -v air >/dev/null 2>&1 || go install github.com/cosmtrek/air@latest
	air

test: deps
	$(GO) test -v -cover ./...

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -rf bin/

fmt:
	$(GO) fmt ./...

lint: deps
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

migrate:
	@echo "Running migrations..."
	psql $${DATABASE_URL:-postgres://postgres:password@localhost:5432/todo_db} < migrations/01_create_todos_table.sql
	@echo "Migrations complete"

.DEFAULT_GOAL := help
