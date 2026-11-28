package leagueclasses

// LeagueStatsValueInfo is one league scoring rule from StatValueTopic, as serialized by the
// Yahoo API's LeagueStatsValueModel: the fantasy points per unit of a Yahoo stat id, e.g.
// {"LeagueStatKey":4708388594,"GameId":470,"LeagueId":838859,"StatId":4,"StatValue":0.04}.
type LeagueStatsValueInfo struct {
	// The Yahoo API's concatenation of game id, league id and stat id. Not stored.
	LeagueStatKey int64    `json:"LeagueStatKey"`
	GameId        int      `json:"GameId"`
	LeagueId      int      `json:"LeagueId"`
	StatId        int      `json:"StatId"`
	StatValue     *float64 `json:"StatValue"`
}
