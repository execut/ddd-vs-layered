recreate_db:
	@docker compose down -v
	@docker compose up -d
	@sleep 2

build:
	@go build -o bin/main  steps/main.go

all_checks: recreate_db build
	DATABASE_URL=postgres://postgres@localhost:15432/labels go test ./...
	golangci-lint run ./steps/...