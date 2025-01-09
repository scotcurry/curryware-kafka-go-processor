package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func InsertPlayerStats(statsJson []fantasyclasses.StatsInfo) {

	psqlInfo := GetDatabaseInformation()
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

	insertValues := ""
	for _, stat := range statsJson {
		statId := stat.StatId
		playerId := stat.PlayerID
		gameKey := stat.GameKey
		weekKey := stat.WeekKey
		statValue := stat.StatValue

		valueLine := "(" + strconv.Itoa(playerId) + "," + strconv.Itoa(gameKey) + "," + strconv.Itoa(weekKey) + "," + strconv.Itoa(statId) + "," + strconv.FormatFloat(statValue, 'f', 2, 64) + "),"
		insertValues = insertValues + valueLine
	}

	insertTemplate := GetSqlTemplate("multiple_play_stats_input_statement")
	sqlStatement := strings.ReplaceAll(insertTemplate, "{insert_values}", insertValues)
	sqlStatement = sqlStatement[:len(sqlStatement)-1]

	res, err := db.Exec(sqlStatement)
	if err != nil {
		fmt.Println("Error inserting player stats")
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
