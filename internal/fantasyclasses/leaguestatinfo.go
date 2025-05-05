package fantasyclasses

type LeagueStatInfo struct {
	LeagueStatId     int    `json:"LeagueStatKey"`
	StatId           int    `json:"StatId"`
	StatEnabled      bool   `json:"StatEnabled"`
	StatName         string `json:"StatName"`
	StatDisplayName  string `json:"StatDisplayName"`
	StatGroup        string `json:"StatGroup"`
	StatAbbreviation string `json:"StatAbbreviation"`
	StatSortOrder    int    `json:"StatSortOrder"`
	StatPositionType string `json:"StatPositionType"`
	StatSortPosition int    `json:"StatSortPosition"`
}
