.DEFAULT_GOAL := build-ci

BIN_DIR := ./bin
.PHONY: build build-gowrap build-init clean test

lint:
	golangci-lint run

fmt:
	go fmt ./...

build-init:
	mkdir -p $(BIN_DIR)

build-gowrap: build-init
	go build -o $(BIN_DIR)/gowrap cmd/gowrap/main.go

build-go: build-init
	go build -o $(BIN_DIR)/go cmd/go/main.go

build-ci: lint test build-gowrap build-go
build: fmt build-ci

clean:
	rm -r $(BIN_DIR)

test:
	go test -v ./...

generate-versions-file: build-gowrap
	$(BIN_DIR)/gowrap versions-file generate --file data/versions.json
