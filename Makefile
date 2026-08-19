.RECIPEPREFIX := >

-include .env
export

MIGRATE_BIN := $(shell go env GOPATH)/bin/migrate
MIGRATIONS_DIR := migrations

.PHONY: run fmt test vet check
.PHONY: db-up db-down db-status
.PHONY: migrate-up migrate-down migration seed
.PHONY: import-content




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
>@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
>@$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migration:
>@test -n "$(name)" || (echo "usage: make migration name=create_example" && exit 1)
>$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

seed:
>docker compose exec -T postgres psql -U "$(POSTGRES_USER)" -d "$(POSTGRES_DB)" < db/seeds/development.sql

import-content:
>@test -n "$(file)" || (echo "usage: make import-content file=content/example.json" && exit 1)
>@go run ./cmd/importcontent -file "$(file)"


.PHONY: extract-slides serve-aligner

SLIDE_PDF ?= source-material/liturgy.pdf
SLIDE_AUDIO ?= source-material/audio/apostles.mp3
SLIDE_OUT ?= content/generated/apostles-slides.json

extract-slides:
>@go run ./cmd/extractslides -pdf "$(SLIDE_PDF)" -audio "$(SLIDE_AUDIO)" -out "$(SLIDE_OUT)"

serve-aligner:
>@python3 -m http.server 4173 --directory tools/content-aligner
