package postgreshandlers

import (
	"context"
	"curryware-kafka-go-processor/internal/fantasyclasses/playerclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"strconv"
	"strings"
)

// inactiveStatusFull is the PlayerStatusFull value for players that are not on a roster.
// Any player whose PlayerStatusFull contains this text is skipped and not persisted.
const inactiveStatusFull = "Inactive: Coach's Decision or Not on Roster"

// InsertAvailablePlayerRecord persists available players from the AllAvailablePlayersTopic into
// the player_info table. Players whose PlayerStatusFull indicates they are inactive / not on a
// roster are skipped. It returns the number of records inserted.
func InsertAvailablePlayerRecord(playerInfo []playerclasses.PlayerInfo) int {
	sqlStatement := `INSERT INTO player_info (player_id, player_season_key, player_name, player_status,
                         player_status_full, player_url, player_team, player_bye_week, player_uniform_number,
                         player_position, player_headshot, player_injury_notes)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT DO NOTHING`

	insertedCount := 0
	skippedCount := 0
	for counter := 0; counter < len(playerInfo); counter++ {
		playerStatusFull := playerInfo[counter].PlayerStatusFull
		if strings.Contains(playerStatusFull, inactiveStatusFull) {
			skippedCount++
			continue
		}

		playerId := playerInfo[counter].PlayerID
		playerKey := playerInfo[counter].PlayerSeasonId
		playerName := playerInfo[counter].PlayerName
		playerStatus := playerInfo[counter].PlayerStatus
		playerPosition := playerInfo[counter].PlayerPosition
		playerTeam := playerInfo[counter].PlayerTeam
		playerByeWeek := playerInfo[counter].PlayerByeWeek
		playerUniformNumber := playerInfo[counter].PlayerUniformNumber
		playerUrl := playerInfo[counter].PlayerUrl
		playerHeadshot := playerInfo[counter].PlayerHeadshot
		playerInjuryNotes := playerInfo[counter].PlayerInjuryNotes

		count, err := ExecStatement(sqlStatement, playerId, playerKey, playerName, playerStatus, playerStatusFull, playerUrl,
			playerTeam, playerByeWeek, playerUniformNumber, playerPosition, playerHeadshot, playerInjuryNotes)
		if err != nil {
			logger.LogError(context.Background(), "Error inserting available player record", "error", err.Error(), "player_id", playerId)
			continue
		}
		insertedCount += int(count)
		logger.LogInfo(context.Background(), "Rows affected", "count", strconv.Itoa(int(count)))
	}
	logger.LogInfo(context.Background(), "Done inserting available player records", "total_records", strconv.Itoa(len(playerInfo)),
		"inserted", strconv.Itoa(insertedCount), "skipped_inactive", strconv.Itoa(skippedCount))
	return insertedCount
}
