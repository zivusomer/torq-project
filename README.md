# torq-project

Starter Go project scaffold.

## Prerequisites

- Go 1.22+ installed

## Project layout

- `cmd/torq-project/`: application entrypoint
- `internal/config/`: env-based configuration
- `internal/logging/`: structured logger setup

## Environment variables

- `APP_NAME` (default: `torq-project`)
- `APP_ENV` (default: `development`)
- `LOG_LEVEL` (default: `info`)

## Quick start

```bash
make run
```

## Commands

```bash
make build
make run
make test
make lint
```

`make lint` runs project hygiene checks: Go formatting, module dependency cleanup (`go mod tidy`), static analysis (`go vet`), and tests.
