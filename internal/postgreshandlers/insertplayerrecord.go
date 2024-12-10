package postgreshandlers

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"strconv"
)

func InsertPlayerRecord(playerId int, playerKey string, playerName string, team string, position string, playerUrl string,
	playerHeadshot string) {

	//
	postgresServer := os.Getenv("POSTGRES_SERVER")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USERNAME")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDb := os.Getenv("POSTGRES_DATABASE")

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

	portInteger, err := strconv.ParseInt(postgresPort, 10, 64)
	if err != nil {
		fmt.Println("Error parsing port")
		panic(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		postgresServer, portInteger, postgresUser, postgresPassword, postgresDb)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error opening postgres connection")
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing postgres connection")
		}
	}(db)

	sqlStatement := `INSERT INTO player_info (player_id, player_season_key, player_name, player_position, player_team, 
                         player_url, player_headshot) 
					VALUES ($1, $2, $3, $4, $5, $6, $7)`

	res, err := db.Exec(sqlStatement, playerId, playerKey, playerName, position, team, playerUrl, playerHeadshot)
	if err != nil {
		fmt.Println("Error inserting player record")
		panic(err)
	} else {
		count, err := res.RowsAffected()
		if err != nil {
			fmt.Println("Error getting rows affected")
			panic(err)
		} else {
			fmt.Println("Rows affected: " + strconv.Itoa(int(count)))
		}
	}
}
