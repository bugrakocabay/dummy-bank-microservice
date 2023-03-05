GATEWAY_BINARY=gatewayApp
ACCOUNT_SERVICE_BINARY=accountServiceApp
LOGGER_SERVICE_BINARY=loggerServiceApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_gateway build_account_service build_logger_service
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## account_db_migrate_up: creates a migration to account database
account_db_migrate_up:
	migrate -path ./../account-service/db/migration -database "postgresql://postgres:postgres@localhost:5432/accounts?sslmode=disable" -verbose up

## account_db_migrate_down: deletes the migration for account database
account_db_migrate_down:
	migrate -path ./../account-service/db/migration -database "postgresql://postgres:postgres@localhost:5432/accounts?sslmode=disable" -verbose down

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_gateway: builds the gateway binary as a linux executable
build_gateway:
	@echo "Building gateway binary..."
	cd ../gateway && env GOOS=linux CGO_ENABLED=0 go build -o ${GATEWAY_BINARY} ./cmd/api
	@echo "Done!"

## build_account-service: builds the account-service binary as a linux executable
build_account_service:
	@echo "Building account-service binary..."
	cd ../account-service && env GOOS=linux CGO_ENABLED=0 go build -o ${ACCOUNT_SERVICE_BINARY} ./cmd/api
	@echo "Done!"

## build_logger-service: builds the logger-service binary as a linux executable
build_logger_service:
	@echo "Building logger-service binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_SERVICE_BINARY} ./cmd/api
	@echo "Done!"

## sqlc: generates sql code from template
sqlc:
	sqlc generate