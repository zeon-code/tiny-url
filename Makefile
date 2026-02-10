# Load variables from .envs/local.env
include .envs/local.env

export $(shell sed 's/=.*//' .envs/local.env)
APP_VERSION ?= $(shell git describe --abbrev=0 --tags 2>/dev/null || echo v0.0.1)

# App
export APP_BINARY_PATH ?= /tmp/tiny-url

.PHONY: new-migration

infrastructure:
	@docker-compose up -d

dev: infrastructure
	@go run cmd/api/main.go

run:
	@docker-compose --profile app up --build --force-recreate -d
	@docker-compose logs api -f

test:
	@go test -coverprofile=coverage.out ./...

build:
	@echo "Building app version $(APP_VERSION)"
	@go build -o $(APP_BINARY_PATH) -ldflags "-X main.version=$(APP_VERSION)" cmd/api/main.go

api-create-urls:
	@echo "\nUsage example, make api-create-urls"
	@curl -X POST "localhost:8080/api/v1/url/" --json '{"target":"https://google.com"}'

api-get-by-id:
	@echo "\nUsage example, make api-get-by-id id=1"
	@curl "localhost:8080/api/v1/url/$(id)" -H "Content-Type: application/json"

api-list-urls:
	@echo "\nUsage example, make api-list-urls"
	@curl "localhost:8080/api/v1/url/" -H "Content-Type: application/json"

new-migration:
	@echo "This command requires goloang migrate; https://github.com/golang-migrate/migrate"
	@echo "Usage example, make new-migration name=create_urls_table"
	@migrate create -ext sql -dir migration -seq $(name)

migrate:
	@migrate -database "postgres://postgres:postgres@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -path ./migration/ up
