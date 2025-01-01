package main

import (
	"curryware-kafka-go-processor/internal/kafkahandlers"
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// The documentation for the Kafka libraries are at https://pkg.go.dev/github.com/confluentinc/confluent-kafka-go/kafka
func main() {

	// Set up the Datadog Tracer.
	tracer.Start(tracer.WithAgentAddr("localhost:8126"),
		tracer.WithService("curryware-kafka-go-processor"),
		tracer.WithServiceVersion("1.0.0"),
		tracer.WithEnv("prod"),
		tracer.WithTraceEnabled(true),
	)
	// The explanation for the defer keyword is at https://read.amazon.com/?asin=B0184N7WWS&ref_=kwl_kr_iv_rec_2
	defer tracer.Stop()

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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Launching curryware-kafka-go-processor")
	log.Printf("Launching curryware-kafka-go-processor %v\n", span)

	defer logger.Info("Exiting curryware-kafka-go-processor")

	// Code to make sure I can run this code locally with a Kafka Docker image.
	debug := false
	if len(os.Args) > 1 {
		argString := os.Args[1]
		if "-debug" == strings.ToLower(argString) {
			debug = true
		}
	}

	// This is going to get the servername from the environment variable to make debugging easier.
	server := kafkahandlers.GetKafkaServer()
	// This should reflect the type of message we want to receive and how to receive it. CurrywareTopic is for
	// debugging.
	topic := "PlayerTopic2"

	// This code is just to make sure that messages can be produced.  This module is set to pull messages off the
	// topic that are produced by the curryware-yahoo-api service.
	if debug {
		topicExists := kafkahandlers.ValidateTopicExists(topic, server)
		if topicExists {
			fmt.Printf("DEBUG - Topic exists: %s", topic)
			kafkahandlers.ProduceMessage(topic, server)
		} else {
			fmt.Println("Topic does not exist")
			os.Exit(0)
		}
		fmt.Printf("Topic %s\n", topic)
	}
	fmt.Printf("Topic %s\n", topic)

	// This code runs in a loop that is always true.  The syscall.SIGTERM above is the handler for breaking out
	// of this code.
	kafkahandlers.ConsumeMessages(topic, server)
}
