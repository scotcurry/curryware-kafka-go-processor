package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"fmt"
	"strconv"
	"strings"
)

// InsertPlayerStats - This does just what it says.  The only thing that makes it unique is that since there are lots
// (test case pulls 310 stats) is it creates a single values string.  There is something with Go compiler check that
// case it to throw an error on the insert statement if it isn't complete filled out, so there is code to go pull
// a template from a file in sqltemplate.txt and use that.
func InsertPlayerStats(statsJson []fantasyclasses.StatsInfo) {
	// Use the singleton database connection pool
	db := GetDB()

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
			logger.LogError("Error inserting player stats", "error", err.Error(), "sql_statement", sqlStatement)
			panic(err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			logger.LogError("Error getting rows affected", "error", err.Error())
			panic(err)
		}
		logger.LogInfo("Rows affected", "count", count)
	}
}
