package teamclasses

type TeamWeeklyStandingsInfo struct {
	TeamKey       string  `json:"TeamKey"`
	TeamId        int     `json:"TeamId"`
	WeekKey       int     `json:"WeekKey"`
	Rank          int     `json:"Rank"`
	PointsFor     float64 `json:"PointsFor"`
	PointsAgainst float64 `json:"PointsAgainst"`
	Wins          int     `json:"Wins"`
	Losses        int     `json:"Losses"`
	Ties          int     `json:"Ties"`
	NumberOfMoves int     `json:"NumberOfMoves"`
}
