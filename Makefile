.PHONY: build run test lint

run:
	go run ./cmd/torq-project

build:
	go build ./cmd/torq-project

test:
	go test ./...

lint:
	go fmt ./...
	go mod tidy
	go vet ./...
	go test ./...
