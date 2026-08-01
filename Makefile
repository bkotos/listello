.PHONY: test run bdd-report bdd-steps

REPORTS_DIR := reports
CUCUMBER_JSON := $(CURDIR)/$(REPORTS_DIR)/cucumber.json

test:
	go test ./...

run:
	go run ./cmd/listello

# Run domain Gherkin suite with pretty output + Cucumber JSON report.
# Undefined/pending/ambiguous steps fail the run (-godog.strict).
bdd-report:
	mkdir -p $(REPORTS_DIR)
	go test ./internal/listello-domain -v -count=1 \
		-godog.format=pretty,cucumber:$(CUCUMBER_JSON) \
		-godog.strict
	@echo "Cucumber JSON written to $(CUCUMBER_JSON)"

# List registered step definitions (implemented steps).
bdd-steps:
	go test ./internal/listello-domain -v -count=1 -godog.definitions
