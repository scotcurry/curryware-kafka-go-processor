package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"strconv"
)

func InsertAllTeamInformation(teamInfo []leagueclasses.AllTeamInformation) int {
	sqlStatement := `INSERT INTO all_team_information (league_key, team_key, team_id, team_name, team_logo,
                         previous_season_team_rank, number_of_moves, number_of_trades, draft_position, draft_grade,
                         manager_nicknames)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT DO NOTHING`

	for counter := 0; counter < len(teamInfo); counter++ {
		team := teamInfo[counter]
		count, err := ExecStatement(sqlStatement,
			team.LeagueKey,
			team.TeamKey,
			team.TeamId,
			team.TeamName,
			team.TeamLogo,
			team.PreviousSeasonTeamRank,
			team.NumberOfMoves,
			team.NumberOfTrades,
			team.DraftPosition,
			team.DraftGrade,
			team.ManagerNicknames,
		)
		if err != nil {
			logger.LogError("Error inserting all team information record", "error", err.Error(), "team_key", team.TeamKey)
			continue
		}
		logger.LogInfo("Rows affected", "count", strconv.Itoa(int(count)))
	}
	logger.LogInfo("Done inserting all team information records", "total_records", strconv.Itoa(len(teamInfo)))
	return len(teamInfo)
}
