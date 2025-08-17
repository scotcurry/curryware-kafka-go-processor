package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestInsertPlayerRecord(t *testing.T) {

	singlePlayerRecord := `{"Key":"461.p.29269","Id":29269,"FullName":"Hunter Henry","Url":"https://sports.yahoo.com/nfl/players/29269","Status":"Active","PlayerStatusFull":"Active","Team":"NE","ByeWeek":14,"UniformNumber":85,"Position":"TE","Headshot":"https://s.yimg.com/iu/api/res/1.2/1TaZL9zISaBQw.W1qxfzFw--~C/YXBwaWQ9eXNwb3J0cztjaD0yMzM2O2NyPTE7Y3c9MTc5MDtkeD04NTc7ZHk9MDtmaT11bGNyb3A7aD02MDtxPTEwMDt3PTQ2/https://s.yimg.com/xe/i/us/sp/v/nfl_cutout/players_l/08292024/29269.png","InjuryNotes":"NA"}`
	encodedPlayerRecord := base64.StdEncoding.EncodeToString([]byte(singlePlayerRecord))
	var playerRecord = jsonhandlers.ParsePlayerInfo(encodedPlayerRecord)

	playerArray := []fantasyclasses.PlayerInfo{playerRecord}
	InsertPlayerRecord(playerArray)
	fmt.Println("Player record inserted")
	t.Log("Test passed")
}
