package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func consumeMessages(topic string, server string) {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": server,
		"group.id":          "curryware-group",
		"auto.offset.reset": "earliest",
	})
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
				fmt.Println("Failed to get watermark offsets: %s\n", err)
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
