create_db:
	docker-compose exec postgres psql "user=postgres password=postgres" -c "create database development;"
drop_db:
	docker-compose exec postgres psql "user=postgres password=postgres" -c "drop database development;"
reset_db: drop_db create_db

DB_PATH?=postgres://postgres:postgres@localhost:5442/development?sslmode=disable
MIGRATIONS_PATH?=./internal/server/repository/postgres/internal/migrations

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_PATH) up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_PATH) down

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database $(DB_PATH) version
