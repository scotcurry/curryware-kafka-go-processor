package leaguestatstests

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"encoding/base64"
	"testing"
)

// The exact shape the Yahoo API's LeagueStatValueParser serializes (System.Text.Json, decimal
// StatValue, nullable).
const leagueStatValuesJson = `[{"LeagueStatKey":4708388594,"GameId":470,"LeagueId":838859,"StatId":4,"StatValue":0.04},` +
	`{"LeagueStatKey":4708388596,"GameId":470,"LeagueId":838859,"StatId":6,"StatValue":-1},` +
	`{"LeagueStatKey":47083885957,"GameId":470,"LeagueId":838859,"StatId":57,"StatValue":null}]`

func TestParseLeagueStatValues(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte(leagueStatValuesJson))

	result, err := jsonhandlers.ParseJSON[[]leagueclasses.LeagueStatsValueInfo](encoded)
	if err != nil {
		t.Fatalf("Error parsing JSON: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("Expected 3 league stat values, got %d", len(result))
	}
	first := result[0]
	if first.GameId != 470 || first.LeagueId != 838859 || first.StatId != 4 {
		t.Errorf("Unexpected ids: %+v", first)
	}
	if first.LeagueStatKey != 4708388594 {
		t.Errorf("Expected LeagueStatKey 4708388594, got %d", first.LeagueStatKey)
	}
	if first.StatValue == nil || *first.StatValue != 0.04 {
		t.Errorf("Expected StatValue 0.04, got %v", first.StatValue)
	}
	if result[1].StatValue == nil || *result[1].StatValue != -1 {
		t.Errorf("Expected StatValue -1, got %v", result[1].StatValue)
	}
	if result[2].StatValue != nil {
		t.Errorf("Expected a null StatValue, got %v", *result[2].StatValue)
	}
}
