recreate_db:
	@docker compose down -v
	@docker compose up -d
	@sleep 2

build:
	@go build -o bin/1-layered  steps/1-layered/main.go
	@go build -o bin/2-ddd-event-sourcing  steps/2-ddd-event-sourcing/main.go