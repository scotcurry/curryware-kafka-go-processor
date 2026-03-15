package transactiontests

import (
	"curryware-kafka-go-processor/internal/fantasyclasses/transactionclasses"
	"curryware-kafka-go-processor/internal/jsonhandlers"
	"encoding/base64"
	"testing"
)

func TestParseTransactionInfo(t *testing.T) {

	transactionInfo := `{"TransactionCount":56,"Transactions":[{"TransactionKey":"461.l.460188.tr.61","TransactionId":61,"TransactionType":"add/drop","TransactionStatus":"successful","TransactionTimestamp":1758958222,"PlayersInvolved":[{"PlayerKey":"461.p.100002","PlayerId":100002,"TransactionKey":"461.l.460188.tr.61","TransactionType":"add","TransactionSource":"waivers","DestinationType":"team","DestinationTeamId":"461.l.460188.t.9"},{"PlayerKey":"461.p.100028","PlayerId":100028,"TransactionKey":"461.l.460188.tr.61","TransactionType":"drop","TransactionSource":"team","DestinationType":"waivers","DestinationTeamId":""}]}]}`
	encodedTransaction := base64.StdEncoding.EncodeToString([]byte(transactionInfo))

	result, err := jsonhandlers.ParseJSON[transactionclasses.TransactionInfoWithCount](encodedTransaction)
	if err != nil {
		t.Errorf("Error parsing JSON: %v", err)
	}

	if result.TransactionCount != 56 {
		t.Errorf("Expected TransactionCount 56, got %d", result.TransactionCount)
	}

	if len(result.Transactions) == 0 {
		t.Error("Expected at least one Transactions entry")
	}
}
