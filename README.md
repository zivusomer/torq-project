# torq-project

IP-to-country service in Go with a pluggable datastore and built-in rate limiting.

## Prerequisites

Set these up before running project commands:

- Go (recommended: 1.22+)
- Make
- Docker (required by `make run` / `make debug` because they ensure local Redis via container)
- curl (for local API checks)

For debugging:

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

For pre-commit lint execution:

```bash
./scripts/install-git-hooks.sh
```

## Endpoint

- `GET /v1/find-country?ip=2.22.233.255`
- Success response:
  - `{"country":"Israel","city":"Tel Aviv"}`
- Error response format:
  - `{"error":"..."}`

## Code flow

### Startup flow

1. `cmd/torq-project/main.go` calls `config.LoadFromEnv()`.
2. `internal/config` selects preset by `APP_ENV` and applies env var overrides.
3. `internal/logging` initializes the global logger service.
4. `internal/store/factory` selects datastore implementation (currently `csv`).
5. `internal/store/csvstore` loads CSV rows into an in-memory map.
6. `internal/ratelimit` initializes the configured limiter backend (`inmemory` or `redis`).
7. `internal/api/findcountry` creates endpoint handler and `internal/httpserver` registers routes.
8. `http.ListenAndServe` starts serving on configured `PORT`.

### Request flow (`GET /v1/find-country?ip=...`)

1. Route handler validates method (`GET` only).
2. Rate limiter checks allowed requests per second.
3. Handler validates `ip` query parameter and parses it.
4. Store lookup is executed through `store.Resolver`.
5. Response mapping:
   - found -> `200` with `{"country":"...","city":"..."}`
   - missing ip -> `400`
   - invalid ip format -> `400`
   - not found -> `404`
   - rate limited -> `429`
   - unexpected internal error -> `500`

### Sequence overview

```text
Startup:
main -> config.LoadFromEnv -> presetForEnv/env overrides
main -> logging.Logger.Info/Warn/Error
main -> store/factory.New -> csvstore.New(load CSV)
main -> ratelimit.Init
main -> findcountry.NewHandler -> httpserver.New -> Handler() -> ListenAndServe

Request:
client -> /v1/find-country
handler -> method check -> rate limit check -> ip validation
handler -> store.FindByIP
handler -> JSON success/error response
```

## Environment variables

- `APP_ENV` (default: `development`) selects a strongly typed Go preset (supported: `development`, `production`)
- `APP_NAME` (default: `torq-project`)
- `PORT` (default: `8080`)
- `DATASTORE_TYPE` (default: `csv`)
- `DATASTORE_PATH` (default: `data/ip_locations.csv`)
- `REQUESTS_PER_SECOND` (default: `10`)
- `RATE_LIMIT_BACKEND` (default: `inmemory`, supported: `inmemory`, `redis`)
- `REDIS_ADDR` (default: `localhost:6379`, used when backend is `redis`)
- `REDIS_PASSWORD` (default: empty, used when backend is `redis`)
- `REDIS_DB` (default: `0`, used when backend is `redis`)
- `REDIS_KEY_PREFIX` (default: `torq:ratelimit`, used when backend is `redis`)

Environment variables override values from the selected environment preset.

## Commands

```bash
make build
make run
make debug
make test
make lint
```

`make lint` runs formatting (`go fmt`), dependency cleanup (`go mod tidy`), static checks (`go vet`), and tests.

`make run` and `make debug` automatically run `ensure-redis` first:

- Ensures a local Redis container named `torq-redis` exists and is running.
- If missing, it creates one with `redis:7-alpine` on `localhost:6379`.
- If it exists but is stopped, it starts it.

`make debug` starts the app under Delve on `localhost:2345` (headless, ready for debugger attach and curl testing).

Debug flow:

1. Run `make debug`.
2. Attach debugger to `localhost:2345` from Cursor/VS Code Go debugger.
   - Configure a local debug profile with:
     - request: `attach`
     - mode: `remote`
     - host: `127.0.0.1`
     - port: `2345`
3. Set breakpoints and call:
   - `curl "http://localhost:8080/v1/find-country?ip=2.22.233.255"`

## Local example

```bash
make run
curl "http://localhost:8080/v1/find-country?ip=2.22.233.255"
```

## Docker build script

```bash
./scripts/build-docker.sh
```

Optional image naming:

```bash
IMAGE_NAME=my-org/ip2country IMAGE_TAG=v1 ./scripts/build-docker.sh
```

Optional build controls:

```bash
NO_CACHE=true ./scripts/build-docker.sh
DOCKERFILE_PATH=Dockerfile BUILD_CONTEXT=. ./scripts/build-docker.sh
```

Note: `./scripts/build-docker.sh` only builds a Docker image. It does not run the service.
For local execution, use `make run` (or `make debug`), which runs the app on your host and auto-ensures local Redis is available.

## Step-by-step local run

If you want to run everything in the common local flow, use:

1. Build Go service binary on host:

```bash
make build
```

2. Run or Debug service on host (also ensures local Redis container is ready):

```bash
make run
# or
make debug
```

3. Call the API:

```bash
curl "http://localhost:8080/v1/find-country?ip=2.22.233.255"
```
