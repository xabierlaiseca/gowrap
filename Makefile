.DEFAULT_GOAL := ci
WRAPPED_COMMANDS := go gofmt
VERSION ?= "0.0.0"

BIN_DIR := ./bin
DIST_DIR := ./dist

.PHONY: bin build build-archive build-init ci clean cmd-wrappers generate-versions-file gowrap-cmd lint fmt test

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
	rm -rf $(DIST_DIR)

test:
	go test -v ./...

generate-versions-file: gowrap-cmd
	$(BIN_DIR)/gowrap versions-file generate --file data/versions.json

build-archive: bin
	mkdir -p $(DIST_DIR)
	tar -czvf $(DIST_DIR)/gowrap-$(VERSION)-$$GOOS-$$GOARCH.tar.gz -C $(BIN_DIR) .
