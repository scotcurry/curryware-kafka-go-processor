package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"strconv"
)

func InsertTeamInformation(teamInfo []leagueclasses.TeamInformation) int {
	sqlStatement := `INSERT INTO team_information (team_key, team_name, team_url, team_logo_url,
                         draft_position, draft_grade, manager_nickname, manager_image_url, manager_felo_score)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT DO NOTHING`

	for counter := 0; counter < len(teamInfo); counter++ {
		teamKey := teamInfo[counter].TeamKey
		teamName := teamInfo[counter].TeamName
		teamUrl := teamInfo[counter].TeamUrl
		teamLogoUrl := teamInfo[counter].TeamLogoUrl
		draftPosition := teamInfo[counter].DraftPosition
		draftGrade := teamInfo[counter].DraftGrade
		managerNickname := teamInfo[counter].ManagerNickname
		managerImageUrl := teamInfo[counter].ManagerImageUrl
		managerFeloScore := teamInfo[counter].ManagerFeloScore

		count, err := ExecStatement(sqlStatement, teamKey, teamName, teamUrl, teamLogoUrl,
			draftPosition, draftGrade, managerNickname, managerImageUrl, managerFeloScore)
		if err != nil {
			logger.LogError("Error inserting team information record", "error", err.Error(), "team_key", teamKey)
			continue
		}
		logger.LogInfo("Rows affected", "count", strconv.Itoa(int(count)))
	}
	logger.LogInfo("Done inserting team information records", "total_records", strconv.Itoa(len(teamInfo)))
	return len(teamInfo)
}
