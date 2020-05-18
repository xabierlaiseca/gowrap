.DEFAULT_GOAL := ci
WRAPPED_COMMANDS := go gofmt

BIN_DIR := ./bin

.PHONY: bin build build-init ci clean cmd-wrappers generate-versions-file gowrap-cmd lint fmt test test-integration

lint:
	golangci-lint run

fmt:
	go fmt ./...

build-init:
	mkdir -p $(BIN_DIR)

gowrap-cmd: build-init
	go build -o $(BIN_DIR)/gowrap cmd/gowrap/main.go

cmd-wrappers: build-init
	for cmd in $(WRAPPED_COMMANDS); do \
		go build -ldflags="-X 'main.wrappedCmd=$$cmd'" -o $(BIN_DIR)/$$cmd cmd/generic-cmd-wrapper/main.go; \
	done

bin: gowrap-cmd cmd-wrappers

ci: lint test
build: bin fmt ci

clean:
	rm -rf $(BIN_DIR)

test:
	go test -v ./...

test-integration:
	go test -v --tags=integration pkg/integration/*_test.go

generate-versions-file: gowrap-cmd
	$(BIN_DIR)/gowrap versions-file generate --file data/versions.json
