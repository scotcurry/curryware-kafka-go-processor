package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/logging"
	"database/sql"
	"fmt"
	"strconv"
)

func InsertLeagueStatValueInfo(leagueStatValueInfo []fantasyclasses.LeagueStatsValueInfo) int {

	psqlInfo := GetDatabaseInformation()
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

	sqlStatement := `INSERT INTO league_scoring_value_information (league_stat_key, league_game_id, league_id, 
                                              league_stat_id, league_stat_value) 
					VALUES ($1, $2, $3, $4, $5)`

	for counter := 0; counter < len(leagueStatValueInfo); counter++ {

		leagueStatKey := leagueStatValueInfo[counter].LeagueStatId
		leagueStatId := leagueStatValueInfo[counter].StatId
		leagueGameId := leagueStatValueInfo[counter].GameId
		leagueId := leagueStatValueInfo[counter].LeagueId
		leagueStatValue := leagueStatValueInfo[counter].StatValue

		res, err := db.Exec(sqlStatement, leagueStatKey, leagueStatId, leagueGameId, leagueId, leagueStatValue)
		if err != nil {
			logging.LogError("Error inserting stat record")
			fmt.Println("Error inserting stat record")
			panic(err)
		} else {
			count, err := res.RowsAffected()
			if err != nil {
				logging.LogError("Error getting stat rows affected")
				fmt.Println("Error getting stat rows affected")
				panic(err)
			} else {
				logging.LogInfo("Stat rows affected", strconv.Itoa(int(count)))
				fmt.Println("Stat rows affected: " + strconv.Itoa(int(count)))
			}
		}
	}
	logging.LogInfo(fmt.Sprintf("Done inserting player records. Total records: %s",
		strconv.Itoa(len(leagueStatValueInfo))))

	return len(leagueStatValueInfo)
}
