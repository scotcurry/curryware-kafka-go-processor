package jsonhandlers

import (
	"encoding/base64"
	"testing"
)

func TestParsePlayerStatValueInfo(t *testing.T) {

	playerStatValueJson := `[{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":16,"statId":4,"statValue":295},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":16,"statId":5,"statValue":5},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":16,"statId":6,"statValue":1},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":16,"statId":8,"statValue":2},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":16,"statId":9,"statValue":11},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":16,"statId":10,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":16,"statId":78,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":16,"statId":11,"statValue":0},{"playerId":34218,"playerGameKey":"461.p.34218","playerStatWeek":16,"statId":12,"statValue":0}]`
	playerStatValueBase64 := base64.StdEncoding.EncodeToString([]byte(playerStatValueJson))

	playerStatValueClass, err := ParseLeagueStatValues(playerStatValueBase64)
	if err != nil {
		t.Fatal(err)
	}
	if playerStatValueClass[0].PlayerId != 34218 && playerStatValueClass[0].StatValue != 295 {
		t.Errorf("Expected 295, got %.2f", playerStatValueClass[0].StatValue)
	}
}
