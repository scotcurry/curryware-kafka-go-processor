package kafkahandlers

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
)

func ValidateTopicExists(topic string, server string) bool {

	versionInt, libraryVersion := kafka.LibraryVersion()
	fmt.Printf("Kafka library version: %s, %d\n", libraryVersion, versionInt)

	// https://pkg.go.dev/github.com/confluentinc/confluent-kafka-go/kafka#NewAdminClient
	// I'm not sure what the & is for.  Need to research.
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": server,
	})

	if err != nil {
		fmt.Println("Error building admin client")
		panic(err)
	} else {
		fmt.Println("Admin client created, named - " + adminClient.String())
		defer adminClient.Close()
	}

	// Reading the documentation for this call - func (a *AdminClient) GetMetadata(topic *string, allTopics bool, timeoutMs int)
	// (*Metadata, error)
	// The last part of the documentation are the return values (metadata, err)
	var metadata, metaDataErr = adminClient.GetMetadata(nil, true, 10000)
	if metaDataErr != nil {
		return false
	}
	allTopics := metadata.Topics
	return topic == allTopics[topic].Topic
}

func GetKafkaServer() string {

	return os.Getenv("KAFKA_SERVER")
}

// CreateTopic This code doesn't work.
//func CreateTopic(topic string, server string) {
//
//	// Usual Kafka client creation
//	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
//		"bootstrap.servers": server,
//	})
//
//	if err != nil {
//		fmt.Println("Error building admin client")
//		panic(err)
//	} else {
//		fmt.Println("Admin client created, named - " + adminClient.String())
//		defer adminClient.Close()
//	}
//
//	// To create a topic you have to build this out.
//	topicSpec := []kafka.TopicSpecification{
//		{
//			Topic:             topic,
//			NumPartitions:     1,
//			ReplicationFactor: 1,
//		},
//	}
//
//	// Create the topic.
//	topics, err := adminClient.CreateTopics(nil, topicSpec)
//	if err != nil {
//		return
//	}
//	for _, topic := range topics {
//		fmt.Printf("Topic Created: %s\n", topic.String())
//	}
//}
