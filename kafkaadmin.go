package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func validateTopicExists(topic string, server string) bool {

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
	}

	// Reading the documentation for this call - func (a *AdminClient) GetMetadata(topic *string, allTopics bool, timeoutMs int)
	// (*Metadata, error)
	// The last part of the documentation are the return values (metadata, err)
	var metadata, metaDataErr = adminClient.GetMetadata(nil, true, 10000)
	if metaDataErr != nil {
		return false
	} else {
		allTopics := metadata.Topics
		if allTopics[topic].Topic == topic {
			return true
		} else {
			return false
		}
	}
	return false
}
