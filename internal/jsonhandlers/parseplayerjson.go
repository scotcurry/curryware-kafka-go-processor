package jsonhandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"encoding/json"
	"fmt"
)

func ParsePlayerInfo(playerInfo string) fantasyclasses.PlayerInfo {

	var playerInfoStruct fantasyclasses.PlayerInfo
	err := json.Unmarshal([]byte(playerInfo), &playerInfoStruct)
	if err != nil {
		fmt.Println("Error parsing player info")
	}

	return playerInfoStruct
}

func ParseMultiplePlayerInfo(playerInfo string) []fantasyclasses.PlayerInfo {

	var multiPlayerInfoStruct []fantasyclasses.PlayerInfo
	err := json.Unmarshal([]byte(playerInfo), &multiPlayerInfoStruct)
	if err != nil {
		fmt.Println("Error parsing player info")
	}

	return multiPlayerInfoStruct
}
