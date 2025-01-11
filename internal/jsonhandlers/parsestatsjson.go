package jsonhandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/logging"
	"encoding/json"
	"fmt"
)

func ParseMultipleStatInfo(statsInfo string) []fantasyclasses.StatsInfo {

	var multipleStatInfo fantasyclasses.StatsJson
	err := json.Unmarshal([]byte(statsInfo), &multipleStatInfo)
	if err != nil {
		fmt.Println("Error parsing player info")
		logging.LogError("Error parsing player info: ", err)
	}
	statsArray := multipleStatInfo.PlayerStats
	return statsArray
}
