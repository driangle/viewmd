.PHONY: build test lint clean

BINARY := viewmd
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/viewmd

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR)
