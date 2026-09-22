APP_NAME := developers-cv-api
GO := go

-include .env
export

MIGRATE_IMAGE ?= migrate/migrate:v4.18.3
MIGRATIONS_DIR := db/migrations

POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432
POSTGRES_DB ?= developers_cv
POSTGRES_USER ?= developers
POSTGRES_SSLMODE ?= disable

MIGRATE_DATABASE_URL := postgres://$(POSTGRES_USER)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)

.PHONY: help fmt vet test test-unit build run lint generate docker-build compose-up compose-down check migrate-up migrate-down migrate-version

help:
	@printf '%s\n' \
		'make fmt             - format Go source files' \
		'make vet             - run go vet' \
		'make test            - run all Go tests' \
		'make test-unit       - run unit tests' \
		'make build           - build API binary' \
		'make run             - run API locally' \
		'make lint            - run golangci-lint' \
		'make generate        - regenerate Wire code' \
		'make docker-build    - build backend image' \
		'make compose-up      - start local infrastructure' \
		'make compose-down    - stop local infrastructure' \
		'make migrate-up      - apply database migrations' \
		'make migrate-down    - rollback last database migration' \
		'make migrate-version - show current database migration version' \
		'make check           - run format, vet, test and lint'

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

test-unit:
	$(GO) test ./internal/...

build:
	$(GO) build ./cmd/api

run:
	$(GO) run ./cmd/api

lint:
	golangci-lint run ./...

generate:
	$(GO) run github.com/google/wire/cmd/wire ./...

docker-build:
	docker build -t $(APP_NAME):local .

compose-up:
	docker compose up -d --build

compose-down:
	docker compose down

check: fmt vet test lint

migrate-up:
	@test -n "$(POSTGRES_PASSWORD)" || (printf '%s\n' 'POSTGRES_PASSWORD is required. Set it in .env.' >&2; exit 1)
	PGPASSWORD="$(POSTGRES_PASSWORD)" docker run --rm \
		--network host \
		-v "$(PWD)/$(MIGRATIONS_DIR):/migrations:ro" \
		$(MIGRATE_IMAGE) \
		-path=/migrations \
		-database "$(MIGRATE_DATABASE_URL)" \
		up

migrate-down:
	@test -n "$(POSTGRES_PASSWORD)" || (printf '%s\n' 'POSTGRES_PASSWORD is required. Set it in .env.' >&2; exit 1)
	PGPASSWORD="$(POSTGRES_PASSWORD)" docker run --rm \
		--network host \
		-v "$(PWD)/$(MIGRATIONS_DIR):/migrations:ro" \
		$(MIGRATE_IMAGE) \
		-path=/migrations \
		-database "$(MIGRATE_DATABASE_URL)" \
		down 1

migrate-version:
	@test -n "$(POSTGRES_PASSWORD)" || (printf '%s\n' 'POSTGRES_PASSWORD is required. Set it in .env.' >&2; exit 1)
	PGPASSWORD="$(POSTGRES_PASSWORD)" docker run --rm \
		--network host \
		-v "$(PWD)/$(MIGRATIONS_DIR):/migrations:ro" \
		$(MIGRATE_IMAGE) \
		-path=/migrations \
		-database "$(MIGRATE_DATABASE_URL)" \
		version