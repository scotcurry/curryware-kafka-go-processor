package leagueclasses

type AllTeamInformation struct {
	LeagueKey              string `json:"LeagueKey"`
	TeamKey                string `json:"TeamKey"`
	TeamId                 int    `json:"TeamId"`
	TeamName               string `json:"TeamName"`
	TeamLogo               string `json:"TeamLogo"`
	PreviousSeasonTeamRank *int   `json:"PreviousSeasonTeamRank"`
	NumberOfMoves          int    `json:"NumberOfMoves"`
	NumberOfTrades         int    `json:"NumberOfTrades"`
	DraftPosition          int    `json:"DraftPosition"`
	DraftGrade             string `json:"DraftGrade"`
	ManagerNicknames       string `json:"ManagerNicknames"`
}
