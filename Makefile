recreate_db:
	@docker compose down -v
	@docker compose up -d
	@sleep 2

build:
	@go build -o bin/main  steps/main.go