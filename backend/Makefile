APP_NAME := developers-cv-api
GO := go
GOLANGCI_LINT_VERSION := 2.13.2


.PHONY: help fmt fmt-check verify-deps vet test test-unit test-devauth check-integration build run run-mock-auth check-lint-version lint generate generate-check docker-build compose-up compose-down check release-check migrate-up migrate-down migrate-version seed-dev

help:
	@printf '%s\n' \
		'make fmt             - format Go source files' \
		'make fmt-check       - verify Go formatting without modifying files' \
		'make verify-deps     - verify downloaded Go modules' \
		'make vet             - run go vet' \
		'make test            - run all Go tests without cache' \
		'make test-unit       - run internal unit tests without cache' \
		'make check-integration - run integration tests (requires Docker/Testcontainers)' \
		'make build           - build API binary' \
		'make run             - run API locally (without mock auth)' \
		'make run-mock-auth   - run API locally with mock auth (-tags=devauth)' \
		'make test-devauth    - run tests with mock auth compiled in' \
		'make check-lint-version - verify the pinned golangci-lint version' \
		'make lint            - run golangci-lint with the pinned version' \
		'make generate        - regenerate Wire code' \
		'make generate-check  - verify Wire output is reproducible without changing files' \
		'make docker-build    - build backend image' \
		'make compose-up      - start local infrastructure' \
		'make compose-down    - stop local infrastructure' \
		'make migrate-up      - apply database migrations' \
		'make migrate-down    - rollback last database migration' \
		'make migrate-version - show current database migration version' \
		'make seed-dev        - insert optional local demo users/profiles' \
		'make check           - run non-mutating local quality checks' \
		'make release-check   - run quality, integration, and container build gates'

fmt:
	$(GO) fmt ./...

fmt-check:
	@files=$$(gofmt -l $$(find . -type f -name '*.go' -not -path './vendor/*')); \
	if [ -n "$$files" ]; then \
		printf 'Go files need formatting:\n%s\n' "$$files"; \
		exit 1; \
	fi

verify-deps:
	$(GO) mod verify

vet:
	$(GO) vet ./...

test:
	$(GO) test -count=1 ./...

test-unit:
	$(GO) test -count=1 ./internal/...

check-integration:
	$(GO) test -tags integration -count=1 ./...

build:
	$(GO) build ./cmd/api

run:
	$(GO) run ./cmd/api

run-mock-auth:
	$(GO) run -tags=devauth ./cmd/api

test-devauth:
	$(GO) test -tags=devauth -count=1 ./...

check-lint-version:
	@command -v golangci-lint >/dev/null 2>&1 || { echo 'golangci-lint is not installed; expected version $(GOLANGCI_LINT_VERSION)' >&2; exit 1; }
	@golangci-lint version | grep -Eq "(^|[[:space:]])v?$(GOLANGCI_LINT_VERSION)([[:space:]]|$$)" || { echo 'golangci-lint version mismatch; CI uses $(GOLANGCI_LINT_VERSION)' >&2; golangci-lint version; exit 1; }

lint: check-lint-version
	golangci-lint run ./...

generate:
	$(GO) tool github.com/google/wire/cmd/wire ./internal/app

generate-check:
	@set -eu; \
	backup=$$(mktemp); \
	cp internal/app/wire_gen.go "$$backup"; \
	restore() { cp "$$backup" internal/app/wire_gen.go; rm -f "$$backup"; }; \
	trap restore EXIT HUP INT TERM; \
	if ! $(GO) tool github.com/google/wire/cmd/wire ./internal/app; then exit 1; fi; \
	if ! cmp -s "$$backup" internal/app/wire_gen.go; then \
		echo 'Wire output is stale: internal/app/wire_gen.go differs after generation.'; \
		diff -u "$$backup" internal/app/wire_gen.go || true; \
		exit 1; \
	fi

docker-build:
	docker build -t $(APP_NAME):local .

compose-up:
	docker compose up -d --build

compose-down:
	docker compose down

check:
	$(MAKE) fmt-check
	$(MAKE) verify-deps
	$(MAKE) generate-check
	$(MAKE) vet
	$(MAKE) test
	$(MAKE) lint

release-check:
	$(MAKE) check
	$(MAKE) check-integration
	$(MAKE) docker-build

migrate-up:
	docker compose run --rm migrations

migrate-down:
	docker compose run --rm \
		--entrypoint /bin/sh \
		migrations \
		-c \
			'exec migrate -path=/migrations -database "$$MIGRATE_DATABASE_URL" down 1'

migrate-version:
	docker compose run --rm \
		--entrypoint /bin/sh \
		migrations \
		-c \
		'exec migrate -path=/migrations -database "$$MIGRATE_DATABASE_URL" version'

seed-dev:
	docker compose exec -T postgres sh -c 'psql -v ON_ERROR_STOP=1 -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"' < db/seeds/development.sql
