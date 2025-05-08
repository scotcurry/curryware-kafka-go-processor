package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/logging"
	"database/sql"
	"fmt"
	"strconv"
)

func InsertLeagueStatInfo(leagueStatInfo []fantasyclasses.LeagueStatInfo) int {

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

	sqlStatement := `INSERT INTO league_scoring_information (league_stat_key, league_stat_id, league_stat_enabled, league_stat_name, 
                                        league_stat_display_name, league_stat_group, league_stat_abbreviation, league_stat_sort_order,
                                        league_stat_position_type, league_stat_sort_position) 
										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	for counter := 0; counter < len(leagueStatInfo); counter++ {

		leagueStatKey := leagueStatInfo[counter].LeagueStatId
		leagueStatId := leagueStatInfo[counter].StatId
		leagueStatEnabled := leagueStatInfo[counter].StatEnabled
		leagueStatName := leagueStatInfo[counter].StatName
		leagueStatDisplayName := leagueStatInfo[counter].StatDisplayName
		leagueStatGroup := leagueStatInfo[counter].StatGroup
		leagueStatAbbreviation := leagueStatInfo[counter].StatAbbreviation
		leagueStatSortOrder := leagueStatInfo[counter].StatSortOrder
		leagueStatPositionType := leagueStatInfo[counter].StatPositionType
		leagueStatSortPosition := leagueStatInfo[counter].StatSortPosition

		logging.LogInfo(fmt.Sprintf("League Stat Key: %d", leagueStatKey))
		res, err := db.Exec(sqlStatement, leagueStatKey, leagueStatId, leagueStatEnabled, leagueStatName, leagueStatDisplayName,
			leagueStatGroup, leagueStatAbbreviation, leagueStatSortOrder, leagueStatPositionType, leagueStatSortPosition)

		if err != nil {
			logging.LogError("Error inserting player record")
			fmt.Println("Error inserting player record")
			panic(err)
		} else {
			count, err := res.RowsAffected()
			if err != nil {
				logging.LogError("Error getting rows affected")
				fmt.Println("Error getting rows affected")
				panic(err)
			} else {
				logging.LogInfo("Rows affected", strconv.Itoa(int(count)))
				fmt.Println("Rows affected: " + strconv.Itoa(int(count)))
			}
		}
	}
	logging.LogInfo(fmt.Sprintf("Done inserting player records. Total records: %s", strconv.Itoa(len(leagueStatInfo))))
	return len(leagueStatInfo)
}
