package kafkahandlers

import (
	"curryware-kafka-go-processor/internal/logging"
	"fmt"
	"os"
)

// This is handled by querying the Kafka server for topics
//func ValidateTopicExists(topics []string, server string) []string {
//
//	versionInt, libraryVersion := kafka.LibraryVersion()
//	logging.LogInfo("Kafka library version: %s, %d", libraryVersion, versionInt)
//	fmt.Printf("Kafka library version: %s, %d\n", libraryVersion, versionInt)
//
//	// https://pkg.go.dev/github.com/confluentinc/confluent-kafka-go/kafka#NewAdminClient
//	// I'm not sure what the & is for.  Need to research.
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
//	// Reading the documentation for this call - func (a *AdminClient) GetMetadata(topic *string, allTopics bool, timeoutMs int)
//	// (*Metadata, error)
//	// The last part of the documentation are the return values (metadata, err)
//	var existingTopics []string
//	var metadata, metaDataErr = adminClient.GetMetadata(nil, true, 10000)
//	if metaDataErr != nil {
//		return existingTopics
//	}
//
//	for _, currentTopic := range topics {
//		for _, currentMetadataTopic := range metadata.Topics {
//			if currentTopic == currentMetadataTopic.Topic {
//				existingTopics = append(existingTopics, currentTopic)
//				logging.LogInfo("Topic %s exists", currentTopic)
//			}
//		}
//	}
//
//	return existingTopics
//}

func GetKafkaServer() string {

	kafkaServer := os.Getenv("KAFKA_SERVER")
	logging.LogInfo(fmt.Sprintf("KAFKA_SERVER: %s", kafkaServer))
	return kafkaServer
}
