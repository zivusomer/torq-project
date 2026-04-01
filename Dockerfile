FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
COPY data ./data

RUN go build -o /bin/torq-service ./cmd/torq-project

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /bin/torq-service /usr/local/bin/torq-service
COPY data ./data

ENV PORT=8080
ENV DATASTORE_TYPE=csv
ENV DATASTORE_PATH=/app/data/ip_locations.csv
ENV REQUESTS_PER_SECOND=10

EXPOSE 8080

CMD ["torq-service"]
