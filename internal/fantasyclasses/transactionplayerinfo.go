package fantasyclasses

type TransactionPlayerInfo struct {
	PlayerKey         string `json:"PlayerKey"`
	PlayerId          int    `json:"PlayerId"`
	TransactionKey    string `json:"TransactionKey"`
	TransactionSource string `json:"TransactionSource"`
	DestinationType   string `json:"DestinationType"`
	DestinationTeamId string `json:"DestinationTeamId"`
}
