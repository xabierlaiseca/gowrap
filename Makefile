BIN_DIR := bin
.PHONY: build build-gowrap build-init clean test

build-init:
	mkdir -p $(BIN_DIR)

build-gowrap: build-init
	go build -o $(BIN_DIR)/gowrap cmd/gowrap/main.go

build: build-gowrap

clean:
	rm -r $(BIN_DIR)

test:
	go test -v ./...
