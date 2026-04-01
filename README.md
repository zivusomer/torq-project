# torq-project

IP-to-country service in Go with a pluggable datastore and built-in rate limiting.

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
3. `internal/logging` builds logger from `LOG_LEVEL`.
4. `internal/store/factory` selects datastore implementation (currently `csv`).
5. `internal/store/csvstore` loads CSV rows into an in-memory map.
6. `internal/ratelimit` creates the per-second limiter.
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
main -> logging.New
main -> store/factory.New -> csvstore.New(load CSV)
main -> ratelimit.New
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
- `LOG_LEVEL` (default: `info`)
- `PORT` (default: `8080`)
- `DATASTORE_TYPE` (default: `csv`)
- `DATASTORE_PATH` (default: `data/ip_locations.csv`)
- `REQUESTS_PER_SECOND` (default: `10`)

Environment variables override values from the selected environment preset.

## Commands

```bash
make build
make run
make test
make lint
```

`make lint` runs formatting (`go fmt`), dependency cleanup (`go mod tidy`), static checks (`go vet`), and tests.

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
