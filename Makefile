.PHONY: run compose-up migrate-up migrate-down

# Путь к схеме миграции
MIGRATION_PATH=./migrations
# Строка подключения к базе данных
DATABASE_URL=postgres://postgres:admin@order-service-db:5432/postgres?sslmode=disable

run: compose-up

compose-up:
	docker-compose -f docker-compose.yaml up --build

migrate-up:
	docker-compose run --rm migrator -path /migrations -database ${DATABASE_URL} up

migrate-down:
	docker-compose run --rm migrator -path /migrations -database ${DATABASE_URL} down