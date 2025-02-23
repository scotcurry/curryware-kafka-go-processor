package postgreshandlers

import (
	logger "curryware-kafka-go-processor/internal/logging"
	"fmt"
	"os"
	"strconv"
)

func GetDatabaseInformation() string {

	postgresServer := os.Getenv("POSTGRES_SERVER")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USERNAME")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DATABASE")

	portInteger, err := strconv.ParseInt(postgresPort, 10, 64)
	if err != nil {
		fmt.Println("Error parsing port")
		panic(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		postgresServer, portInteger, postgresUser, postgresPassword, postgresDb)

	logger.LogInfo(fmt.Sprintf("psqlInfo %s", psqlInfo))

	return psqlInfo
}
