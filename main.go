package main

import (
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"os"
)

// The documentation for the Kafka libraries are at https://pkg.go.dev/github.com/confluentinc/confluent-kafka-go/kafka
func main() {

	tracer.Start(tracer.WithAgentAddr("localhost:8126"),
		tracer.WithService("curryware-kafka-go-processor"),
		tracer.WithServiceVersion("1.0.0"),
		tracer.WithEnv("prod"))

	defer tracer.Stop()

	server := "ubuntu-postgres.curryware.org:9092"
	topic := "PlayerTopic"

	topicExists := validateTopicExists(topic, server)
	if topicExists {
		fmt.Println("Topic exists")
		produceMessage(topic, server)
	} else {
		fmt.Println("Topic does not exist")
		os.Exit(0)
	}
	fmt.Printf("Topic %s\n", topic)
	consumeMessages(topic, server)
}
