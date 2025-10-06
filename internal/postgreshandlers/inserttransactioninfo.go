package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	"regexp"
	"time"
)

func InsertTransactionInfo(transactionJson fantasyclasses.TransactionInfoWithCount) int64 {

	// sqlStatement := "INSERT INTO latest_transaction_info (transaction_id, latest_transaction, last_transaction_date) VALUES ($1, $2, $3)"
	getLastTransactionStatement := "SELECT latest_transaction FROM latest_transaction_info WHERE transaction_id = $1"
	sqlParams := make([]interface{}, 0)
	leagueTransactionId := buildTransactionID(transactionJson)
	sqlParams = append(sqlParams, leagueTransactionId)

	rows := ExecuteSqlStatement(getLastTransactionStatement, sqlParams)

	// Since we got no rows back, we need to insert the latest transaction info.
	if rows == 0 {
		currentDateTime := time.Now()
		sqlStatement := "INSERT INTO latest_transaction_info (transaction_id, latest_transaction, last_transaction_date) VALUES ($1, $2, $3)"

		if transactionJson.TransactionCount > 0 {
			transactionId := buildTransactionID(transactionJson)
			insertParams := make([]interface{}, 0)
			insertParams = append(insertParams, transactionId)
			insertParams = append(insertParams, transactionJson.TransactionCount)
			insertParams = append(insertParams, currentDateTime)
			rows := ExecuteSqlStatement(sqlStatement, insertParams)

			return rows
		}
	}
	return rows
}

func buildTransactionID(transactionInfo fantasyclasses.TransactionInfoWithCount) string {

	transactionIdComplete := transactionInfo.Transactions[0].TransactionKey
	pattern := `^(\d+\.l\.\d+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(transactionIdComplete)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}
