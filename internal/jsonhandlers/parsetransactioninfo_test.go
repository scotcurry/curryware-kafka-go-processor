package jsonhandlers

import (
	"encoding/base64"
	"testing"
)

func TestParseTransactionInfo(t *testing.T) {

	transactionInfo := `{"TransactionCount":56,"Transactions":[{"TransactionKey":"461.l.460188.tr.61","TransactionId":61,"TransactionType":"add/drop","TransactionStatus":"successful","TransactionTimestamp":1758958222,"PlayersInvolved":[{"PlayerKey":"461.p.100002","PlayerId":100002,"TransactionKey":"461.l.460188.tr.61","TransactionType":"add","TransactionSource":"waivers","DestinationType":"team","DestinationTeamId":"461.l.460188.t.9"},{"PlayerKey":"461.p.100028","PlayerId":100028,"TransactionKey":"461.l.460188.tr.61","TransactionType":"drop","TransactionSource":"team","DestinationType":"waivers","DestinationTeamId":""}]},{"TransactionKey":"461.l.460188.tr.60","TransactionId":60,"TransactionType":"add/drop","TransactionStatus":"successful","TransactionTimestamp":1758813598,"PlayersInvolved":[{"PlayerKey":"461.p.100017","PlayerId":100017,"TransactionKey":"461.l.460188.tr.60","TransactionType":"add","TransactionSource":"freeagents","DestinationType":"team","DestinationTeamId":"461.l.460188.t.3"},{"PlayerKey":"461.p.100026","PlayerId":100026,"TransactionKey":"461.l.460188.tr.60","TransactionType":"drop","TransactionSource":"team","DestinationType":"waivers","DestinationTeamId":""}]},{"TransactionKey":"461.l.460188.tr.59","TransactionId":59,"TransactionType":"add/drop","TransactionStatus":"successful","TransactionTimestamp":1758813564,"PlayersInvolved":[{"PlayerKey":"461.p.29269","PlayerId":29269,"TransactionKey":"461.l.460188.tr.59","TransactionType":"add","TransactionSource":"freeagents","DestinationType":"team","DestinationTeamId":"461.l.460188.t.3"},{"PlayerKey":"461.p.33443","PlayerId":33443,"TransactionKey":"461.l.460188.tr.59","TransactionType":"drop","TransactionSource":"team","DestinationType":"waivers","DestinationTeamId":""}]},{"TransactionKey":"461.l.460188.tr.58","TransactionId":58,"TransactionType":"add","TransactionStatus":"successful","TransactionTimestamp":1758811390,"PlayersInvolved":[{"PlayerKey":"461.p.29236","PlayerId":29236,"TransactionKey":"461.l.460188.tr.58","TransactionType":"add","TransactionSource":"freeagents","DestinationType":"team","DestinationTeamId":"461.l.460188.t.4"}]},{"TransactionKey":"461.l.460188.tr.57","TransactionId":57,"TransactionType":"drop","TransactionStatus":"successful","TransactionTimestamp":1758811280,"PlayersInvolved":[{"PlayerKey":"461.p.25785","PlayerId":25785,"TransactionKey":"461.l.460188.tr.57","TransactionType":"drop","TransactionSource":"team","DestinationType":"waivers","DestinationTeamId":""}]},{"TransactionKey":"461.l.460188.tr.56","TransactionId":56,"TransactionType":"add/drop","TransactionStatus":"successful","TransactionTimestamp":1758798125,"PlayersInvolved":[{"PlayerKey":"461.p.100023","PlayerId":100023,"TransactionKey":"461.l.460188.tr.56","TransactionType":"add","TransactionSource":"freeagents","DestinationType":"team","DestinationTeamId":"461.l.460188.t.2"}]}]}`
	transactionInfo64 := base64.StdEncoding.EncodeToString([]byte(transactionInfo))

	transactionValueClass := ParseTransactionInfo(transactionInfo64)
	if transactionValueClass.TransactionCount != 56 {
		t.Error("TransactionCount not parsed correctly")
	}
}
