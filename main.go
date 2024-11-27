package main

import (
	"fmt"
	"os"
)

// The documentation for the Kafka libraries are at https://pkg.go.dev/github.com/confluentinc/confluent-kafka-go/kafka
func main() {

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
