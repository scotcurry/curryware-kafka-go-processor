package fantasyclasses

type StatsInfo struct {
	PlayerID  int     `json:"PlayerId"`
	GameKey   int     `json:"GameKey"`
	WeekKey   int     `json:"WeekKey"`
	StatId    int     `json:"StatId"`
	StatValue float64 `json:"StatValue"`
}

type StatsJson struct {
	PlayerStats []StatsInfo `json:"PlayerStats"`
}
