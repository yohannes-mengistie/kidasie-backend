.RECIPEPREFIX := >

-include .env
export

MIGRATE_BIN := $(shell go env GOPATH)/bin/migrate
MIGRATIONS_DIR := migrations

.PHONY: run fmt test vet check
.PHONY: db-up db-down db-status
.PHONY: migrate-up migrate-down migration seed



run:
>go run ./cmd/api

fmt:
>go fmt ./...

test:
>go test ./...

vet:
>go vet ./...

check: fmt test vet

db-up:
>docker compose up -d --wait postgres

db-down:
>docker compose down

db-status:
>docker compose ps

migrate-up:
>$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
>$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migration:
>@test -n "$(name)" || (echo "usage: make migration name=create_example" && exit 1)
>$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

seed:
>docker compose exec -T postgres psql -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)" < db/seeds/development.sql
