USER := postgres
PASSWORD := password
HOST := localhost
PORT := 5432
DB := simple-bank-db
MIGRATION_DIR := db/migration
MIGRATION_NAME := init_schema

migrate_up: 
	@echo "Migrating up..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://$(USER):$(PASSWORD)@$(HOST):$(PORT)/$(DB)?sslmode=disable" -verbose up
	@echo "Migration up completed"

migrate_down:
	@echo "Migrating down..."
	migrate -path $(MIGRATION_DIR) -database "postgresql://$(USER):$(PASSWORD)@$(HOST):$(PORT)/$(DB)?sslmode=disable" -verbose down
	@echo "Migration down completed"

migrate_create:
	@echo "Creating migration..."
	migrate create -ext sql -dir db/migration -seq $(MIGRATION_NAME)
	@echo "Migration created"

sqlc:
	@echo "Generating sqlc..."
	sqlc generate
	@echo "sqlc generated"

test: 
	go test -v -cover ./... 
