package jsonhandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/logging"
	"encoding/base64"
	"encoding/json"
)

func ParseLeagueStatValue(statsInfo string) []fantasyclasses.LeagueStatsValueInfo {

	decodedBytes, err := base64.StdEncoding.DecodeString(statsInfo)
	statsInfo = string(decodedBytes)

	var statValueStruct []fantasyclasses.LeagueStatsValueInfo
	err = json.Unmarshal([]byte(statsInfo), &statValueStruct)
	if err != nil {
		logging.LogError("Error Parsing Stat Values JSON ", err)
	}

	return statValueStruct
}
