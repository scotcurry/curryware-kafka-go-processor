package fantasyclasses

type LeagueStatsValueInfo struct {
	LeagueStatId int     `json:"LeagueStatKey"`
	GameId       int     `json:"GameId"`
	LeagueId     int     `json:"LeagueId"`
	StatId       int     `json:"StatId"`
	StatValue    float64 `json:"StatValue"`
}
