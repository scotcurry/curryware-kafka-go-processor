package jsonhandlers

import (
	"fmt"
	"testing"
)

func TestParseSinglePlayerInfo(t *testing.T) {

	singlePlayer := `{"Key":"449.p.33500","Id":33500,"FullName":"Amon-Ra St. Brown","Url":"https://sports.yahoo.com/nfl/players/33500","Status":"","Team":"Det","ByeWeek":0,"UniformNumber":14,"Position":"WR","Headshot":"https://s.yimg.com/iu/api/res/1.2/04AG_FTYoS61qWZEVJDaSg--~C/YXBwaWQ9eXNwb3J0cztjaD0yMzM2O2NyPTE7Y3c9MTc5MDtkeD04NTc7ZHk9MDtmaT11bGNyb3A7aD02MDtxPTEwMDt3PTQ2/https://s.yimg.com/xe/i/us/sp/v/nfl_cutout/players_l/08142024/33500.png"}`
	singlePlayerTestResult := ParsePlayerInfo(singlePlayer)
	if singlePlayerTestResult.PlayerID != 33500 {
		t.Errorf("Expected 33500, got %d", singlePlayerTestResult.PlayerID)
	}
	fmt.Println(singlePlayerTestResult)
	t.Log("Test passed")
}

func TestParseMultiplePlayerInfo(t *testing.T) {

	multiPlayerJson := `[
{"Key":"449.p.40300","Id":40300,"FullName":"Cephus Johnson III","Url":"https://sports.yahoo.com/nfl/players/40300","Status":"NA","Team":"TB","ByeWeek":0,"UniformNumber":0,"Position":"WR","Headshot":"https://s.yimg.com/iu/api/res/1.2/8YBaY7mta1OqebkDNJ9JUg--~C/YXBwaWQ9eXNwb3J0cztjaD0yMzM2O2NyPTE7Y3c9MTc5MDtkeD04NTc7ZHk9MDtmaT11bGNyb3A7aD02MDtxPTEwMDt3PTQ2/https://s.yimg.com/xe/i/us/sp/v/nfl_cutout/players_l/08092023/40300.png"},
{"Key":"449.p.33500","Id":33500,"FullName":"Amon-Ra St. Brown","Url":"https://sports.yahoo.com/nfl/players/33500","Status":"","Team":"Det","ByeWeek":0,"UniformNumber":14,"Position":"WR","Headshot":"https://s.yimg.com/iu/api/res/1.2/04AG_FTYoS61qWZEVJDaSg--~C/YXBwaWQ9eXNwb3J0cztjaD0yMzM2O2NyPTE7Y3c9MTc5MDtkeD04NTc7ZHk9MDtmaT11bGNyb3A7aD02MDtxPTEwMDt3PTQ2/https://s.yimg.com/xe/i/us/sp/v/nfl_cutout/players_l/08142024/33500.png"}
]`

	multiPlayerTestResult := ParseMultiplePlayerInfo(multiPlayerJson)
	firstPlayer := multiPlayerTestResult[0]
	if firstPlayer.PlayerID != 40300 {
		t.Errorf("Expected 40300, got %d", firstPlayer.PlayerID)
	}
	fmt.Println(multiPlayerTestResult[1].PlayerName)
	t.Log(firstPlayer.PlayerName)
}
