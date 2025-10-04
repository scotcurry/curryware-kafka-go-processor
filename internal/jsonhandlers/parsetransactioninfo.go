package jsonhandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"curryware-kafka-go-processor/internal/logging"
	"encoding/base64"
	"encoding/json"
)

func ParseTransactionInfo(transactionInfo string) fantasyclasses.TransactionInfoWithCount {

	decodedBytes, err := base64.StdEncoding.DecodeString(transactionInfo)
	transactionInfo = string(decodedBytes)

	var transactionValueStruct fantasyclasses.TransactionInfoWithCount
	err = json.Unmarshal([]byte(transactionInfo), &transactionValueStruct)
	if err != nil {
		logging.LogError("Error Parsing Stat Values JSON :", err)
	}

	return transactionValueStruct
}
