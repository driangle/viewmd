.PHONY: build test lint clean setup install install-dev install-dev-full

build:
	$(MAKE) -C apps/cli build

test:
	$(MAKE) -C apps/cli test

lint:
	$(MAKE) -C apps/cli lint

install:
	$(MAKE) -C apps/cli install

install-dev:
	$(MAKE) -C apps/cli install-dev

install-dev-full: install-dev
	@# TODO: $(MAKE) -C apps/web install-dev

clean:
	$(MAKE) -C apps/cli clean

setup:
	git config core.hooksPath .githooks
	@echo "Git hooks configured."
