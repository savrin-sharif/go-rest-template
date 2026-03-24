APP_NAME ?= go-rest-template
BIN ?= bin/server

.PHONY: build test lint tidy run clean

build:
	go build -o $(BIN) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

clean:
	rm -rf $(BIN)

