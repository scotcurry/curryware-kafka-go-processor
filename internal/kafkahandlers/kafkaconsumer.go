package kafkahandlers

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	ddkafka "gopkg.in/DataDog/dd-trace-go.v1/contrib/confluentinc/confluent-kafka-go/kafka.v2"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ConsumeMessages(topic string, server string) {

	// Logging code for Datadog.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info(
		"Starting consumer",
		slog.String("topic", topic),
		slog.String("server", server),
	)
	defer logger.Info(
		"Stopping consumer",
		slog.String("topic", topic),
		slog.String("server", server),
	)

	consumer, err := ddkafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": server,
		"group.id":          "curryware-group",
		"auto.offset.reset": "earliest",
	}, ddkafka.WithDataStreams())
	if err != nil {
		fmt.Println("Error building consumer")
		panic(err)
	}

	err = consumer.Subscribe(topic, nil)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Jumping int the event loop")
	run := true
	for run {
		select {
		case sig := <-signalChannel:
			fmt.Printf("Caught signal %v, exiting\n", sig)
			run = false
		default:
			event, eventError := consumer.ReadMessage(1 * time.Second)
			if eventError != nil {
				continue
			}

			lowOffset, highOffSet, err := consumer.QueryWatermarkOffsets(
				*event.TopicPartition.Topic, event.TopicPartition.Partition, 3000)
			if err != nil {
				fmt.Printf("Failed to get watermark offsets: %s\n", err)
			} else {
				fmt.Printf("Low offset: %d, high offset: %d\n", lowOffset, highOffSet)
			}

			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				event.TopicPartition, string(event.Key), string(event.Value))
		}
	}

	closeError := consumer.Close()
	if closeError != nil {
		fmt.Println("Error closing consumer")
		panic(closeError)
	}
}
