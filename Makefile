.PHONY: help proto buf-generate buf-lint buf-format buf-breaking sqlc run build docker-build docker-up docker-down docker-logs deps migrate-up migrate-down migrate-status migrate-new migrate-reset seed seed-clear

# Default target
help:
	@echo "Available commands:"
	@echo ""
	@echo "  Protocol Buffers:"
	@echo "    buf-generate     Generate Go code from proto files"
	@echo "    buf-lint         Lint proto files"
	@echo "    buf-format       Format proto files"
	@echo "    buf-breaking     Check for breaking changes"
	@echo "    proto            Legacy protoc generation"
	@echo ""
	@echo "  Database:"
	@echo "    sqlc             Generate Go code from SQL queries"
	@echo "    migrate-up       Apply pending migrations"
	@echo "    migrate-down     Rollback last migration"
	@echo "    migrate-status   Show migration status"
	@echo "    migrate-new      Create new migration (usage: make migrate-new name=migration_name)"
	@echo "    migrate-reset    Reset database (rollback all, then apply all)"
	@echo "    seed             Seed database with sample companies"
	@echo "    seed-clear       Clear seed data from database"
	@echo ""
	@echo "  Development:"
	@echo "    run              Run the server locally"
	@echo "    build            Build the server binary"
	@echo "    deps             Update Go dependencies"
	@echo ""
	@echo "  Docker:"
	@echo "    docker-build     Build Docker image"
	@echo "    docker-up        Start containers with build"
	@echo "    docker-down      Stop containers"
	@echo "    docker-logs      View container logs"
	@echo "    docker-clean     Stop containers and remove volumes"
	@echo ""
	@echo "  Prerequisites for migrations:"
	@echo "    export DATABASE_URL=postgres://admin:password@localhost:5432/saas_microcourses?sslmode=disable"

# Buf commands (recommended)
buf-generate:
	buf generate

buf-lint:
	buf lint

buf-format:
	buf format -w

buf-breaking:
	buf breaking --against '.git#branch=main'

# Legacy protoc command (kept for compatibility)
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/company.proto

sqlc:
	sqlc generate

run:
	go run cmd/server/main.go

build:
	go build -o server cmd/server/main.go

docker-build:
	docker build -t saas-company-service .

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v
	docker rmi saas-company-service || true

deps:
	go mod tidy

# Database migration commands (requires dbmate and DATABASE_URL)
# Install dbmate: curl -fsSL -o ~/.local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64 && chmod +x ~/.local/bin/dbmate
# Example: export DATABASE_URL=postgres://admin:password@localhost:5432/saas_microcourses?sslmode=disable

# Define dbmate binary path (check common locations)
DBMATE := $(shell which dbmate 2>/dev/null || echo ~/.local/bin/dbmate)

migrate-up:
	$(DBMATE) up

migrate-down:
	$(DBMATE) rollback

migrate-status:
	$(DBMATE) status

migrate-new:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-new name=migration_name"; \
		echo "Example: make migrate-new name=add_users_table"; \
		exit 1; \
	fi
	$(DBMATE) new $(name)

migrate-reset:
	$(DBMATE) down && $(DBMATE) up

# Seed commands
seed:
	@echo "Seeding database with sample companies..."
	go run cmd/seed/main.go

seed-clear:
	@echo "Clearing seed data from database..."
	go run cmd/seed/main.go --clear