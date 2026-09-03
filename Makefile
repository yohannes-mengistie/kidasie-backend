.RECIPEPREFIX := >

-include .env
export

MIGRATE_BIN := $(shell go env GOPATH)/bin/migrate
MIGRATIONS_DIR := migrations

.PHONY: run fmt test vet check
.PHONY: db-up db-down db-status
.PHONY: migrate-up migrate-down migration seed
.PHONY: import-content
.PHONY: separate-ethiopic
.PHONY: publish-content
.PHONY: prepare-st-mary import-st-mary publish-st-mary
.PHONY: integrate-st-mary
.PHONY: prepare-anaphoras import-anaphoras publish-anaphoras
.PHONY: integrate-anaphoras
.PHONY: prepare-additional-anaphoras import-additional-anaphoras
.PHONY: publish-additional-anaphoras integrate-additional-anaphoras
.PHONY: prepare-liturgy-guide import-liturgy-guide
.PHONY: publish-liturgy-guide integrate-liturgy-guide
.PHONY: integrate-additional-content
.PHONY: import-catalog publish-catalog integrate-catalog
.PHONY: prepare-all-content import-all-content publish-all-content
.PHONY: set-audio
.PHONY: upsert-announcement
.PHONY: export-offline




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

upsert-announcement:
>@test -n "$(file)" || (echo "usage: make upsert-announcement file=content/announcements/example.json confirm=example-announcement" && exit 1)
>@test -n "$(confirm)" || (echo "confirm is required" && exit 1)
>@go run ./cmd/upsertannouncement -file "$(file)" -confirm "$(confirm)"

separate-ethiopic:
>@test -n "$(file)" || (echo "usage: make separate-ethiopic file=content/generated/source.json pdf=source-material/source.pdf [metadata_only=true]" && exit 1)
>@test -n "$(pdf)" || (echo "usage: make separate-ethiopic file=content/generated/source.json pdf=source-material/source.pdf [metadata_only=true]" && exit 1)
>@go run ./cmd/separateethiopic \
>  -file "$(file)" \
>  -pdf "$(pdf)" \
>  $(if $(filter true,$(metadata_only)),-metadata-only,)

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

APOSTLES_SOURCE ?= content/generated/Anaphora-of-the-Apostles.json
APOSTLES_IMPORT ?= content/generated/apostles-import.json
OUR_LORD_SOURCE ?= content/generated/Anaphora-of-Our-Lord-Jesus-Christ.json
OUR_LORD_IMPORT ?= content/generated/our-lord-jesus-christ-import.json
ST_ATHANASIUS_SOURCE ?= content/generated/Anaphora-of-St-Athanasius.json
ST_ATHANASIUS_IMPORT ?= content/generated/st-athanasius-import.json

prepare-anaphoras:
>@go run ./cmd/convertanaphora \
>  -file "$(APOSTLES_SOURCE)" \
>  -out "$(APOSTLES_IMPORT)" \
>  -slug "apostles" \
>  -name "Anaphora of the Apostles" \
>  -name-am "የሐዋርያት ቅዳሴ" \
>  -section-title "Complete Liturgy and Anaphora of the Apostles" \
>  -section-title-am "ሙሉ ሥርዓተ ቅዳሴና የሐዋርያት ቅዳሴ"
>@go run ./cmd/convertanaphora \
>  -file "$(OUR_LORD_SOURCE)" \
>  -out "$(OUR_LORD_IMPORT)" \
>  -slug "our-lord-jesus-christ" \
>  -name "Anaphora of Our Lord Jesus Christ" \
>  -name-am "የጌታችን የኢየሱስ ክርስቶስ ቅዳሴ" \
>  -section-title "Complete Anaphora of Our Lord Jesus Christ" \
>  -section-title-am "ሙሉ የጌታችን የኢየሱስ ክርስቶስ ቅዳሴ"
>@go run ./cmd/convertanaphora \
>  -file "$(ST_ATHANASIUS_SOURCE)" \
>  -out "$(ST_ATHANASIUS_IMPORT)" \
>  -slug "st-athanasius" \
>  -name "Anaphora of St. Athanasius" \
>  -name-am "የቅዱስ አትናቴዎስ ቅዳሴ" \
>  -section-title "Complete Liturgy and Anaphora of St. Athanasius" \
>  -section-title-am "ሙሉ ሥርዓተ ቅዳሴና የቅዱስ አትናቴዎስ ቅዳሴ"

import-anaphoras: prepare-anaphoras
>@go run ./cmd/importcontent -file "$(APOSTLES_IMPORT)"
>@go run ./cmd/importcontent -file "$(OUR_LORD_IMPORT)"
>@go run ./cmd/importcontent -file "$(ST_ATHANASIUS_IMPORT)"

publish-anaphoras:
>@go run ./cmd/publishcontent \
>  -slug "apostles" \
>  -confirm "apostles" \
>  -allow-review-required
>@go run ./cmd/publishcontent \
>  -slug "our-lord-jesus-christ" \
>  -confirm "our-lord-jesus-christ" \
>  -allow-review-required
>@go run ./cmd/publishcontent \
>  -slug "st-athanasius" \
>  -confirm "st-athanasius" \
>  -allow-review-required

integrate-anaphoras: import-anaphoras publish-anaphoras

ST_EPIPHANIUS_SOURCE ?= content/generated/Anaphora-of-St-Epiphanius.json
ST_EPIPHANIUS_IMPORT ?= content/generated/st-epiphanius-import.json
ST_JOHN_CHRYSOSTOM_SOURCE ?= content/generated/Anaphora-of-St-John-Chrysostom.json
ST_JOHN_CHRYSOSTOM_IMPORT ?= content/generated/st-john-chrysostom-import.json
ST_DIOSCORUS_SOURCE ?= content/generated/Anaphora-of-St-Dioscorus.json
ST_DIOSCORUS_IMPORT ?= content/generated/st-dioscorus-import.json
ST_DIOSCORUS_SEASONAL_SOURCE ?= content/generated/Saint-Dioscorous-Fasika-pentcost.json
ST_DIOSCORUS_SEASONAL_IMPORT ?= content/generated/st-dioscorus-fasika-pentecost-import.json
ST_GREGORY_SOURCE ?= content/generated/Anaphora-of-St-Gregory.json
ST_GREGORY_IMPORT ?= content/generated/st-gregory-import.json

prepare-additional-anaphoras:
>@go run ./cmd/convertanaphora \
>  -file "$(ST_EPIPHANIUS_SOURCE)" \
>  -out "$(ST_EPIPHANIUS_IMPORT)" \
>  -slug "st-epiphanius" \
>  -name "Anaphora of St. Epiphanius" \
>  -name-am "የቅዱስ ኤጲፋንዮስ ቅዳሴ" \
>  -section-title "Complete Liturgy and Anaphora of St. Epiphanius" \
>  -section-title-am "ሙሉ ሥርዓተ ቅዳሴና የቅዱስ ኤጲፋንዮስ ቅዳሴ"
>@go run ./cmd/convertanaphora \
>  -file "$(ST_JOHN_CHRYSOSTOM_SOURCE)" \
>  -out "$(ST_JOHN_CHRYSOSTOM_IMPORT)" \
>  -slug "st-john-chrysostom" \
>  -name "Anaphora of St. John Chrysostom" \
>  -name-am "የቅዱስ ዮሐንስ አፈወርቅ ቅዳሴ" \
>  -section-title "Complete Liturgy and Anaphora of St. John Chrysostom" \
>  -section-title-am "ሙሉ ሥርዓተ ቅዳሴና የቅዱስ ዮሐንስ አፈወርቅ ቅዳሴ"
>@go run ./cmd/convertanaphora \
>  -file "$(ST_DIOSCORUS_SOURCE)" \
>  -out "$(ST_DIOSCORUS_IMPORT)" \
>  -slug "st-dioscorus" \
>  -name "Anaphora of St. Dioscorus" \
>  -name-am "የቅዱስ ዲዮስቆሮስ ቅዳሴ" \
>  -section-title "Complete Liturgy and Anaphora of St. Dioscorus" \
>  -section-title-am "ሙሉ ሥርዓተ ቅዳሴና የቅዱስ ዲዮስቆሮስ ቅዳሴ"
>@go run ./cmd/convertanaphora \
>  -file "$(ST_DIOSCORUS_SEASONAL_SOURCE)" \
>  -out "$(ST_DIOSCORUS_SEASONAL_IMPORT)" \
>  -slug "st-dioscorus-fasika-pentecost" \
>  -name "St. Dioscorus for Fasika and Pentecost" \
>  -name-am "የቅዱስ ዲዮስቆሮስ የፋሲካና የጰራቅሊጦስ ቅዳሴ" \
>  -section-title "Complete St. Dioscorus Liturgy for Fasika and Pentecost" \
>  -section-title-am "ሙሉ የቅዱስ ዲዮስቆሮስ የፋሲካና የጰራቅሊጦስ ቅዳሴ"
>@go run ./cmd/convertanaphora \
>  -file "$(ST_GREGORY_SOURCE)" \
>  -out "$(ST_GREGORY_IMPORT)" \
>  -slug "st-gregory" \
>  -name "Anaphora of St. Gregory" \
>  -name-am "የቅዱስ ጎርጎርዮስ ቅዳሴ" \
>  -section-title "Complete Liturgy and Anaphora of St. Gregory" \
>  -section-title-am "ሙሉ ሥርዓተ ቅዳሴና የቅዱስ ጎርጎርዮስ ቅዳሴ"

import-additional-anaphoras: prepare-additional-anaphoras
>@go run ./cmd/importcontent -file "$(ST_EPIPHANIUS_IMPORT)"
>@go run ./cmd/importcontent -file "$(ST_JOHN_CHRYSOSTOM_IMPORT)"
>@go run ./cmd/importcontent -file "$(ST_DIOSCORUS_IMPORT)"
>@go run ./cmd/importcontent -file "$(ST_DIOSCORUS_SEASONAL_IMPORT)"
>@go run ./cmd/importcontent -file "$(ST_GREGORY_IMPORT)"

publish-additional-anaphoras:
>@go run ./cmd/publishcontent -slug "st-epiphanius" -confirm "st-epiphanius" -allow-review-required
>@go run ./cmd/publishcontent -slug "st-john-chrysostom" -confirm "st-john-chrysostom" -allow-review-required
>@go run ./cmd/publishcontent -slug "st-dioscorus" -confirm "st-dioscorus" -allow-review-required
>@go run ./cmd/publishcontent -slug "st-dioscorus-fasika-pentecost" -confirm "st-dioscorus-fasika-pentecost" -allow-review-required
>@go run ./cmd/publishcontent -slug "st-gregory" -confirm "st-gregory" -allow-review-required

integrate-additional-anaphoras: import-additional-anaphoras publish-additional-anaphoras

LITURGY_GUIDE_SOURCE ?= content/generated/liturgy.json
LITURGY_GUIDE_IMPORT ?= content/generated/liturgy-guide-import.json

prepare-liturgy-guide:
>@go run ./cmd/convertguide \
>  -file "$(LITURGY_GUIDE_SOURCE)" \
>  -out "$(LITURGY_GUIDE_IMPORT)" \
>  -slug "liturgy-guide" \
>  -name "Introduction to the Divine Liturgy" \
>  -name-am "የሥርዓተ ቅዳሴ መግቢያ" \
>  -section-title "Understanding the Divine Liturgy" \
>  -section-title-am "ሥርዓተ ቅዳሴን መረዳት"

import-liturgy-guide: prepare-liturgy-guide
>@go run ./cmd/importcontent -file "$(LITURGY_GUIDE_IMPORT)"

publish-liturgy-guide:
>@go run ./cmd/publishcontent \
>  -slug "liturgy-guide" \
>  -confirm "liturgy-guide" \
>  -allow-review-required

integrate-liturgy-guide: import-liturgy-guide publish-liturgy-guide

integrate-additional-content: integrate-additional-anaphoras integrate-liturgy-guide

CATALOG_FILES := \
	content/catalog/st-basil.json \
	content/catalog/three-hundred.json \
	content/catalog/st-cyril.json \
	content/catalog/st-jacob-of-serough.json

import-catalog:
>@for file in $(CATALOG_FILES); do \
>	go run ./cmd/importcontent -file "$$file" || exit 1; \
>done

publish-catalog:
>@for slug in st-basil three-hundred st-cyril st-jacob-of-serough; do \
>	go run ./cmd/publishcontent -slug "$$slug" -confirm "$$slug" || exit 1; \
>done

integrate-catalog: import-catalog publish-catalog

# The fourteen liturgies now come from content/updated. The targets above that
# build the same slugs out of content/generated are kept for reference but are
# deliberately no longer part of these aggregates: running both would overwrite
# the newer content with the older OCR drafts.
#
# import-catalog is excluded for the same reason. Its four stubs carry empty
# section lists, and importing a document replaces a liturgy's sections
# wholesale, so it would now erase the real content that st-basil, st-cyril,
# st-jacob-of-serough and three-hundred just gained.

.PHONY: prepare-dioscorus-seasonal import-dioscorus-seasonal
.PHONY: publish-dioscorus-seasonal integrate-dioscorus-seasonal

# Saint-Dioscorous-Fasika-pentcost.json has no content/updated counterpart, so
# this one liturgy still runs through the older flat-anaphora pipeline.
prepare-dioscorus-seasonal:
>@go run ./cmd/convertanaphora \
>  -file "$(ST_DIOSCORUS_SEASONAL_SOURCE)" \
>  -out "$(ST_DIOSCORUS_SEASONAL_IMPORT)" \
>  -slug "st-dioscorus-fasika-pentecost" \
>  -name "St. Dioscorus for Fasika and Pentecost" \
>  -name-am "የቅዱስ ዲዮስቆሮስ የፋሲካና የጰራቅሊጦስ ቅዳሴ" \
>  -section-title "Complete St. Dioscorus Liturgy for Fasika and Pentecost" \
>  -section-title-am "ሙሉ የቅዱስ ዲዮስቆሮስ የፋሲካና የጰራቅሊጦስ ቅዳሴ"

import-dioscorus-seasonal: prepare-dioscorus-seasonal
>@go run ./cmd/importcontent -file "$(ST_DIOSCORUS_SEASONAL_IMPORT)"

publish-dioscorus-seasonal:
>@go run ./cmd/publishcontent \
>  -slug "st-dioscorus-fasika-pentecost" \
>  -confirm "st-dioscorus-fasika-pentecost" \
>  -allow-review-required

integrate-dioscorus-seasonal: import-dioscorus-seasonal publish-dioscorus-seasonal

prepare-all-content: prepare-updated-liturgies prepare-dioscorus-seasonal prepare-liturgy-guide

import-all-content: import-updated-liturgies import-dioscorus-seasonal import-liturgy-guide

publish-all-content: publish-updated-liturgies publish-dioscorus-seasonal publish-liturgy-guide

set-audio:
>@test -n "$(slug)" || (echo "usage: make set-audio slug=example file=recording.mp3 url=https://... duration_ms=123 confirm=example" && exit 1)
>@test -n "$(file)" || (echo "file is required" && exit 1)
>@test -n "$(url)" || (echo "url is required" && exit 1)
>@test -n "$(duration_ms)" || (echo "duration_ms is required" && exit 1)
>@go run ./cmd/setaudio \
>  -slug "$(slug)" \
>  -file "$(file)" \
>  -url "$(url)" \
>  -duration-ms "$(duration_ms)" \
>  -confirm "$(confirm)"

export-offline:
>@test -n "$(out)" || (echo "usage: make export-offline out=../mobile/assets/offline" && exit 1)
>@go run ./cmd/exportoffline -out "$(out)"

# --- Normalized liturgy sources (content/updated) -------------------------
#
# The anaphoras under content/updated no longer carry the shared opening.
# Qidase_serate.json holds it once and every conversion joins it in front of
# the anaphora, renumbering the pages continuously across the join. Pages are
# packed to a rune budget so a service runs to roughly half as many pages as
# its source groups without any one page turning into a long scroll.

UPDATED_DIR ?= content/updated
UPDATED_BEGINNING ?= $(UPDATED_DIR)/Qidase_serate.json
UPDATED_MANIFEST ?= content/updated-liturgies.tsv
UPDATED_IMPORT_DIR ?= content/generated/updated
UPDATED_TARGET_RUNES ?= 4000
UPDATED_SLUGS := $(shell awk -F'\t' 'NF > 1 && $$0 !~ /^#/ { print $$1 }' content/updated-liturgies.tsv)

.PHONY: prepare-updated-liturgies import-updated-liturgies
.PHONY: publish-updated-liturgies integrate-updated-liturgies

prepare-updated-liturgies:
>@go run ./cmd/convertliturgy \
>  -beginning "$(UPDATED_BEGINNING)" \
>  -manifest "$(UPDATED_MANIFEST)" \
>  -source-dir "$(UPDATED_DIR)" \
>  -out-dir "$(UPDATED_IMPORT_DIR)" \
>  -target-runes "$(UPDATED_TARGET_RUNES)"

import-updated-liturgies: prepare-updated-liturgies
>@for slug in $(UPDATED_SLUGS); do \
>	go run ./cmd/importcontent \
>	  -file "$(UPDATED_IMPORT_DIR)/$$slug-import.json" || exit 1; \
>done

# Five entries in the shared beginning are still flagged needs_review, so this
# publishes as an explicit preview until they are cleared.
publish-updated-liturgies:
>@for slug in $(UPDATED_SLUGS); do \
>	go run ./cmd/publishcontent \
>	  -slug "$$slug" \
>	  -confirm "$$slug" \
>	  -allow-review-required || exit 1; \
>done

integrate-updated-liturgies: import-updated-liturgies publish-updated-liturgies
