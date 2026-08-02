.DEFAULT_GOAL := list

.PHONY: list test mocks build run bdd-report bdd-gaps bdd-steps

REPORTS_DIR := reports
CUCUMBER_JSON := $(CURDIR)/$(REPORTS_DIR)/cucumber.json
MOCKERY ?= $(shell go env GOPATH)/bin/mockery
BIN := bin/listello

list:
	@echo "Available targets:"
	@echo "  make list         Show this list of targets"
	@echo "  make test         Run all Go tests"
	@echo "  make mocks        Regenerate testify mocks (requires mockery v3)"
	@echo "  make build        Compile CLI to $(BIN)"
	@echo "  make run          Run CLI (pass args via ARGS=...)"
	@echo "  make bdd-report   Run implemented BDD features + write Cucumber JSON"
	@echo "  make bdd-gaps     Run all BDD features including @wip (shows undefined steps)"
	@echo "  make bdd-steps    List registered BDD step definitions"

test:
	go test ./...

mocks:
	@test -x "$(MOCKERY)" || { echo "mockery not found. Install with: go install github.com/vektra/mockery/v3@latest"; exit 1; }
	"$(MOCKERY)"

build:
	mkdir -p bin
	go build -o $(BIN) ./cmd/listello

run: build
	./$(BIN) $(ARGS)

# Run implemented domain Gherkin suite (excludes @wip) + Cucumber JSON report.
bdd-report:
	mkdir -p $(REPORTS_DIR)
	go test ./internal/listello-domain -v -count=1 \
		-godog.format=pretty,cucumber:$(CUCUMBER_JSON) \
		-godog.strict
	@echo "Cucumber JSON written to $(CUCUMBER_JSON)"

# Include @wip scenarios so undefined steps show up in the report.
bdd-gaps:
	mkdir -p $(REPORTS_DIR)
	go test ./internal/listello-domain -v -count=1 \
		-godog.tags= \
		-godog.format=pretty,cucumber:$(CUCUMBER_JSON)
	@echo "Cucumber JSON written to $(CUCUMBER_JSON)"

# List registered step definitions (implemented steps).
bdd-steps:
	go test ./internal/listello-domain -v -count=1 -godog.definitions
