package teamrostertests

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/teamclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"encoding/base64"
	"testing"
)

func TestParseTeamRosterInfo(t *testing.T) {

	teamRosterJson := `[{"TeamKey":"461.l.460188.t.5","TeamId":5,"ManagerId":5,"PlayerKey":"461.p.9265","PlayerId":9265,"TeamFullName":"Los Angeles Rams","ByeWeek":8,"PrimaryPosition":"QB","HasPlayerNotes":false,"LastPlayerNoteTimestamp":null},{"TeamKey":"461.l.460188.t.5","TeamId":5,"ManagerId":5,"PlayerKey":"461.p.31010","PlayerId":31010,"TeamFullName":"Denver Broncos","ByeWeek":12,"PrimaryPosition":"WR","HasPlayerNotes":true,"LastPlayerNoteTimestamp":1783725756}]`
	encodedTeamRoster := base64.StdEncoding.EncodeToString([]byte(teamRosterJson))

	result, err := jsonhandlers.ParseJSON[[]teamclasses.TeamRosterInfo](encodedTeamRoster)
	if err != nil {
		t.Errorf("Error parsing JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 roster entries, got %d", len(result))
	}

	if result[0].TeamKey != "461.l.460188.t.5" {
		t.Errorf("Expected TeamKey 461.l.460188.t.5, got %s", result[0].TeamKey)
	}

	if result[0].PlayerId != 9265 {
		t.Errorf("Expected PlayerId 9265, got %d", result[0].PlayerId)
	}

	if result[0].LastPlayerNoteTimestamp != nil {
		t.Errorf("Expected nil LastPlayerNoteTimestamp, got %v", *result[0].LastPlayerNoteTimestamp)
	}

	if result[1].LastPlayerNoteTimestamp == nil || *result[1].LastPlayerNoteTimestamp != 1783725756 {
		t.Errorf("Expected LastPlayerNoteTimestamp 1783725756, got %v", result[1].LastPlayerNoteTimestamp)
	}

	if !result[1].HasPlayerNotes {
		t.Errorf("Expected HasPlayerNotes true for second entry")
	}
}
