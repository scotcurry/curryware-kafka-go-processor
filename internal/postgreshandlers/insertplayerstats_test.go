package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/statsclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestInsertPlayerStatsRecord(t *testing.T) {

	allPlayersJson := `[{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":4,"statValue":295},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":5,"statValue":3}]`
	encodedStats := base64.StdEncoding.EncodeToString([]byte(allPlayersJson))
	allPlayersArray, err := jsonhandlers.ParseJSON[[]statsclasses.PlayerWeeklyStatsInfo](encodedStats)
	if err != nil {
		t.Fatalf("Error parsing player stats: %v", err)
	}
	fmt.Println(allPlayersArray)
	InsertPlayerStats(allPlayersArray)
}
