package playertests

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/playerclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"encoding/base64"
	"testing"
)

func TestParseSinglePlayer(t *testing.T) {

	singlePlayer := `{"Key":"461.p.28556","Id":28556,"FullName":"Michael Burton","Url":"https://sports.yahoo.com/nfl/players/28556","Status":"Q","Team":"Den","ByeWeek":12,"UniformNumber":20,"Position":"RB","Headshot":"https://s.yimg.com/iu/api/res/1.2/2ysfDyrGSW9R1igZmyqqow--~C/YXBwaWQ9eXNwb3J0cztjaD0yMzM2O2NyPTE7Y3c9MTc5MDtkeD04NTc7ZHk9MDtmaT11bGNyb3A7aD02MDtxPTEwMDt3PTQ2/https://s.yimg.com/xe/i/us/sp/v/nfl_cutout/players_l/08122025/28556.png","InjuryNotes":"NA"}`
	encodedPlayer := base64.StdEncoding.EncodeToString([]byte(singlePlayer))

	result, err := jsonhandlers.ParseJSON[playerclasses.PlayerInfo](encodedPlayer)
	if err != nil {
		t.Errorf("Error parsing JSON: %v", err)
	}

	if result.PlayerID != 28556 {
		t.Errorf("Expected PlayerID 28556, got %d", result.PlayerID)
	}
}

func TestParseMultiplePlayers(t *testing.T) {

	multiPlayerJson := `[{"Key":"461.p.29269","Id":29269,"FullName":"Hunter Henry","Url":"https://sports.yahoo.com/nfl/players/29269","Status":"Active","PlayerStatusFull":"Active","Team":"NE","ByeWeek":14,"UniformNumber":85,"Position":"TE","Headshot":"https://s.yimg.com/iu/api/res/1.2/1TaZL9zISaBQw.W1qxfzFw--~C/YXBwaWQ9eXNwb3J0cztjaD0yMzM2O2NyPTE7Y3c9MTc5MDtkeD04NTc7ZHk9MDtmaT11bGNyb3A7aD02MDtxPTEwMDt3PTQ2/https://s.yimg.com/xe/i/us/sp/v/nfl_cutout/players_l/08292024/29269.png","InjuryNotes":"NA"},{"Key":"461.p.29274","Id":29274,"FullName":"Sterling Shepard","Url":"https://sports.yahoo.com/nfl/players/29274","Status":"Active","PlayerStatusFull":"Active","Team":"TB","ByeWeek":9,"UniformNumber":17,"Position":"WR","Headshot":"https://s.yimg.com/iu/api/res/1.2/6EWNAkZk4_O0nmpjwQp7Iw--~C/YXBwaWQ9eXNwb3J0cztjaD0yMzM2O2NyPTE7Y3c9MTc5MDtkeD04NTc7ZHk9MDtmaT11bGNyb3A7aD02MDtxPTEwMDt3PTQ2/https://s.yimg.com/xe/i/us/sp/v/nfl_cutout/players_l/08222023/29274.1.png","InjuryNotes":"NA"}]`
	encodedPlayers := base64.StdEncoding.EncodeToString([]byte(multiPlayerJson))

	result, err := jsonhandlers.ParseJSON[[]playerclasses.PlayerInfo](encodedPlayers)
	if err != nil {
		t.Errorf("Error parsing JSON: %v", err)
	}

	if result[0].PlayerID != 29269 {
		t.Errorf("Expected first PlayerID 29269, got %d", result[0].PlayerID)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 players, got %d", len(result))
	}
}
