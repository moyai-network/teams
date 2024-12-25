.PHONY: mocks

lint:
	./tests/scripts/lint.sh
mocks:
	./tests/scripts/mocks.sh
tests: lint mocks
	./tests/scripts/tests.sh
