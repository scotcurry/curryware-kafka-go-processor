package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"fmt"
	"testing"
)

func TestInsertPlayerRecord(t *testing.T) {

	singlePlayerRecord := `{"Key":"449.p.40300","Id":40300,"FullName":"Cephus Johnson III","Url":"https://sports.yahoo.com/nfl/players/40300","Status":"NA","Team":"TB","ByeWeek":0,"UniformNumber":0,"Position":"WR","Headshot":"https://s.yimg.com/iu/api/res/1.2/8YBaY7mta1OqebkDNJ9JUg--~C/YXBwaWQ9eXNwb3J0cztjaD0yMzM2O2NyPTE7Y3c9MTc5MDtkeD04NTc7ZHk9MDtmaT11bGNyb3A7aD02MDtxPTEwMDt3PTQ2/https://s.yimg.com/xe/i/us/sp/v/nfl_cutout/players_l/08092023/40300.png"}`
	var playerRecord = jsonhandlers.ParsePlayerInfo(singlePlayerRecord)

	playerArray := []fantasyclasses.PlayerInfo{playerRecord}
	InsertPlayerRecord(playerArray)
	fmt.Println("Player record inserted")
	t.Log("Test passed")
}
