package main

import (
	"curryware-kafka-go-processor/internal/kafkahandlers"
	"curryware-kafka-go-processor/internal/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"os"
	"os/signal"
	"syscall"
)

// The documentation for the Kafka libraries are at https://pkg.go.dev/github.com/confluentinc/confluent-kafka-go/kafka
func main() {

	// Set up the Datadog Tracer.
	tracer.Start(tracer.WithService("curryware-kafka-go-processor"),
		tracer.WithServiceVersion("1.0.1"),
		tracer.WithEnv("prod"),
		tracer.WithTraceEnabled(true),
	)
	// The explanation for the defer keyword is at https://read.amazon.com/?asin=B0184N7WWS&ref_=kwl_kr_iv_rec_2
	defer tracer.Stop()

	logging.LogDebug("Debug message just because")
	// Create a channel that looks for all os Signals (https://pkg.go.dev/os/signal)
	sigChan := make(chan os.Signal, 1)
	// This limits the notification to the SIGTERM signal
	signal.Notify(sigChan, syscall.SIGTERM)
	// When the SIGTERM signal is received, it sends it to this go routine, which stops the tracer and closes the app.
	go func() {
		<-sigChan
		tracer.Stop()
		os.Exit(0)
	}()

	span := tracer.StartSpan("main")
	defer span.Finish()

	// Setting up logging.  JSON format.
	logging.LogInfo("Launching curryware-kafka-go-processor")

	// This is going to get the servername from the environment variable to make debugging easier.
	server := kafkahandlers.GetKafkaServer()

	// This code runs in a loop that is always true.  The syscall.SIGTERM above is the handler for breaking out
	// of this code.
	// logging.LogInfo("Starting Message Consumer for Topic: %s", topic)
	topicsToMonitor := kafkahandlers.GetTopicNames("localhost:9092")
	kafkahandlers.ConsumeMessages(topicsToMonitor, server)
}
