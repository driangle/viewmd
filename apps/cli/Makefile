.PHONY: build test lint clean install

BINARY := viewmd
BUILD_DIR := bin
PREFIX ?= $(shell go env GOPATH)

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/viewmd

test:
	go test ./...

lint:
	go vet ./...

install: build
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 $(BUILD_DIR)/$(BINARY) $(DESTDIR)$(PREFIX)/bin/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)
