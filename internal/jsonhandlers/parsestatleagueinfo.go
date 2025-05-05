package jsonhandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"encoding/base64"
	"encoding/json"
)

func ParsLeagueStatInfo(leaguestatisticinfo string) []fantasyclasses.LeagueStatInfo {

	decodedBytes, err := base64.StdEncoding.DecodeString(leaguestatisticinfo)
	leaguestatisticinfo = string(decodedBytes)

	var leagueinfo []fantasyclasses.LeagueStatInfo
	err = json.Unmarshal([]byte(leaguestatisticinfo), &leagueinfo)
	if err != nil {
		logger.LogError("Error parsing player info: ", err)
	}

	return leagueinfo
}
