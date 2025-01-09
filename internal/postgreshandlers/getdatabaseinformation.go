package postgreshandlers

import (
	"fmt"
	"log/slog"
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Launching curryware-kafka-go-processor")
	logger.Debug(
		"Environment variables",
		slog.String("POSTGRES_SERVER", postgresServer),
		slog.String("POSTGRES_PORT", postgresPort),
		slog.String("POSTGRES_USERNAME", postgresUser),
		slog.String("POSTGRES_PASSWORD", postgresPassword),
		slog.String("POSTGRES_DATABASE", postgresDb),
	)
	defer logger.Info(
		"Stopping curryware-kafka-go-processor",
	)

	return psqlInfo
}
