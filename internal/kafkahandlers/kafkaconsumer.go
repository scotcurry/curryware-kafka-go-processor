package kafkahandlers

import (
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"curryware-kafka-go-processor/internal/logging"
	"curryware-kafka-go-processor/internal/postgreshandlers"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	ddkafka "gopkg.in/DataDog/dd-trace-go.v1/contrib/confluentinc/confluent-kafka-go/kafka.v2"
)

func ConsumeMessages(topics []string, server string) {

	// Logging code for Datadog.
	logging.LogInfo("Launching curryware-kafka-go-processor")
	for i := 0; i < len(topics); i++ {
		fmt.Println(fmt.Sprintf("Consuming Message(s): %s", topics[i]))
	}

	// Builds the consumer. Group ID will change for different types of statistics.
	consumer, err := ddkafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  server,
		"group.id":           "curryware-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "true",
	}, ddkafka.WithDataStreams())
	defer func(consumer *ddkafka.Consumer) {
		closeErr := consumer.Close()
		if closeErr != nil {
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
				continue
			} else {
				logging.LogInfo("Event received %d bytes", len(event.Value))
			}

			switch *event.TopicPartition.Topic {
			case "PlayerStats":
				logging.LogInfo("Processing PlayerStats")
				statPackage := string(event.Value)
				statsInfo := jsonhandlers.ParseMultipleStatInfo(statPackage)
				postgreshandlers.InsertPlayerStats(statsInfo)
				break
			case "PlayerTopic2":
				logging.LogInfo("Processing PlayerTopic2")
				playerPackage := string(event.Value)
				playersToAdd := jsonhandlers.ParseMultiplePlayerInfo(playerPackage)
				postgreshandlers.InsertPlayerRecord(playersToAdd)
				break
			case "StatTopic":
				logging.LogInfo("Processing StatTopic")
				statInfoPackage := string(event.Value)
				statInfoToAdd := jsonhandlers.ParseLeagueStatInfo(statInfoPackage)
				postgreshandlers.InsertLeagueStatInfo(statInfoToAdd)
				break
			case "StatValueTopic":
				logging.LogInfo("Processing StatValueTopic")
				statValuePackage := string(event.Value)
				statValuesToAdd := jsonhandlers.ParseLeagueStatValue(statValuePackage)
				postgreshandlers.InsertLeagueStatValueInfo(statValuesToAdd)
				break
			case "PlayerTopicDaily":
				logging.LogInfo("Processing PlayerTopicDaily")
				playerPackage := string(event.Value)
				logging.LogInfo("Player package length: ", len(playerPackage))
				break
			case "DatadogValidationTopic":
				logging.LogInfo("Processing DatadogValidationTopic")
				dataValidationPackage := string(event.Value)
				logging.LogInfo("Data validation package length: ", len(dataValidationPackage))
				break
			case "TransactionTopic":
				logging.LogInfo("Processing TransactionTopic")
				transactionPackage := string(event.Value)
				transactionJson := jsonhandlers.ParseTransactionInfo(transactionPackage)
				transactionCount := postgreshandlers.ProcessTransactionInfo(transactionJson)
				logging.LogDebug("NEEDS TO BE UPDATED: Transaction count: ", transactionCount)
				logging.LogInfo("Transaction package length: ", len(transactionPackage))
			default:
				logging.LogError(fmt.Sprintf("Unknown topic - %s", *event.TopicPartition.Topic))
				fmt.Println(fmt.Sprintf("Unknown topic - %s", *event.TopicPartition.Topic))
			}
		}
	}
}
