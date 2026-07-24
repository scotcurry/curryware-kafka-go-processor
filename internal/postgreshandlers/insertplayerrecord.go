package postgreshandlers

import (
	"context"
	"curryware-kafka-go-processor/internal/fantasyclasses/playerclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"strconv"
)

func InsertPlayerRecord(ctx context.Context, playerInfo []playerclasses.PlayerInfo) {
	sqlStatement := `INSERT INTO player_info_snapshot (player_id, player_season_key, player_name, player_status,
                         player_status_full, player_url,  player_team, player_bye_week, player_uniform_number,
                         player_position, player_headshot, player_injury_notes)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT DO NOTHING`
	for counter := 0; counter < len(playerInfo); counter++ {
		playerId := playerInfo[counter].PlayerID
		playerKey := playerInfo[counter].PlayerSeasonId
		playerName := playerInfo[counter].PlayerName
		playerStatus := playerInfo[counter].PlayerStatus
		playerStatusFull := playerInfo[counter].PlayerStatusFull
		playerPosition := playerInfo[counter].PlayerPosition
		playerTeam := playerInfo[counter].PlayerTeam
		playerByeWeek := playerInfo[counter].PlayerByeWeek
		playerUniformNumber := playerInfo[counter].PlayerUniformNumber
		playerUrl := playerInfo[counter].PlayerUrl
		playerHeadshot := playerInfo[counter].PlayerHeadshot
		playerInjuryNotes := playerInfo[counter].PlayerInjuryNotes

		count, err := ExecStatement(ctx, sqlStatement, playerId, playerKey, playerName, playerStatus, playerStatusFull, playerUrl,
			playerTeam, playerByeWeek, playerUniformNumber, playerPosition, playerHeadshot, playerInjuryNotes)
		if err != nil {
			logger.LogError(ctx, "Error inserting player record", "error", err.Error(), "player_id", playerId)
			continue
		}
		logger.LogInfo(ctx, "Rows affected", "count", strconv.Itoa(int(count)))
	}
	logger.LogInfo(ctx, "Done inserting player records", "total_records", strconv.Itoa(len(playerInfo)))
}
