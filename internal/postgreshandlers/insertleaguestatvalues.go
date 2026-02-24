package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/logging"
	"database/sql"
	"fmt"
	"strconv"
)

// InsertLeagueStatValueInfo TODO:  This is complete function needs to be revisited.  Completely wrong.
func InsertLeagueStatValueInfo(leagueStatValueInfo []fantasyclasses.PlayerStatValueInfo) int {

	psqlInfo, variableError := GetDatabaseInformation()
	if variableError != nil {
		logging.LogError("Error getting database information")
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logging.LogInfo("Error opening postgres connection")
		fmt.Println("Error opening postgres connection")
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing postgres connection")
		}
	}(db)

	sqlStatement := `INSERT INTO league_scoring_value_information (game_id, league_id, 
                                              league_stat_id, league_stat_value) 
					VALUES ($1, $2, $3, $4)`

	for counter := 0; counter < len(leagueStatValueInfo); counter++ {

		leagueGameId := leagueStatValueInfo[counter].PlayerGameKey
		leagueId := leagueStatValueInfo[counter].StatId
		leagueStatId := leagueStatValueInfo[counter].StatId
		leagueStatValue := leagueStatValueInfo[counter].StatValue

		res, err := db.Exec(sqlStatement, leagueGameId, leagueId, leagueStatId, leagueStatValue)
		if err != nil {
			logging.LogError("Error inserting stat record")
			fmt.Println("Error inserting stat record")
			panic(err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			logging.LogError("Error getting stat rows affected")
			fmt.Println("Error getting stat rows affected")
			panic(err)
		}

		logging.LogInfo("Stat rows affected", strconv.Itoa(int(count)))
		fmt.Println("Stat rows affected: " + strconv.Itoa(int(count)))
	}
	logging.LogInfo(fmt.Sprintf("Done inserting player records. Total records: %s",
		strconv.Itoa(len(leagueStatValueInfo))))

	return len(leagueStatValueInfo)
}
