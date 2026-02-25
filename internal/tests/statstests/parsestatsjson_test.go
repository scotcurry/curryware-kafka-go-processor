package statstests

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/statsclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestParsePlayerWeeklyStats(t *testing.T) {

	multiStatJson := `[{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":4,"statValue":295},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":5,"statValue":3},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":6,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":8,"statValue":7},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":9,"statValue":44},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":10,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":78,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":11,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":12,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":13,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":15,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":16,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":18,"statValue":1},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":57,"statValue":0},{"playerId":33389,"playerGameKey":"461.p.33389","playerStatWeek":15,"statId":4,"statValue":330}]`

	encodedStats := base64.StdEncoding.EncodeToString([]byte(multiStatJson))
	// This is how you use a generic type parameter to parse JSON into a slice of PlayerWeeklyStatsInfo.
	statTestResult, err := jsonhandlers.ParseJSON[[]statsclasses.PlayerWeeklyStatsInfo](encodedStats)

	if err != nil {
		t.Errorf("Error parsing JSON: %v", err)
	}

	if statTestResult[0].PlayerId != 34218 {
		t.Errorf("Expected playerId 34218, got %d", statTestResult[0].PlayerId)
	}

	if statTestResult[0].StatValue != 295 {
		t.Errorf("Expected statValue 295, got %f", statTestResult[0].StatValue)
	}
	fmt.Println(statTestResult)
}
