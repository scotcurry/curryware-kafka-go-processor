# Build using docker build -t curryware-kafka-processor .

# Run local with Docker container agent.
#docker run -d --name dd-agent \
#-e DD_API_KEY=<API Key> \
#-e DD_SITE="datadoghq.com" \
#-e DD_LOG_LEVEL="debug" \
#-e DD_APM_ENABLED=true \
#-e DD_LOGS_ENABLED=true \
#-e DD_LOGS_CONFIG_CONTAINER_COLLECT_ALL=true \
#-e DD_LOGS_CONFIG_DOCKER_CONTAINER_USE_FILE=true \
#-v /var/run/docker.sock:/var/run/docker.sock:ro \
#-v /proc/:/host/proc/:ro \
#-v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
#-v /var/lib/docker/containers:/var/lib/docker/containers:ro \
#-v /opt/datadog-agent/run:/opt/datadog-agent/run:rw \
#-p 127.0.0.1:8126:8126/tcp \
#gcr.io/datadoghq/agent:7

FROM golang:1.23
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o curryware-kafka-go-processor .
CMD ["./curryware-kafka-go-processor"]