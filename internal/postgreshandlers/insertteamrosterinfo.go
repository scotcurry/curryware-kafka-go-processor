package postgreshandlers

import (
	"context"
	"curryware-kafka-go-processor/internal/fantasyclasses/teamclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"strconv"
)

func InsertTeamRosterInfo(rosterInfo []teamclasses.TeamRosterInfo) int {
	sqlStatement := `INSERT INTO team_roster_information (team_key, team_id, manager_id, player_key, player_id,
                         team_full_name, bye_week, primary_position, has_player_notes, last_player_note_timestamp)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                ON CONFLICT (team_key, player_key) DO UPDATE SET
                    team_id = EXCLUDED.team_id,
                    manager_id = EXCLUDED.manager_id,
                    player_id = EXCLUDED.player_id,
                    team_full_name = EXCLUDED.team_full_name,
                    bye_week = EXCLUDED.bye_week,
                    primary_position = EXCLUDED.primary_position,
                    has_player_notes = EXCLUDED.has_player_notes,
                    last_player_note_timestamp = EXCLUDED.last_player_note_timestamp`

	for counter := range rosterInfo {
		teamKey := rosterInfo[counter].TeamKey
		teamId := rosterInfo[counter].TeamId
		managerId := rosterInfo[counter].ManagerId
		playerKey := rosterInfo[counter].PlayerKey
		playerId := rosterInfo[counter].PlayerId
		teamFullName := rosterInfo[counter].TeamFullName
		byeWeek := rosterInfo[counter].ByeWeek
		primaryPosition := rosterInfo[counter].PrimaryPosition
		hasPlayerNotes := rosterInfo[counter].HasPlayerNotes
		lastPlayerNoteTimestamp := rosterInfo[counter].LastPlayerNoteTimestamp

		count, err := ExecStatement(sqlStatement, teamKey, teamId, managerId, playerKey, playerId,
			teamFullName, byeWeek, primaryPosition, hasPlayerNotes, lastPlayerNoteTimestamp)
		if err != nil {
			logger.LogError(context.Background(), "Error inserting team roster record", "error", err.Error(), "team_key", teamKey, "player_key", playerKey)
			continue
		}
		logger.LogInfo(context.Background(), "Rows affected", "count", strconv.Itoa(int(count)))
	}
	logger.LogInfo(context.Background(), "Done inserting team roster records", "total_records", strconv.Itoa(len(rosterInfo)))
	return len(rosterInfo)
}
