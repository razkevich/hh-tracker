.PHONY: watch
watch:
	@./bin/watch

.PHONY: start
start:
	@./bin/start

.PHONY: stop
stop:
	@./bin/stop

.PHONY: component-test
component-test:
	@./bin/component-test

.PHONY: lint
lint:
	@./bin/lint

.PHONY: test-coverage
test-coverage:
	@./bin/test-coverage

.PHONY: test
test:
	@./bin/test

.PHONY: destroy
destroy:
	@./bin/destroy

.PHONY: prepare
prepare:
	@./bin/prepare

.PHONY: check
check:
	@./bin/prepare
	@./bin/test
	@./bin/component-test
