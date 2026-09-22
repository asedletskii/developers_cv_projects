APP_NAME := developers-cv-api
GO := go

.PHONY: help fmt vet test test-unit build run lint generate docker-build compose-up compose-down check

help:
	@printf '%s\n' \
		'make fmt          - format Go source files' \
		'make vet          - run go vet' \
		'make test         - run all Go tests' \
		'make test-unit    - run unit tests' \
		'make build        - build API binary' \
		'make run          - run API locally' \
		'make lint         - run golangci-lint' \
		'make generate     - regenerate Wire code' \
		'make docker-build - build backend image' \
		'make compose-up   - start local infrastructure' \
		'make compose-down - stop local infrastructure' \
		'make check        - run format, vet, test and lint'

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
