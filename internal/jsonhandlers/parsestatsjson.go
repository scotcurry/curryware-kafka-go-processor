package jsonhandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"encoding/json"
	"fmt"
)

func ParseMultipleStatInfo(statsInfo string) []fantasyclasses.StatsInfo {

	var multipleStatInfo fantasyclasses.StatsJson
	err := json.Unmarshal([]byte(statsInfo), &multipleStatInfo)
	if err != nil {
		fmt.Println("Error parsing player info")
	}
	statsArray := multipleStatInfo.PlayerStats
	return statsArray
}
