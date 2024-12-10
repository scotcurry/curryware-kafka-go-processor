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

	tracer.Start(tracer.WithAgentAddr("localhost:8126"),
		tracer.WithService("curryware-kafka-go-processor"),
		tracer.WithServiceVersion("1.0.0"),
		tracer.WithEnv("prod"),
		tracer.WithTraceEnabled(true),
	)

	defer tracer.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
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

	debug := false
	if len(os.Args) > 1 {
		argString := os.Args[1]
		if "-debug" == strings.ToLower(argString) {
			debug = true
		}
	}

	server := "ubuntu-postgres.curryware.org:9092"
	topic := "PlayerTopic"

	// This code is just to make sure that messages can be produced.  This module is set to pull messages off the
	// topic that are produced by the curryware-yahoo-api service.
	if debug {
		topicExists := kafkahandlers.ValidateTopicExists(topic, server)
		if topicExists {
			fmt.Println("Topic exists")
			kafkahandlers.ProduceMessage(topic, server)
		} else {
			fmt.Println("Topic does not exist")
			os.Exit(0)
		}
		fmt.Printf("Topic %s\n", topic)
	}

	kafkahandlers.ConsumeMessages(topic, server)
}
