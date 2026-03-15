package leaguestatstests

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/leagueclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"encoding/base64"
	"testing"
)

func TestParseLeagueStatInfo(t *testing.T) {

	leagueStatInfo := `[{"LeagueStatKey":4493606174,"StatId":4,"StatEnabled":true,"StatName":"Passing Yards","StatDisplayName":"Pass Yds","StatGroup":"passing","StatAbbreviation":"Yds","StatSortOrder":1,"StatPositionType":"O","StatSortPostion":0},{"LeagueStatKey":4493606175,"StatId":5,"StatEnabled":true,"StatName":"Passing Touchdowns","StatDisplayName":"Pass TD","StatGroup":"passing","StatAbbreviation":"TD","StatSortOrder":1,"StatPositionType":"O","StatSortPostion":0},{"LeagueStatKey":4493606176,"StatId":6,"StatEnabled":true,"StatName":"Interceptions","StatDisplayName":"Int","StatGroup":"passing","StatAbbreviation":"Int","StatSortOrder":0,"StatPositionType":"O","StatSortPostion":0}]`
	encodedStats := base64.StdEncoding.EncodeToString([]byte(leagueStatInfo))

	result, err := jsonhandlers.ParseJSON[[]leagueclasses.LeagueStatDescriptionInfo](encodedStats)
	if err != nil {
		t.Errorf("Error parsing JSON: %v", err)
	}

	if result[0].StatId != 4 {
		t.Errorf("Expected first StatId 4, got %d", result[0].StatId)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 league stats, got %d", len(result))
	}
}
