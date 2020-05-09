.DEFAULT_GOAL := ci

BIN_DIR := ./bin
.PHONY: bin build build-init ci clean generate-versions-file go-cmd gowrap-cmd lint fmt test

lint:
	golangci-lint run

fmt:
	go fmt ./...

build-init:
	mkdir -p $(BIN_DIR)

gowrap-cmd: build-init
	go build -o $(BIN_DIR)/gowrap cmd/gowrap/main.go

go-cmd: build-init
	go build -o $(BIN_DIR)/go cmd/go/main.go

bin: gowrap-cmd go-cmd

ci: lint test
build: bin fmt ci

clean:
	rm -rf $(BIN_DIR)

test:
	go test -v ./...

generate-versions-file: gowrap-cmd
	$(BIN_DIR)/gowrap versions-file generate --file data/versions.json
