package jsonhandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"encoding/base64"
	"encoding/json"
)

func ParseLeagueStatInfo(leagueStatsInfo string) ([]fantasyclasses.LeagueStatInfo, error) {

	decodedBytes, err := base64.StdEncoding.DecodeString(leagueStatsInfo)
	if err != nil {
		logger.LogError("Error parsing player info: ", err)
		return nil, err
	} else {
		logger.LogInfo("Decoded bytes: ", decodedBytes[:20])
	}
	leagueStatsInfo = string(decodedBytes)

	var statInfos []fantasyclasses.LeagueStatInfo
	err = json.Unmarshal([]byte(leagueStatsInfo), &statInfos)
	if err != nil {
		logger.LogError("Error parsing player info: ", err)
		return nil, err
	}

	return statInfos, nil
}
