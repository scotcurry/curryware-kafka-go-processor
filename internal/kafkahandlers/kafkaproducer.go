package kafkahandlers

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	ddkafka "gopkg.in/DataDog/dd-trace-go.v1/contrib/confluentinc/confluent-kafka-go/kafka.v2"
	"strings"
)

func ProduceMessage(topic string, server string) {

	// Link to the sample code to create a producer - https://developer.confluent.io/get-started/go/#build-producer
	// Link to the producer documentation - https://pkg.go.dev/github.com/confluentinc/confluent-kafka-go/kafka#Producer
	producer, err := ddkafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": server,
	}, ddkafka.WithDataStreams())
	if err != nil {
		panic(err)
	}

	// The defer keyword is the equivalent of the using keyword in .NET.  It will run the code when
	// the produceMessages gets to the end.
	defer producer.Close()

	deliveryChan := make(chan kafka.Event)

	key := "producer_key"
	keyValue := `[{"Key":"449.p.40300","Id":40300,"FullName":"Cephus Johnson III","Url":"https://sports.yahoo.com/nfl/players/40300","Status":"NA","Team":"TB","ByeWeek":0,"UniformNumber":0,"Position":"WR","Headshot":"https://s.yimg.com/iu/api/res/1.2/8YBaY7mta1OqebkDNJ9JUg--~C/YXBwaWQ9eXNwb3J0cztjaD0yMzM2O2NyPTE7Y3c9MTc5MDtkeD04NTc7ZHk9MDtmaT11bGNyb3A7aD02MDtxPTEwMDt3PTQ2/https://s.yimg.com/xe/i/us/sp/v/nfl_cutout/players_l/08092023/40300.png"}]`

	var topicPartition = kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}
	err = producer.Produce(&kafka.Message{
		Key:            []byte(key),
		Value:          []byte(keyValue),
		TopicPartition: topicPartition,
	}, deliveryChan)

	if err != nil {
		fmt.Println("Failed to produce message:", err)
		panic(err)
	}

	e := <-deliveryChan
	msg := e.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", msg.TopicPartition)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
	}
	messageLabel := "Producer Name:"
	message := strings.Join([]string{messageLabel, producer.String()}, " ")
	fmt.Println(message)
}
