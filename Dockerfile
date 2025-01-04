# Build using docker build -t curryware-kafka-processor .
FROM golang:1.23
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o curryware-kafka-go-processor .
CMD ["./curryware-kafka-go-processor"]