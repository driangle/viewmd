.PHONY: build test lint clean setup

build:
	$(MAKE) -C apps/cli build

test:
	$(MAKE) -C apps/cli test

lint:
	$(MAKE) -C apps/cli lint
	pylint viewmd.py tests/

clean:
	$(MAKE) -C apps/cli clean

setup:
	git config core.hooksPath .githooks
	@echo "Git hooks configured."
