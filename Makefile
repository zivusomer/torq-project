.PHONY: build run debug test lint

DLV := $(shell go env GOPATH)/bin/dlv

run:
	go run ./cmd/torq-project

debug:
	$(DLV) debug ./cmd/torq-project --headless --listen=:2345 --api-version=2 --accept-multiclient --continue

build:
	go build ./cmd/torq-project

test:
	go test ./...

lint:
	./scripts/lint-custom.sh
	go fmt ./...
	go mod tidy
	go vet ./...
	go test ./...
