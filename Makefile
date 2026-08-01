.DEFAULT_GOAL := list

.PHONY: list test run bdd-report bdd-gaps bdd-steps

REPORTS_DIR := reports
CUCUMBER_JSON := $(CURDIR)/$(REPORTS_DIR)/cucumber.json

list:
	@echo "Available targets:"
	@echo "  make list         Show this list of targets"
	@echo "  make test         Run all Go tests"
	@echo "  make run          Create a list via cmd/listello (optional: NAME=\"Someday\")"
	@echo "  make bdd-report   Run implemented BDD features + write Cucumber JSON"
	@echo "  make bdd-gaps     Run all BDD features including @wip (shows undefined steps)"
	@echo "  make bdd-steps    List registered BDD step definitions"

test:
	go test ./...

run:
	go run ./cmd/listello $(NAME)

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
