.PHONY: mocks

lint:
	./tests/scripts/lint.sh
mocks:
	./tests/scripts/mocks.sh
cloc:
	./tests/scripts/cloc.sh
tests: mocks lint
	./tests/scripts/tests.sh