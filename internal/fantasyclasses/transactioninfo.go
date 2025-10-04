package fantasyclasses

type TransactionInfoWithCount struct {
	TransactionCount int               `json:"TransactionCount"`
	Transactions     []TransactionInfo `json:"Transactions"`
}

type TransactionInfo struct {
	TransactionKey       string                  `json:"TransactionKey"`
	TransactionId        int                     `json:"TransactionId"`
	TransactionType      string                  `json:"TransactionType"`
	TransactionStatus    string                  `json:"TransactionStatus"`
	TransactionTimestamp int                     `json:"TransactionTimestamp"`
	PlayersInvolved      []TransactionPlayerInfo `json:"PlayersInvolved"`
}
