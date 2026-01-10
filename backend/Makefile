.PHONY: help build run test test-integration clean migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build            - Build binaries"
	@echo "  make run              - Run API server"
	@echo "  make run-indexer      - Run indexer"
	@echo "  make test             - Run unit tests"
	@echo "  make test-integration - Run integration tests"
	@echo "  make migrate-up       - Run migrations"
	@echo "  make migrate-down     - Rollback migrations"
	@echo "  make clean            - Clean build artifacts"

build:
	@echo "Building..."
	go build -o bin/api cmd/api/main.go
	go build -o bin/indexer cmd/indexer/main.go

run:
	@echo "Starting API server..."
	go run cmd/api/main.go

run-indexer:
	@echo "Starting indexer..."
	go run cmd/indexer/main.go

test:
	@echo "Running unit tests..."
	go test -v -cover ./internal/...

test-integration:
	@echo "Running integration tests..."
	go test -v ./test/...

migrate-up:
	@echo "Running migrations..."
	migrate -path internal/database/migrations \
		-database "postgresql://eth_user:eth_pass_dev_only@localhost:5432/ethereum_explorer?sslmode=disable" up

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path internal/database/migrations \
		-database "postgresql://eth_user:eth_pass_dev_only@localhost:5432/ethereum_explorer?sslmode=disable" down 1

clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

lint:
	@echo "Linting..."
	golangci-lint run

fmt:
	@echo "Formatting..."
	go fmt ./...