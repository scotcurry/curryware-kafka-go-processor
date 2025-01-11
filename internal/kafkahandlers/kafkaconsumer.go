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

func ConsumeMessages(topic string, server string) {

	// Logging code for Datadog.
	logging.LogInfo("Launching curryware-kafka-go-processor")

	// Builds the consumer. Group ID will change for different types of statistics.
	consumer, err := ddkafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  server,
		"group.id":           "curryware-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "true",
	}, ddkafka.WithDataStreams())
	if err != nil {
		logging.LogError("Error building consumer")
		fmt.Println("Error building consumer")
		panic(err)
	}

	// List where the last commit happened.  To do this you need to have to pass in the TopicPartitions, so get those
	// first.
	err = consumer.SubscribeTopics([]string{topic}, nil)
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
				continue
			}
			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				event.TopicPartition, string(event.Key), string(event.Value))
			logging.LogInfo("Consumed event from topic %s: key = %-10s value = %s",
				event.TopicPartition, string(event.Key), string(event.Value))

			valueAsString := string(event.Value)
			printValue := valueAsString[10 : 20+60]
			fmt.Println(printValue)
			playersToAdd := jsonhandlers.ParseMultiplePlayerInfo(valueAsString)
			postgreshandlers.InsertPlayerRecord(playersToAdd)

			offsets := []kafka.TopicPartition{
				{Topic: event.TopicPartition.Topic, Partition: event.TopicPartition.Partition,
					Offset: event.TopicPartition.Offset + 1},
			}
			_, err = consumer.CommitOffsets(offsets)
			if err != nil {
				logging.LogError("Failed to commit offsets: %s", err)
				fmt.Printf("Failed to commit offsets: %s\n", err)
			}

		}
	}

	closeError := consumer.Close()
	if closeError != nil {
		fmt.Println("Error closing consumer")
		logging.LogError("Error closing consumer")
		panic(closeError)
	}
}
