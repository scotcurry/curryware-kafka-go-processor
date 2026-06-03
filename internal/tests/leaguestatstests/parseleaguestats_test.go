package leaguestatstests

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"encoding/base64"
	"testing"
)

func TestParseLeagueStatInfo(t *testing.T) {

	leagueStatInfo := `[{"LeagueStatKeyId":"4493606174_4","GameId":1,"LeagueId":4493606,"StatId":4,"StatEnabled":true,"StatName":"Passing Yards","StatDisplayName":"Pass Yds","StatGroupDisplayName":"passing","StatAbbreviation":"Yds","StatSortOrder":1,"StatPositionType":"O","StatSortPosition":0},{"LeagueStatKeyId":"4493606175_5","GameId":1,"LeagueId":4493606,"StatId":5,"StatEnabled":true,"StatName":"Passing Touchdowns","StatDisplayName":"Pass TD","StatGroupDisplayName":"passing","StatAbbreviation":"TD","StatSortOrder":1,"StatPositionType":"O","StatSortPosition":0},{"LeagueStatKeyId":"4493606176_6","GameId":1,"LeagueId":4493606,"StatId":6,"StatEnabled":true,"StatName":"Interceptions","StatDisplayName":"Int","StatGroupDisplayName":"passing","StatAbbreviation":"Int","StatSortOrder":0,"StatPositionType":"O","StatSortPosition":0}]`
	encodedStats := base64.StdEncoding.EncodeToString([]byte(leagueStatInfo))

	result, err := jsonhandlers.ParseJSON[[]leagueclasses.LeagueStatDescriptionInfo](encodedStats)
	if err != nil {
		t.Errorf("Error parsing JSON: %v", err)
	}

	if result[0].StatId != 4 {
		t.Errorf("Expected first StatId 4, got %d", result[0].StatId)
	}

	if result[0].LeagueStatKeyId != "4493606174_4" {
		t.Errorf("Expected LeagueStatKeyId '4493606174_4', got %s", result[0].LeagueStatKeyId)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 league stats, got %d", len(result))
	}
}
