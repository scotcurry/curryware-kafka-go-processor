package fantasyclasses

type PlayerStatValueInfo struct {
	PlayerId       int     `json:"playerId"`
	PlayerGameKey  int     `json:"playerGameKey"`
	PlayerStatWeek int     `json:"playerStatWeek"`
	StatId         int     `json:"statId"`
	StatValue      float64 `json:"statValue"`
}
