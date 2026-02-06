# Load variables from .envs/local.env
include .envs/local.env
export $(shell sed 's/=.*//' .envs/local.env)

# App
export APP_BINARY_PATH ?= /tmp/tiny-url

.PHONY: new-migration

dev:
	@go run cmd/api/main.go

test:
	@go test -coverprofile=coverage.out ./...

build:
	@go build -v -o $(APP_BINARY_PATH) cmd/api/main.go
	@chmod -X $(APP_BINARY_PATH)

run: build
	@$(APP_BINARY_PATH)

api-create-urls:
	@echo "Usage example, make api-create-urls"
	@curl -X POST "localhost:8080/api/v1/url/" --json '{"target":"https://google.com"}'

api-get-by-id:
	@echo "Usage example, make api-get-by-id id=1"
	@curl "localhost:8080/api/v1/url/$(id)" -H "Content-Type: application/json"

api-list-urls:
	@echo "Usage example, make api-list-urls"
	@curl "localhost:8080/api/v1/url/" -H "Content-Type: application/json"

new-migration:
	@echo "This command requires goloang migrate; https://github.com/golang-migrate/migrate"
	@echo "Usage example, make new-migration name=create_urls_table"
	@migrate create -ext sql -dir migration -seq $(name)

migrate:
	@migrate -database "postgres://postgres:postgres@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -path ./migration/ up
