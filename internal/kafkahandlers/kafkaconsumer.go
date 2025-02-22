package kafkahandlers

import (
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"curryware-kafka-go-processor/internal/logging"
	"curryware-kafka-go-processor/internal/postgreshandlers"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	ddkafka "gopkg.in/DataDog/dd-trace-go.v1/contrib/confluentinc/confluent-kafka-go/kafka.v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ConsumeMessages(topics []string, server string) {

	// Logging code for Datadog.
	logging.LogInfo("Launching curryware-kafka-go-processor")

	// Builds the consumer. Group ID will change for different types of statistics.
	consumer, err := ddkafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  server,
		"group.id":           "curryware-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "true",
	}, ddkafka.WithDataStreams())
	defer func(consumer *ddkafka.Consumer) {
		err := consumer.Close()
		if err != nil {
			logging.LogError("Error closing consumer")
		}
	}(consumer)
	if err != nil {
		logging.LogError("Error building consumer")
		fmt.Println("Error building consumer")
		panic(err)
	}

	// List where the last commit happened.  To do this you need to have to pass in the TopicPartitions, so get those
	// first.
	err = consumer.SubscribeTopics(topics, nil)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// This is the loop that will run forever.  Need to use Datadog to see how much processor this actually takes.
	fmt.Println("Jumping into the event loop")
	run := true
	for run {
		select {
		case sig := <-signalChannel:
			fmt.Printf("Caught signal %v, exiting\n", sig)
			run = false
		default:
			event, eventError := consumer.ReadMessage(20 * time.Second)
			if eventError != nil {
				logging.LogError("Error reading message: %s", eventError)
				continue
			}

			switch *event.TopicPartition.Topic {
			case "PlayerStats":
				statPackage := string(event.Value)
				statsInfo := jsonhandlers.ParseMultipleStatInfo(statPackage)
				postgreshandlers.InsertPlayerStats(statsInfo)
				break
			case "PlayerTopics2":
				playerPackage := string(event.Value)
				playersToAdd := jsonhandlers.ParseMultiplePlayerInfo(playerPackage)
				postgreshandlers.InsertPlayerRecord(playersToAdd)
				break
			default:
				fmt.Println(fmt.Sprintf("Unknown topic - %s", *event.TopicPartition.Topic))
			}
		}
	}
}
