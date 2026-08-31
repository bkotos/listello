.PHONY: run-api run-ui test-api test-ui

run-api:
	$(MAKE) -C api serve

run-ui:
	npm run dev -w ui

test-api:
	$(MAKE) -C api test

test-ui:
	npm test -w ui

test: test-api test-ui
