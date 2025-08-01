.PHONY: proto buf-generate buf-lint buf-format buf-breaking sqlc run build docker-build docker-up docker-down docker-logs deps

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