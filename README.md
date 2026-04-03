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

## Design notes

### Identity and rate limiting (IP-based)

- This API does not require user authentication tokens.
- Caller identity is derived from network metadata (`X-Forwarded-For`, then `X-Real-IP`, then `RemoteAddr`).
- Rate limiting is applied per caller identity key (currently the resolved caller IP).
- On limit exceed, the API returns `429` and includes standard throttling headers such as `Retry-After` and `RateLimit-*`.
- This design keeps the API simple while still protecting backend resources from abuse.

### Backend rate-limit implementations

- Two interchangeable backends are supported behind one facade:
  - in-memory token bucket (single-process/local)
  - Redis token bucket (distributed/multi-instance)
- Redis backend uses an atomic Lua script to avoid race conditions under concurrent requests.
- Switching backends is configuration-driven (`RATE_LIMIT_BACKEND`) without changing handler code.

### Datastore design (CSV now, pluggable later)

- Current datastore is CSV (`data/ip_locations.csv`) for easy local testing.
- Business logic reads through the store abstraction (`store.Resolver`), not directly from CSV parsing code.
- The active datastore is selected by configuration (`DATASTORE_TYPE`), so adding a new backend (DB/service/file format) is straightforward and does not require API contract changes.

### Logging behavior

- Logging is centralized via the internal logging package and exposed through simple calls (`Info`, `Warn`, `Error`).
- Logs are used for startup lifecycle, configuration/bootstrap failures, and runtime issues (including middleware/runtime guardrails).
- The current behavior is intentionally verbose enough for local debugging and external validation.

### Architecture at a glance

```text
Client
  -> HTTP route: /v1/find-country
  -> API middleware pipeline
      1) caller identity extraction (IP from headers/remote address)
      2) rate limit check (in-memory or Redis backend)
      3) request input validation (query param ip)
      4) find-country execution via store.Resolver
      5) unified JSON response writer
  -> Response: 200/400/404/429/500 with JSON body

Startup path:
config.LoadFromEnv (preset + env overrides)
  -> app bootstrap
  -> datastore factory (CSV implementation today)
  -> ratelimit backend init (inmemory/redis)
  -> HTTP server start
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
make run # executes `ensure-redis` first (below) 
make debug # executes `ensure-redis` (below), and starts the app under Delve on `localhost:2345`
make test
make lint # runs format (`go fmt`), dependency cleanup (`go mod tidy`), static checks (`go vet`), and tests
ensure-redis # Uses local Redis container or creates one with `redis:7-alpine` on `localhost:6379`
```

`make run` and `make debug` also auto-clean stale listeners before start:

- `make run` frees port `8080` if needed.
- `make debug` frees ports `8080` and `2345` if needed.

## Debug flow:

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
