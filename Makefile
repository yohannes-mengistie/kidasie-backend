.RECIPEPREFIX := >

-include .env
export

MIGRATE_BIN := $(shell go env GOPATH)/bin/migrate
MIGRATIONS_DIR := migrations

.PHONY: run fmt test vet check
.PHONY: db-up db-down db-status
.PHONY: migrate-up migrate-down migration seed
.PHONY: import-content
.PHONY: publish-content
.PHONY: prepare-st-mary import-st-mary publish-st-mary
.PHONY: integrate-st-mary




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

publish-content:
>@test -n "$(slug)" || (echo "usage: make publish-content slug=example [allow_review=true]" && exit 1)
>@go run ./cmd/publishcontent \
>  -slug "$(slug)" \
>  -confirm "$(slug)" \
>  $(if $(filter true,$(allow_review)),-allow-review-required,)

ST_MARY_SOURCE ?= content/generated/st-mary-slides.json
ST_MARY_IMPORT ?= content/generated/st-mary-import.json

prepare-st-mary:
>@go run ./cmd/convertslides \
>  -file "$(ST_MARY_SOURCE)" \
>  -out "$(ST_MARY_IMPORT)" \
>  -slug "st-mary" \
>  -name "Anaphora of St. Mary" \
>  -name-am "የእመቤታችን የድንግል ማርያም ቅዳሴ" \
>  -section-title "Complete Liturgy and Anaphora of St. Mary" \
>  -section-title-am "ሙሉ ሥርዓተ ቅዳሴና የእመቤታችን ቅዳሴ"

import-st-mary: prepare-st-mary
>@go run ./cmd/importcontent -file "$(ST_MARY_IMPORT)"

publish-st-mary:
>@go run ./cmd/publishcontent \
>  -slug "st-mary" \
>  -confirm "st-mary" \
>  -allow-review-required

integrate-st-mary: import-st-mary publish-st-mary
