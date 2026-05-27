package leagueclasses

type LeagueStatDescriptionInfo struct {
	LeagueStatKeyId      string `json:"LeagueStatKeyId"`
	GameId               int    `json:"GameId"`
	LeagueId             int    `json:"LeagueId"`
	StatId               int    `json:"StatId"`
	StatEnabled          bool   `json:"StatEnabled"`
	StatName             string `json:"StatName"`
	StatDisplayName      string `json:"StatDisplayName"`
	StatGroupDisplayName string `json:"StatGroupDisplayName"`
	StatAbbreviation     string `json:"StatAbbreviation"`
	StatSortOrder        int    `json:"StatSortOrder"`
	StatPositionType     string `json:"StatPositionType"`
	StatSortPosition     int    `json:"StatSortPosition"`
}
