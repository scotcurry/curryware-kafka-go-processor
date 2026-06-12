package leagueclasses

type LeagueInformation struct {
	GameKey               string `json:"GameKey"`
	LeagueKey             string `json:"LeagueKey"`
	LeagueId              int    `json:"LeagueId"`
	LeagueName            string `json:"LeagueName"`
	LeagueLogoUrl         string `json:"LeagueLogoUrl"`
	NumberOfTeams         int    `json:"NumberOfTeams"`
	LeagueUpdateTimestamp string `json:"LeagueUpdateTimestamp"`
	StartDate             string `json:"StartDate"`
	EndDate               string `json:"EndDate"`
	Season                int    `json:"Season"`
}
