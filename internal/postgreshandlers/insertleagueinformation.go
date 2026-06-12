package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"strconv"
	"time"
)

func InsertLeagueInformation(leagueInfo []leagueclasses.LeagueInformation) int {
	sqlStatement := `INSERT INTO all_league_information (league_key, league_id, game_id, league_name, league_logo_url,
                         number_of_teams, league_update_timestamp, start_date, end_week, season)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT DO NOTHING`

	for counter := 0; counter < len(leagueInfo); counter++ {
		gameId, err := strconv.Atoi(leagueInfo[counter].GameKey)
		if err != nil {
			logger.LogError("Error converting GameKey to integer", "error", err.Error(), "game_key", leagueInfo[counter].GameKey)
			continue
		}
		leagueKey := leagueInfo[counter].LeagueKey
		leagueId := leagueInfo[counter].LeagueId
		leagueName := leagueInfo[counter].LeagueName
		leagueLogoUrl := leagueInfo[counter].LeagueLogoUrl
		numberOfTeams := leagueInfo[counter].NumberOfTeams
		leagueUpdateTimestamp := parseDateTime(leagueInfo[counter].LeagueUpdateTimestamp)
		startDate := leagueInfo[counter].StartDate
		endDate := leagueInfo[counter].EndDate
		season := leagueInfo[counter].Season

		count, err := ExecStatement(sqlStatement, leagueKey, leagueId, gameId, leagueName, leagueLogoUrl,
			numberOfTeams, leagueUpdateTimestamp, startDate, endDate, season)
		if err != nil {
			logger.LogError("Error inserting league information record", "error", err.Error(), "league_key", leagueKey)
			continue
		}
		if count == 0 {
			logger.LogInfo("League information record already exists, skipping remaining records", "league_key", leagueKey)
			return counter
		}
		logger.LogInfo("Rows affected", "count", strconv.Itoa(int(count)))
	}
	logger.LogInfo("Done inserting league information records", "total_records", strconv.Itoa(len(leagueInfo)))
	return len(leagueInfo)
}

func parseDateTime(dateTimeString string) time.Time {
	// C# DateTime serializes in ISO 8601 format
	parsed, err := time.Parse(time.RFC3339, dateTimeString)
	if err != nil {
		// Try without timezone
		parsed, err = time.Parse("2006-01-02T15:04:05", dateTimeString)
		if err != nil {
			logger.LogError("Error parsing date time", "error", err.Error(), "value", dateTimeString)
			return time.Time{}
		}
	}
	return parsed
}
