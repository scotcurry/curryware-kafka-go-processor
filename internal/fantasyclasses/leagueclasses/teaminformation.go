package leagueclasses

type TeamInformation struct {
	TeamKey          string `json:"TeamKey"`
	TeamName         string `json:"TeamName"`
	TeamUrl          string `json:"TeamUrl"`
	TeamLogoUrl      string `json:"TeamLogoUrl"`
	DraftPosition    int    `json:"DraftPosition"`
	DraftGrade       string `json:"DraftGrade"`
	ManagerNickname  string `json:"ManagerNickname"`
	ManagerImageUrl  string `json:"ManagerImageUrl"`
	ManagerFeloScore int    `json:"ManagerFeloScore"`
}
