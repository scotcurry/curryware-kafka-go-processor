package jsonhandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func ParsePlayerInfo(playerInfo string) fantasyclasses.PlayerInfo {

	decodedBytes, err := base64.StdEncoding.DecodeString(playerInfo)
	playerInfo = string(decodedBytes)

	var playerInfoStruct fantasyclasses.PlayerInfo
	err = json.Unmarshal([]byte(playerInfo), &playerInfoStruct)
	if err != nil {
		logger.LogError("Error parsing player info: ", err)
	}

	return playerInfoStruct
}

func ParseMultiplePlayerInfo(playerInfo string) []fantasyclasses.PlayerInfo {

	decodedBytes, err := base64.StdEncoding.DecodeString(playerInfo)
	playerInfo = string(decodedBytes)

	var multiPlayerInfoStruct []fantasyclasses.PlayerInfo
	err = json.Unmarshal([]byte(playerInfo), &multiPlayerInfoStruct)
	if err != nil {
		fmt.Println("Error parsing player info")
	}

	logger.LogInfo("Number of players parsed: ", len(multiPlayerInfoStruct))
	return multiPlayerInfoStruct
}
