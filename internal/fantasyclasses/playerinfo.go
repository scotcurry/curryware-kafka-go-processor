package fantasyclasses

type PlayerInfo struct {
	PlayerSeasonId      string `json:"Key"`
	PlayerID            int    `json:"Id"`
	PlayerName          string `json:"FullName"`
	PlayerUrl           string `json:"URL"`
	PlayerStatus        string `json:"Status"`
	PlayerStatusFull    string `json:"StatusFull"`
	PlayerTeam          string `json:"Team"`
	PlayerByeWeek       int    `json:"ByeWeek"`
	PlayerUniformNumber int    `json:"UniformNumber"`
	PlayerPosition      string `json:"Position"`
	PlayerHeadshot      string `json:"Headshot"`
	PlayerInjuryNotes   string `json:"InjuryNotes"`
}
