.DEFAULT_GOAL := help

.PHONY: help run run-api run-ui test-api test-ui test-e2e-cli test api-types build

help:
	@echo "Available targets:"
	@echo "  make help         Show this list of targets"
	@echo "  make run          Run API and UI together"
	@echo "  make run-api      Run the API server with auto-reload"
	@echo "  make run-ui       Run the UI dev server"
	@echo "  make test         Run API, UI, and CLI e2e tests"
	@echo "  make test-api     Run API tests"
	@echo "  make test-ui      Run UI tests"
	@echo "  make test-e2e-cli Run CLI e2e tests"
	@echo "  make build        Build API CLI, API server, and UI"
	@echo "  make api-types    Generate TypeScript types from the API"

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
	