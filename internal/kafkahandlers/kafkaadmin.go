package kafkahandlers

import (
	"curryware-kafka-go-processor/internal/logging"
	"fmt"
	"os"
)

func GetKafkaServer() string {

	logging.LogDebug("Reading environment variable: KAFKA_BOOTSTRAP_SERVER")
	kafkaServer := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	logging.LogInfo(fmt.Sprintf("KAFKA_BOOTSTRAP_SERVER: %s", kafkaServer))
	return kafkaServer
}
