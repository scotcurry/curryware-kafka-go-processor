package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/logging"
	"strconv"
)

func InsertPlayerRecord(playerInfo []fantasyclasses.PlayerInfo) {
	// Use the singleton database connection pool
	db := GetDB()
	sqlStatement := `INSERT INTO player_info (player_id, player_season_key, player_name, player_status, 
                         player_status_full, player_url,  player_team, player_bye_week, player_uniform_number, 
                         player_position, player_headshot, player_injury_notes) 
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
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

		res, err := db.Exec(sqlStatement, playerId, playerKey, playerName, playerStatus, playerStatusFull, playerUrl,
			playerTeam, playerByeWeek, playerUniformNumber, playerPosition, playerHeadshot, playerInjuryNotes)
		if err != nil {
			logging.LogError("Error inserting player record", "error", err.Error(), "player_id", playerId)
			panic(err)
		}

		count, err := res.RowsAffected()
		if err != nil {
			logging.LogError("Error getting rows affected", "error", err.Error())
			panic(err)
		}
		logging.LogInfo("Rows affected", "count", strconv.Itoa(int(count)))
	}
	logging.LogInfo("Done inserting player records", "total_records", strconv.Itoa(len(playerInfo)))
}
