package statsclasses

type PlayerWeeklyStatsInfo struct {
	PlayerId       int     `json:"playerId"`
	PlayerGameKey  string  `json:"playerGameKey"`
	PlayerStatWeek int     `json:"playerStatWeek"`
	StatId         int     `json:"statId"`
	StatValue      float64 `json:"statValue"`
}
