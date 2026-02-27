package postgreshandlers

import (
	logger "curryware-kafka-go-processor/internal/logging"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func GetDatabaseInformation() (string, error) {

	postgresServer := os.Getenv("POSTGRES_SERVER")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USERNAME")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DATABASE")

	if postgresServer == "" || postgresPort == "" || postgresUser == "" || postgresPassword == "" || postgresDb == "" {
		return "", errors.New("missing environment variables")
	}

	portInteger, err := strconv.ParseInt(postgresPort, 10, 64)
	if err != nil {
		logger.LogError("Error parsing port", "error", err.Error())
		return "", fmt.Errorf("invalid POSTGRES_PORT value %q: %w", postgresPort, err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		postgresServer, portInteger, postgresUser, postgresPassword, postgresDb)

	logger.LogInfo(fmt.Sprintf("psqlInfo %s", psqlInfo))

	return psqlInfo, nil
}
