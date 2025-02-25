package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// InsertPlayerStats - This does just what it says.  The only thing that makes it unique is that since there are lots
// (test case pulls 310 stats) is it creates a single values string.  There is something with Go compiler check that
// case it to throw an error on the insert statement if it isn't complete filled out, so there is code to go pull
// a template from a file in sqltemplate.txt and use that.
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

	insertTemplate := GetSqlTemplate("multiple_player_stats_input_statement")
	logger.LogDebug(fmt.Sprintf("Insert Template: %s\n", insertTemplate))
	if len(insertTemplate) > 0 {
		sqlStatement := strings.ReplaceAll(insertTemplate, "{insert_values}", insertValues)
		sqlStatement = sqlStatement[:len(sqlStatement)-1]
		logger.LogDebug(fmt.Sprintf("SQL Statement: %s\n", sqlStatement))

		res, err := db.Exec(sqlStatement)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error inserting player stats: SQL Statement: %s\n", sqlStatement))
			logger.LogError(err.Error())
			panic(err)
		} else {
			count, err := res.RowsAffected()
			if err != nil {
				fmt.Println("Error getting rows affected")
				panic(err)
			} else {
				logger.LogInfo(fmt.Sprintf("Rows affected: " + strconv.Itoa(int(count))))
			}
		}
	}
}
