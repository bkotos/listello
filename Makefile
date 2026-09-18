.DEFAULT_GOAL := help

COMPOSE_POSTGRES := docker compose -f docker-compose.postgres.yml
LISTELLO_TEST_POSTGRES_DSN ?= postgres://listello:listello@127.0.0.1:5433/listello?sslmode=disable

.PHONY: help run run-api run-ui test-api test-ui test-e2e-cli test test-api-postgres postgres-up postgres-down api-types build

help:
	@echo "Available targets:"
	@echo "  make help              Show this list of targets"
	@echo "  make run               Run API and UI together"
	@echo "  make run-api           Run the API server with auto-reload"
	@echo "  make run-ui            Run the UI dev server"
	@echo "  make test              Run API, UI, and CLI e2e tests"
	@echo "  make test-api          Run API tests"
	@echo "  make test-api-postgres Run API tests against Compose PostgreSQL"
	@echo "  make postgres-up       Start Compose PostgreSQL for tests"
	@echo "  make postgres-down     Stop Compose PostgreSQL and remove volumes"
	@echo "  make test-ui           Run UI tests"
	@echo "  make test-e2e-cli      Run CLI e2e tests"
	@echo "  make build             Build API CLI, API server, and UI"
	@echo "  make api-types         Generate TypeScript types from the API"

run:
	@trap 'kill 0' EXIT INT TERM; \
	$(MAKE) -C api serve-watch & \
	npm run dev -w ui & \
	wait

run-api:
	$(MAKE) -C api serve-watch

run-ui:
	npm run dev -w ui

test-api:
	$(MAKE) -C api test

postgres-up:
	@command -v docker >/dev/null || { echo "docker not found. Install Docker to test against PostgreSQL."; exit 1; }
	$(COMPOSE_POSTGRES) up -d --wait

postgres-down:
	$(COMPOSE_POSTGRES) down --volumes --remove-orphans

test-api-postgres: postgres-up
	LISTELLO_TEST_POSTGRES_DSN="$(LISTELLO_TEST_POSTGRES_DSN)" $(MAKE) -C api test

test-ui:
	npm test -w ui

test-e2e-cli:
	npm test -w e2e-cli

test: test-api test-ui test-e2e-cli

build:
	$(MAKE) -C api build-cli build-server
	npm run build -w ui

api-types:
	$(MAKE) -C api api-types
	