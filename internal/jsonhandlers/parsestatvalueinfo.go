package jsonhandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/logging"
	"encoding/base64"
	"encoding/json"
)

func ParseLeagueStatValue(statsInfo string) ([]fantasyclasses.LeagueStatsValueInfo, error) {

	decodedBytes, err := base64.StdEncoding.DecodeString(statsInfo)

	if err != nil {
		logging.LogError("Error Parsing Stat Values JSON ", err)
		return nil, err
	}
	var statValueStruct []fantasyclasses.LeagueStatsValueInfo
	err = json.Unmarshal(decodedBytes, &statValueStruct)

	return statValueStruct, nil
}
