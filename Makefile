.DEFAULT_GOAL := build

BIN_DIR := ./bin
.PHONY: build build-gowrap build-init clean test

pre-build:
	go fmt ./...
	go vet ./...
	golangci-lint run

build-init:
	mkdir -p $(BIN_DIR)

build-gowrap: build-init
	go build -o $(BIN_DIR)/gowrap cmd/gowrap/main.go

build-go: build-init
	go build -o $(BIN_DIR)/go cmd/go/main.go

build: pre-build test build-gowrap build-go

clean:
	rm -r $(BIN_DIR)

test:
	go test -v ./...

generate-versions-file: build-gowrap
	$(BIN_DIR)/gowrap versions-file generate --file data/versions.json
