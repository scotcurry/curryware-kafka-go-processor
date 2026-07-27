package teamweeklystandingstests

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/teamclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"encoding/base64"
	"testing"
)

func TestParseTeamWeeklyStandingsInfo(t *testing.T) {

	teamWeeklyStandingsJson := `[{"TeamKey":"461.l.460188.t.1","TeamId":1,"WeekKey":17,"Rank":1,"PointsFor":1652.50,"PointsAgainst":1468.76,"Wins":9,"Losses":5,"Ties":0,"NumberOfMoves":33},{"TeamKey":"461.l.460188.t.3","TeamId":3,"WeekKey":17,"Rank":2,"PointsFor":1702.92,"PointsAgainst":1409.9,"Wins":10,"Losses":4,"Ties":0,"NumberOfMoves":20}]`
	encodedTeamWeeklyStandings := base64.StdEncoding.EncodeToString([]byte(teamWeeklyStandingsJson))

	result, err := jsonhandlers.ParseJSON[[]teamclasses.TeamWeeklyStandingsInfo](encodedTeamWeeklyStandings)
	if err != nil {
		t.Errorf("Error parsing JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 standings entries, got %d", len(result))
	}

	if result[0].TeamKey != "461.l.460188.t.1" {
		t.Errorf("Expected TeamKey 461.l.460188.t.1, got %s", result[0].TeamKey)
	}

	if result[0].WeekKey != 17 {
		t.Errorf("Expected WeekKey 17, got %d", result[0].WeekKey)
	}

	if result[0].Rank != 1 {
		t.Errorf("Expected Rank 1, got %d", result[0].Rank)
	}

	if result[0].PointsFor != 1652.50 {
		t.Errorf("Expected PointsFor 1652.50, got %v", result[0].PointsFor)
	}

	if result[1].Wins != 10 || result[1].Losses != 4 {
		t.Errorf("Expected Wins 10 and Losses 4 for second entry, got Wins %d Losses %d", result[1].Wins, result[1].Losses)
	}

	if result[1].NumberOfMoves != 20 {
		t.Errorf("Expected NumberOfMoves 20, got %d", result[1].NumberOfMoves)
	}
}
