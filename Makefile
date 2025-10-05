.PHONY: help build run test clean docker-up docker-down swagger setup start

help:
	@echo "Available commands:"
	@echo "  make docker-up    - Start PostgreSQL database"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-down  - Stop PostgreSQL database"
	@echo "  make swagger      - Generate Swagger documentation"

docker-up:
	docker compose up -d postgres
	@echo "Waiting for database to be ready..."
	@bash -c 'TRIES=30; while ! docker compose exec -T postgres pg_isready -U spycat -d spycat_agency >/dev/null 2>&1; do TRIES=$$((TRIES-1)); if [ $$TRIES -le 0 ]; then echo "DB not ready"; exit 1; fi; sleep 1; done'
	@echo "Database is ready!"

docker-down:
	docker compose down

build:
	go mod tidy
	go build -o bin/spy-cat-agency cmd/main.go

run: docker-up
	@echo "Starting Spy Cat Agency API..."
	go run cmd/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

swagger:
	swag init -g cmd/main.go -o ./docs

setup:
	@echo "Installing swag (OpenAPI generator) if missing..."
	@which swag >/dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger docs..."
	@$(shell go env GOPATH)/bin/swag init -g cmd/main.go -o ./docs || swag init -g cmd/main.go -o ./docs
	$(MAKE) docker-up

start:
	@echo "Starting API on port $${PORT:-3030}..."
	@PORT=$${PORT:-3030} go run cmd/main.go
