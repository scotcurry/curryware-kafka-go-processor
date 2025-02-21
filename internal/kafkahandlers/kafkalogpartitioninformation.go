package kafkahandlers

import (
	"curryware-kafka-go-processor/internal/logging"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// This is handled by calling the meta data
//func ListKafkaInformation(server string) {
//
//	allTopics := GetTopicNames(server)
//	fmt.Println("Total topics: " + fmt.Sprint(len(allTopics)))
//	printTopicMetaData(server, allTopics)
//}

// Not used anymore
//func printTopicMetaData(server string, topic []string) {
//
//	// Build a consumer.  Noting very interesting here. https://pkg.go.dev/github.com/confluentinc/confluent-kafka-go/v2/kafka@v2.8.0#NewConsumer
//	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
//		"bootstrap.servers": server,
//	})
//	if err != nil {
//		fmt.Println("Error building admin client")
//		panic(err)
//	}
//	defer adminClient.Close()
//
//	for counter := 0; counter < len(topic); counter++ {
//		fmt.Println("Topic: " + topic[counter])
//		metadata, err := adminClient.GetMetadata(&topic[counter], false, 10000)
//		if err != nil {
//			fmt.Println("Error getting metadata")
//		}
//		if metadata != nil {
//			topicPartitions := metadata.Topics
//			for topicCounter := 0; topicCounter < len(topicPartitions); topicCounter++ {
//				topicMetaData := topicPartitions[topic[counter]]
//				partitionMetaData := topicMetaData.Partitions
//				fmt.Println(partitionMetaData)
//			}
//
//		} else {
//			fmt.Println("No metadata")
//		}
//	}
//}

// GetTopicNames If you don't pass in a topic and ask for all topics they are returned using this call.
func GetTopicNames(server string) []string {

	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": server,
	})
	if err != nil {
		fmt.Println("Error building admin client")
		panic(err)
	}
	defer adminClient.Close()

	// Get all the topics and print them out.
	var allTopics []string
	metadata, err := adminClient.GetMetadata(nil, true, 10000)
	if err != nil {
		fmt.Println("Error getting metadata")
	}
	if metadata != nil {
		for topic := range metadata.Topics {
			if topic[0] != '_' {
				logging.LogInfo("Topic name: " + topic)
				allTopics = append(allTopics, topic)
			}
		}
	} else {
		fmt.Println("No metadata")
	}

	return allTopics
}
