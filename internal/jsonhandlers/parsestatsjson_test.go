package jsonhandlers

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestParseMultipleStatInfo(t *testing.T) {

	multiStatJson := `[{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":4,"statValue":295},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":5,"statValue":3},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":6,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":8,"statValue":7},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":9,"statValue":44},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":10,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":78,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":11,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":12,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":13,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":15,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":16,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":18,"statValue":1},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":15,"statId":57,"statValue":0},{"playerId":33389,"playerGameKey":"461.p.33389","playerStatWeek":15,"statId":4,"statValue":330}]`

	encodedStats := base64.StdEncoding.EncodeToString([]byte(multiStatJson))
	multiStatTestResult := ParseMultipleStatInfo(encodedStats)

	if len(multiStatTestResult) == 0 {
		t.Fatal("Expected stats, got empty slice")
	}

	if multiStatTestResult[0].PlayerID != 34218 {
		t.Errorf("Expected 34218, got %d", multiStatTestResult[0].PlayerID)
	}
	fmt.Println(multiStatTestResult)
}
