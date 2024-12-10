package fantasyclasses

type PlayerInfo struct {
	PlayerID            int    `json:"Id"`
	PlayerSeasonId      string `json:"Key"`
	PlayerName          string `json:"FullName"`
	PlayerUrl           string `json:"URL"`
	PlayerTeam          string `json:"Team"`
	PlayerByeWeek       int    `json:"ByeWeek"`
	PlayerUniformNumber int    `json:"UniformNumber"`
	PlayerPosition      string `json:"Position"`
	PlayerHeadshot      string `json:"Headshot"`
}
