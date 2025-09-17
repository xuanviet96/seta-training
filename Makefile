SHELL := powershell.exe

.PHONY: tidy build run clean test

tidy:
	go mod tidy

build:
	go build -o bin/server.exe ./cmd/server

run:
	go run ./cmd/server

clean:
	if (Test-Path bin) { Remove-Item -Recurse -Force bin }

test:
	go test -v ./...

# Docker commands
up:
	docker-compose up -d --build

down:
	docker-compose down -v

logs:
	docker-compose logs -f

# Database migration helpers (requires PostgreSQL client tools)
migrate-up:
	@echo "Please run the migration SQL manually: migrations/0001_init_sql"

migrate-down:
	@echo "Please run migration rollback manually if needed"