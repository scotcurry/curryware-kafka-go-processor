package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/logging"
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

func InsertPlayerRecord(playerInfo []fantasyclasses.PlayerInfo) {

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
	logging.LogInfo(fmt.Sprintf("Done inserting player records. Total records: %s", strconv.Itoa(len(playerInfo))))
}
