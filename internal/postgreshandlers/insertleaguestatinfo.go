package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"strconv"
)

func InsertLeagueStatInfo(leagueStatInfo []leagueclasses.LeagueStatDescriptionInfo) int {
	sqlStatement := `INSERT INTO league_stat_description (league_stat_key_id, game_id, league_id, stat_id, stat_enabled, stat_name,
                                        stat_display_name, stat_group_display_name, stat_abbreviation, stat_sort_order,
                                        stat_position_type, stat_sort_position)
										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for counter := 0; counter < len(leagueStatInfo); counter++ {
		stat := leagueStatInfo[counter]

		logger.LogDebug("Inserting league stat", "league_stat_key_id", stat.LeagueStatKeyId)
		count, err := ExecStatement(sqlStatement, stat.LeagueStatKeyId, stat.GameId, stat.LeagueId, stat.StatId,
			stat.StatEnabled, stat.StatName, stat.StatDisplayName, stat.StatGroupDisplayName,
			stat.StatAbbreviation, stat.StatSortOrder, stat.StatPositionType, stat.StatSortPosition)
		if err != nil {
			logger.LogError("Error inserting league stat record", "error", err.Error())
			continue
		}
		logger.LogInfo("Rows affected", "count", strconv.Itoa(int(count)))
	}
	logger.LogInfo("Done inserting league stat records", "total_records", strconv.Itoa(len(leagueStatInfo)))
	return len(leagueStatInfo)
}
