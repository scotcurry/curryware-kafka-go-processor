package leagueclasses

type TeamInformation struct {
	LeagueKey              string `json:"LeagueKey"`
	TeamKey                string `json:"TeamKey"`
	TeamID                 int    `json:"TeamId"`
	TeamName               string `json:"TeamName"`
	TeamLogo               string `json:"TeamLogo"`
	PreviousSeasonTeamRank int    `json:"PreviousSeasonTeamRank"`
	NumberOfMoves          int    `json:"NumberOfMoves"`
	NumberOfTrades         int    `json:"NumberOfTrades"`
	DraftPosition          int    `json:"DraftPosition"`
	DraftGrade             string `json:"DraftGrade"`
	ManagerNickname        string `json:"ManagerNickname"`
	ManagerFeloScore       int    `json:"ManagerFeloScore"`
}
