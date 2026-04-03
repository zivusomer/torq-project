.PHONY: build run debug test lint ensure-redis

DLV := $(shell go env GOPATH)/bin/dlv

run: ensure-redis
	go run ./cmd/torq-project

debug: ensure-redis
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

ensure-redis:
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "docker is required to auto-start redis for run/debug"; \
		exit 1; \
	fi
	@if docker inspect torq-redis >/dev/null 2>&1; then \
		if ! docker ps --format '{{.Names}}' | rg -x 'torq-redis' >/dev/null; then \
			echo "Starting redis container torq-redis..."; \
			if ! docker start torq-redis >/dev/null; then \
				echo "Cleaning failed torq-redis container and recreating..."; \
				docker rm -f torq-redis >/dev/null 2>&1 || true; \
				docker run --name torq-redis -p 6379:6379 -d redis:7-alpine >/dev/null; \
			fi; \
		fi; \
	else \
		echo "Creating redis container torq-redis..."; \
		if ! docker run --name torq-redis -p 6379:6379 -d redis:7-alpine >/dev/null; then \
			echo "Cleaning conflicting torq-redis container and retrying..."; \
			docker rm -f torq-redis >/dev/null 2>&1 || true; \
			docker run --name torq-redis -p 6379:6379 -d redis:7-alpine >/dev/null; \
		fi; \
	fi
