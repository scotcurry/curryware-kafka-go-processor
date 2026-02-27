package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"strconv"
)

func InsertLeagueStatInfo(leagueStatInfo []leagueclasses.LeagueStatDescriptionInfo) int {
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

		logger.LogDebug("Inserting league stat", "league_stat_key", leagueStatKey)
		count, err := ExecStatement(sqlStatement, leagueStatKey, leagueStatId, leagueStatEnabled, leagueStatName, leagueStatDisplayName,
			leagueStatGroup, leagueStatAbbreviation, leagueStatSortOrder, leagueStatPositionType, leagueStatSortPosition)
		if err != nil {
			logger.LogError("Error inserting league stat record", "error", err.Error())
			continue
		}
		logger.LogInfo("Rows affected", "count", strconv.Itoa(int(count)))
	}
	logger.LogInfo("Done inserting league stat records", "total_records", strconv.Itoa(len(leagueStatInfo)))
	return len(leagueStatInfo)
}
