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

func TestInsertTeamRosterInfo(t *testing.T) {

	teamRosterJson := `[{"TeamKey":"461.l.460188.t.5","TeamId":5,"ManagerId":5,"PlayerKey":"461.p.9265","PlayerId":9265,"TeamFullName":"Los Angeles Rams","ByeWeek":8,"PrimaryPosition":"QB","HasPlayerNotes":false,"LastPlayerNoteTimestamp":null},{"TeamKey":"461.l.460188.t.5","TeamId":5,"ManagerId":5,"PlayerKey":"461.p.31010","PlayerId":31010,"TeamFullName":"Denver Broncos","ByeWeek":12,"PrimaryPosition":"WR","HasPlayerNotes":true,"LastPlayerNoteTimestamp":1783725756}]`
	encodedTeamRoster := base64.StdEncoding.EncodeToString([]byte(teamRosterJson))
	teamRosterInfo, err := jsonhandlers.ParseJSON[[]teamclasses.TeamRosterInfo](encodedTeamRoster)
	if err != nil {
		t.Fatalf("Error parsing team roster info: %v", err)
	}

	postgreshandlers.InsertTeamRosterInfo(context.Background(), teamRosterInfo)
	fmt.Println("Team roster records inserted")
	t.Log("Test passed")
}
