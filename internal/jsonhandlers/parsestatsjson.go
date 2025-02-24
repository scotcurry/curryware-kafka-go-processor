package jsonhandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/logging"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func ParseMultipleStatInfo(statsInfo string) []fantasyclasses.StatsInfo {

	decodedBytes, err := base64.StdEncoding.DecodeString(statsInfo)
	statsInfo = string(decodedBytes)

	var multipleStatInfo fantasyclasses.StatsJson
	err = json.Unmarshal([]byte(statsInfo), &multipleStatInfo)
	if err != nil {
		fmt.Println("Error parsing player info")
		logging.LogError("Error parsing player info: ", err)
	}
	statsArray := multipleStatInfo.PlayerStats
	return statsArray
}
