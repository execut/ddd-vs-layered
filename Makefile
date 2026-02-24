recreate_db:
	@docker compose down -v
	@docker compose up -d
	@sleep 2

build:
	@go build -o bin/main  steps/presentation/cmd/main.go

all_checks: recreate_db build
	USER_ID="323e4567-e89b-12d3-a456-426655440001"
	DATABASE_URL=postgres://postgres@localhost:15432/labels go test ./...
	golangci-lint run ./steps/...