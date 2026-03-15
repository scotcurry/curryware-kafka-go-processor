package main

import (
	"curryware-kafka-go-processor/internal/kafkahandlers"
	"curryware-kafka-go-processor/internal/logging"
	"curryware-kafka-go-processor/internal/postgreshandlers"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

// The documentation for the Kafka libraries are at https://pkg.go.dev/github.com/confluentinc/confluent-kafka-go/kafka
func main() {

	// Set up the Datadog Tracer.
	currentPath, err := os.Getwd()
	if err != nil {
		logging.LogError("Error getting current working directory", "error", err.Error())
	} else {
		logging.LogInfo("Current working directory", "path", currentPath)
	}

	err = tracer.Start(tracer.WithService("curryware-kafka-go-processor"),
		tracer.WithServiceVersion("1.0.1"),
		tracer.WithEnv("prod"),
		tracer.WithTraceEnabled(true),
	)
	if err != nil {
		return
	}
	// The explanation for defer keyword is at https://read.amazon.com/?asin=B0184N7WWS&ref_=kwl_kr_iv_rec_2
	defer tracer.Stop()

	logging.LogDebug("Starting curryware-kafka-go-processor...")
	// Create a channel that looks for all os Signals (https://pkg.go.dev/os/signal)
	sigChan := make(chan os.Signal, 1)
	// This limits the notification to the SIGTERM signal
	signal.Notify(sigChan, syscall.SIGTERM)
	// When the SIGTERM signal is received, it sends it to this go routine, which stops the tracer and closes the app.
	go func() {
		<-sigChan
		logging.LogInfo("Received SIGTERM signal, shutting down gracefully...")
		tracer.Stop()
		_ = postgreshandlers.CloseDB()
		os.Exit(0)
	}()

	span := tracer.StartSpan("main")
	defer span.Finish()

	// Setting up logging.  JSON format.
	logging.LogInfo("Launching curryware-kafka-go-processor...")

	// Log Postgres DNS/network diagnostics to help debug connectivity issues.
	postgresServer := os.Getenv("POSTGRES_SERVER")
	postgresPort := os.Getenv("POSTGRES_PORT")
	if postgresServer != "" {
		logging.LogInfo("Postgres server configured", "server", postgresServer, "port", postgresPort)
		addrs, dnsErr := net.LookupHost(postgresServer)
		if dnsErr != nil {
			logging.LogError("Postgres DNS resolution failed", "server", postgresServer, "error", dnsErr.Error())
		} else {
			logging.LogInfo("Postgres DNS resolved", "server", postgresServer, "addresses", fmt.Sprintf("%v", addrs))
			if postgresPort != "" {
				conn, dialErr := net.DialTimeout("tcp", net.JoinHostPort(postgresServer, postgresPort), 5*time.Second)
				if dialErr != nil {
					logging.LogError("Postgres TCP connection failed", "server", postgresServer, "port", postgresPort, "error", dialErr.Error())
				} else {
					logging.LogInfo("Postgres TCP connection successful", "server", postgresServer, "port", postgresPort)
					err := conn.Close()
					if err != nil {
						return
					}
				}
			}
		}
	} else {
		logging.LogError("POSTGRES_SERVER environment variable is not set")
	}

	// Attempt an early database connection. Log a warning if it fails but continue running
	// so the pod does not crash-loop when Postgres is temporarily unavailable.
	_, dbErr := postgreshandlers.GetDB()
	if dbErr != nil {
		logging.LogError("Postgres is not reachable at startup — the service will retry when a message arrives",
			"error", dbErr.Error())
	} else {
		logging.LogInfo("Postgres connection established at startup")
	}

	// This is going to get the servername from the environment variable to make debugging easier.
	server := kafkahandlers.GetKafkaServer()
	if server == "" {
		logging.LogError("Error getting Kafka server from environment variable (KAFKA_BOOTSTRAP_SERVER)")
		os.Exit(1001)
	}

	// This code runs in a loop that is always true.  The syscall.SIGTERM above is the handler for breaking out
	// of this code.
	topicsToMonitor, err := kafkahandlers.GetTopicNames(server)
	if err != nil {
		logging.LogError("Error getting topic names from Kafka server", "error", err.Error())
		os.Exit(1002)
	}
	kafkahandlers.ConsumeMessages(topicsToMonitor, server)
}
