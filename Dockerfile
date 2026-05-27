# Build:  podman build -t curryware-kafka-go-processor:latest .
# Run:    podman run --env-file .env curryware-kafka-go-processor:latest

# --- build stage ---
FROM golang:1.26 AS builder
WORKDIR /app

ARG DD_GIT_COMMIT_SHA
ARG DD_GIT_REPOSITORY_URL

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux \
    go build \
    -ldflags="-s -w \
      -X main.ddGitCommitSHA=${DD_GIT_COMMIT_SHA} \
      -X main.ddGitRepositoryURL=${DD_GIT_REPOSITORY_URL}" \
    -o curryware-kafka-go-processor .

# --- final stage ---
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/curryware-kafka-go-processor .

ENV DD_GIT_COMMIT_SHA=${DD_GIT_COMMIT_SHA}
ENV DD_GIT_REPOSITORY_URL=${DD_GIT_REPOSITORY_URL}

CMD ["./curryware-kafka-go-processor"]