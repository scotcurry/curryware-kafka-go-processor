package kafkahandlers

import (
	"curryware-kafka-go-processor/internal/logging"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

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
