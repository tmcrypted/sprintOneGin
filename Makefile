DB_HOST := localhost
DB_PORT := 54320
DB_USER := db_user
DB_PASSWORD := pwd123
DB_NAME := db_test

GOOSE_DRIVER := "postgres"
SERVICE_NAME := pgdb

GOOSE_DBSTRING := "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable"



run:
	@go run .\cmd\app\.
	
migrate-up:
	@echo "Применяю миграции..."
	@goose -dir ./internal/migrations/ $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up
	@echo "Миграции применены успешно"

migrate-down:
	@echo "Откатываю последнюю миграцию..."
	@goose -dir ./migrations/sqlmigrations $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down
	@echo "Миграция откачена"