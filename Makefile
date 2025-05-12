run_postgresql:
	@echo "Running PostgreSQL container..."
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:latest
	@echo "PostgreSQL container is running."
create_db:
	@echo "Creating database..."
	docker exec -it postgres createdb --username=root --owner=root simple_bank
	@echo "Database created."
drop_db:
	@echo "Dropping database..."
	docker exec -it postgres dropdb simple_bank
	@echo "Database dropped."
db_migrate_up:
	@echo "Running database migrations..."
	migrate -path db/migrations -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up
	@echo "Database migrations completed."
db_migrate_down:
	@echo "Rolling back database migrations..."
	migrate -path db/migrations -database "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down
	@echo "Database migrations rolled back."
sqlc_generate:
	@echo "Running SQLC..."
	sqlc generate
	@echo "SQLC completed."

.PHONY: run_postgresql create_db drop_db db_migrate_up db_migrate_down sqlc_generate