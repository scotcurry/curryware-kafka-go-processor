package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"strconv"
)

func InsertTeamInformation(teamInfo []leagueclasses.TeamInformation) int {
	sqlStatement := `INSERT INTO all_team_information (league_key, team_key, team_id, team_name, team_logo,
                         previous_season_team_rank, number_of_moves, number_of_trades, draft_position, draft_grade, 
                         draft_recap_url, manager_nicknames, manager_felo_score)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) ON CONFLICT DO NOTHING`

	for counter := 0; counter < len(teamInfo); counter++ {
		leagueKey := teamInfo[counter].LeagueKey
		teamKey := teamInfo[counter].TeamKey
		teamId := teamInfo[counter].TeamID
		teamName := teamInfo[counter].TeamName
		teamUrl := teamInfo[counter].TeamLogo
		previousSeasonTeamRank := teamInfo[counter].PreviousSeasonTeamRank
		numberOfMoves := teamInfo[counter].NumberOfMoves
		numberOfTrades := teamInfo[counter].NumberOfTrades
		draftRecapURL := teamInfo[counter].DraftRecapURL
		draftPosition := teamInfo[counter].DraftPosition
		draftGrade := teamInfo[counter].DraftGrade
		managerNickname := teamInfo[counter].ManagerNicknames
		managerFeloScore := teamInfo[counter].ManagerFeloScore

		count, err := ExecStatement(sqlStatement, leagueKey, teamKey, teamId, teamName, teamUrl, previousSeasonTeamRank,
			numberOfMoves, numberOfTrades, draftPosition, draftGrade, draftRecapURL, managerNickname, managerFeloScore)
		if err != nil {
			logger.LogError("Error inserting team information record", "error", err.Error(), "team_key", teamKey)
			continue
		}
		logger.LogInfo("Rows affected", "count", strconv.Itoa(int(count)))
	}
	logger.LogInfo("Done inserting team information records", "total_records", strconv.Itoa(len(teamInfo)))
	return len(teamInfo)
}
