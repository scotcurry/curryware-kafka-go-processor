package integration_tests

import (
	"context"
	"curryware-kafka-go-processor/internal/fantasyclasses/teamclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"curryware-kafka-go-processor/internal/postgreshandlers"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestInsertTeamWeeklyStandingsInfo(t *testing.T) {

	teamWeeklyStandingsJson := `[{"TeamKey":"461.l.460188.t.1","TeamId":1,"WeekKey":17,"Rank":1,"PointsFor":1652.50,"PointsAgainst":1468.76,"Wins":9,"Losses":5,"Ties":0,"NumberOfMoves":33},{"TeamKey":"461.l.460188.t.3","TeamId":3,"WeekKey":17,"Rank":2,"PointsFor":1702.92,"PointsAgainst":1409.9,"Wins":10,"Losses":4,"Ties":0,"NumberOfMoves":20}]`
	encodedTeamWeeklyStandings := base64.StdEncoding.EncodeToString([]byte(teamWeeklyStandingsJson))
	teamWeeklyStandingsInfo, err := jsonhandlers.ParseJSON[[]teamclasses.TeamWeeklyStandingsInfo](encodedTeamWeeklyStandings)
	if err != nil {
		t.Fatalf("Error parsing team weekly standings info: %v", err)
	}

	postgreshandlers.InsertTeamWeeklyStandingsInfo(context.Background(), teamWeeklyStandingsInfo)
	fmt.Println("Team weekly standings records inserted")
	t.Log("Test passed")
}
