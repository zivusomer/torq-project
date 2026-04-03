.PHONY: build run debug test lint ensure-redis ensure-port-8080-free ensure-port-2345-free

DLV := $(shell go env GOPATH)/bin/dlv

run: ensure-redis ensure-port-8080-free
	go run ./cmd/torq-project

debug: ensure-redis ensure-port-8080-free ensure-port-2345-free
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

ensure-port-8080-free:
	@PIDS="$$(lsof -ti tcp:8080 -sTCP:LISTEN 2>/dev/null || true)"; \
	if [ -n "$$PIDS" ]; then \
		echo "Port 8080 is busy. Stopping listener(s): $$PIDS"; \
		kill $$PIDS || true; \
		sleep 1; \
	fi

ensure-port-2345-free:
	@PIDS="$$(lsof -ti tcp:2345 -sTCP:LISTEN 2>/dev/null || true)"; \
	if [ -n "$$PIDS" ]; then \
		echo "Port 2345 is busy. Stopping listener(s): $$PIDS"; \
		kill $$PIDS || true; \
		sleep 1; \
	fi
