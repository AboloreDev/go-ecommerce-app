.PHONY: help build run dev lint migrate-up migrate-down docker-up docker-down

help: 
	@echo "Available commands:"
	@echo "  build       Build the server binary"
	@echo "  run         Run the server"
	@echo "  dev         Run the server in development mode"
	@echo "  lint        Run linters"
	@echo "  migrate-up  Apply database migrations"
	@echo "  migrate-down Rollback database migrations"
	@echo "  docker-up   Start running the container"
	@echo "  docker-down   Rollback running of the container"

build: 
	go build -o bin/server ./cmd/api

run: 
	go run ./cmd/api

dev:
	go run ./cmd/api

lint: 
	golangci-lint run ./...

migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/armory_db?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/armory_db?sslmode=disable" down

docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down -d
