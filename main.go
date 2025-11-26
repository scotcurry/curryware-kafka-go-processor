package main

import (
	"curryware-kafka-go-processor/internal/kafkahandlers"
	"curryware-kafka-go-processor/internal/logging"
	"curryware-kafka-go-processor/internal/postgreshandlers"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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

	tracer.Start(tracer.WithService("curryware-kafka-go-processor"),
		tracer.WithServiceVersion("1.0.1"),
		tracer.WithEnv("prod"),
		tracer.WithTraceEnabled(true),
	)
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

	// This is to "fail fast" if the database connection fails.
	_ = postgreshandlers.GetDB()

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
