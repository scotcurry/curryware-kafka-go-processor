FROM golang:1.23
WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go build -o curryware-kafka-go-processor .
CMD ["./curryware-kafka-go-processor"]