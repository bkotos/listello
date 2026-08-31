.DEFAULT_GOAL := list

.PHONY: list test mocks clean build run build-wasm run-wasm serve bdd-report bdd-gaps bdd-steps

REPORTS_DIR := reports
CUCUMBER_JSON := $(CURDIR)/$(REPORTS_DIR)/cucumber.json
MOCKERY ?= $(shell go env GOPATH)/bin/mockery
BIN := bin/listello
BIN_WASM := bin/listello.wasm
WASMTIME ?= wasmtime

list:
	@echo "Available targets:"
	@echo "  make list          Show this list of targets"
	@echo "  make test          Run all Go tests"
	@echo "  make mocks         Regenerate testify mocks (requires mockery v3)"
	@echo "  make clean         Remove built binaries"
	@echo "  make build         Clean and compile CLI to $(BIN)"
	@echo "  make run           Run CLI (pass args via ARGS=...)"
	@echo "  make serve         Run HTTP API server (pass args via ARGS=...)"
	@echo "  make build-wasm    Clean and compile CLI to WASI WASM ($(BIN_WASM))"
	@echo "  make run-wasm      Run WASM CLI via wasmtime (pass args via ARGS=...)"
	@echo "  make bdd-report    Run implemented BDD features + write Cucumber JSON"
	@echo "  make bdd-gaps      Run all BDD features including @wip (shows undefined steps)"
	@echo "  make bdd-steps     List registered BDD step definitions"

test:
	go test ./...

mocks:
	@test -x "$(MOCKERY)" || { echo "mockery not found. Install with: go install github.com/vektra/mockery/v3@latest"; exit 1; }
	"$(MOCKERY)"

clean:
	rm -rf bin

build: clean
	mkdir -p bin
	go build -o $(BIN) ./cmd/listello

run: build
	./$(BIN) $(ARGS)

serve: build
	./$(BIN) serve $(ARGS)

build-wasm: clean
	mkdir -p bin
	GOOS=wasip1 GOARCH=wasm go build -tags sqlite3_dotlk -o $(BIN_WASM) ./cmd/listello

run-wasm: build-wasm
	@command -v $(WASMTIME) >/dev/null || { echo "wasmtime not found. Install with: brew install wasmtime"; exit 1; }
	$(WASMTIME) run --dir=.::/ $(BIN_WASM) $(ARGS)

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
