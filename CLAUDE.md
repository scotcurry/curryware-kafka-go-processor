# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Platform and Requirements
- Go 1.26+
- PostgreSQL
- Docker/Podman (Podman machine: `curryware-podman`)

## Build and Validation Commands
**Important: After every code change, validate the build succeeds.**

```bash
# Build
go build -o curryware-kafka-go-processor .

# Run all tests
go test ./...

# Run tests in a specific package
go test ./internal/jsonhandlers/...
go test ./internal/postgreshandlers/...
go test ./internal/kafkahandlers/...

# Run a single test
go test ./internal/jsonhandlers/... -run TestFunctionName

# Run tests in the moved stats tests directory
go test ./internal/tests/statstests/...

# After build and tests pass, verify Podman is running then build and validate the Docker image
podman machine info curryware-podman   # Ensure the curryware-podman machine is running
podman build -t curryware-kafka-go-processor:latest .   # Build the Docker image
podman rmi curryware-kafka-go-processor:latest           # Clean up the image after a successful build
```

## Architecture Overview

This is a Kafka consumer service for a fantasy football data pipeline. It consumes messages from Kafka topics, parses JSON payloads, and persists records to PostgreSQL. The service is instrumented with Datadog APM tracing (dd-trace-go) and deployed as a Docker container to Kubernetes.

### Message Flow

`main.go` → `kafkahandlers.ConsumeMessages()` → topic handler functions → `jsonhandlers` parsers → `postgreshandlers` inserters

The consumer runs an infinite event loop dispatching messages to handlers registered in `topicHandlers` map by topic name:
- `PlayerTopicDaily` → parse player roster → insert player records
- `TransactionTopic` → parse transaction JSON → insert transaction info
- `StatisticsTopic` → parse stats JSON → insert player weekly stats
- `DatadogValidationTopic` → validation/heartbeat (no DB write)

### Package Structure

- **`internal/kafkahandlers/`** — Kafka consumer (with Datadog instrumentation via `ddkafka`), producer, DNS/port validation, topic/partition logging
- **`internal/jsonhandlers/`** — JSON parsing for each entity type; all Kafka message payloads are **base64-encoded JSON**. `parsejson.go` contains the generic `ParseJSON[T]` function using Go generics.
- **`internal/postgreshandlers/`** — DB connection singleton (`getdatabaseconnection.go`), SQL template loader (`getsqltemplate.go`), and per-entity insert functions
- **`internal/fantasyclasses/`** — Data model structs, organized into sub-packages:
  - `playerclasses/` — `PlayerInfo`
  - `statsclasses/` — `PlayerWeeklyStatsInfo`
  - `transactionclasses/` — `TransactionInfo`, `TransactionPlayerInfo`
  - `leagueclasses/` — `LeagueStatDescriptionInfo`, `LeagueStatValues`
- **`internal/logging/`** — Singleton `slog` logger writing JSON to stdout
- **`internal/tests/`** — Test files separated from implementation (e.g., `statstests/`)

### SQL Templates

Rather than embedding SQL strings in Go code, multi-row `INSERT` statements are stored in `internal/postgreshandlers/sqltemplates/sqltemplate.txt`. `GetSqlTemplate(name)` reads this file at runtime by matching a template name prefix on each line. The `{insert_values}` placeholder is replaced with dynamically built value strings.

### Environment Variables

The service requires these at runtime:
- `KAFKA_BOOTSTRAP_SERVER` — Kafka broker address (e.g., `broker:9092`)
- `POSTGRES_SERVER`, `POSTGRES_PORT`, `POSTGRES_USERNAME`, `POSTGRES_PASSWORD`, `POSTGRES_DATABASE` — PostgreSQL connection details

### CI/CD

GitHub Actions (`.github/workflows/build-actions.yaml`) builds a Docker image on push to `master`, pushes to Docker Hub as `scotcurry4/curryware-kafka-go-processor:<run_number>`, then updates the image tag in the `scotcurry/k8s-manifests` repo for GitOps deployment.

Local Docker build: `docker build -t curryware-kafka-processor .`
