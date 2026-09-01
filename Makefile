.PHONY: run run-api run-ui test-api test-ui test api-types

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

test: test-api test-ui

api-types:
	$(MAKE) -C api api-types
