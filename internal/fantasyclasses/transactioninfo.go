package fantasyclasses

type TransactionInfoWithCount struct {
	LeagueKey        string            `json:"LeagueKey"`
	TransactionCount int               `json:"TransactionCount"`
	Transactions     []TransactionInfo `json:"Transactions"`
}

type TransactionInfo struct {
	GameID               int64                   `json:"GameId"`
	LeagueID             int64                   `json:"LeagueId"`
	TransactionKey       string                  `json:"TransactionKey"`
	TransactionId        int                     `json:"TransactionId"`
	TransactionType      string                  `json:"TransactionType"`
	TransactionStatus    string                  `json:"TransactionStatus"`
	TransactionTimestamp int64                   `json:"TransactionTimestamp"`
	PlayersInvolved      []TransactionPlayerInfo `json:"PlayersInvolved"`
}
