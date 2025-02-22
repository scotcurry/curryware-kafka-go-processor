package kafkahandlers

import (
	"curryware-kafka-go-processor/internal/logging"
	"fmt"
	"os"
)

func GetKafkaServer() string {

	kafkaServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	logging.LogInfo(fmt.Sprintf("KAFKA_BOOTSTRAP_SERVER: %s", kafkaServer))
	return kafkaServer
}
