package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"regexp"
	"time"
)

func InsertTransactionInfo(transactionJson fantasyclasses.TransactionInfoWithCount) int64 {

	// sqlStatement := "INSERT INTO latest_transaction_info (transaction_id, latest_transaction, last_transaction_date) VALUES ($1, $2, $3)"
	getLastTransactionStatement := "SELECT league_latest_transaction FROM latest_transaction_id WHERE league_transaction_id = $1"
	sqlParams := make([]interface{}, 0)
	leagueTransactionId := buildTransactionID(transactionJson)
	sqlParams = append(sqlParams, leagueTransactionId)

	rows := ExecuteSqlStatement(getLastTransactionStatement, sqlParams)

	// Since we got no rows back, we need to insert the latest transaction info.
	if rows == 0 {
		currentDateTime := time.Now()
		sqlStatement := "INSERT INTO latest_transaction_id (league_transaction_id, league_latest_transaction, last_transaction_date) VALUES ($1, $2, $3)"

		// This is to check if there are any transactions to insert.
		if transactionJson.TransactionCount > 0 {
			transactionId := buildTransactionID(transactionJson)
			insertParams := make([]interface{}, 0)
			insertParams = append(insertParams, transactionId)
			insertParams = append(insertParams, transactionJson.TransactionCount)
			insertParams = append(insertParams, currentDateTime)
			ExecuteSqlStatement(sqlStatement, insertParams)

			transactionsToInsert := transactionJson.Transactions
			addedTransactions := insertTransactionInfo(transactionsToInsert)

			return addedTransactions
		}
	}
	return rows
}

func insertTransactionInfo(transactionInfo []fantasyclasses.TransactionInfo) int64 {

	transactionInfoSqlStatement := "INSERT INTO transaction_info (transaction_key, transaction_id, transaction_type, transaction_status, transaction_time) VALUES ($1, $2, $3, $4, $5)"
	var totalRowsAdded int64 = 0
	for counter := 0; counter < len(transactionInfo); counter++ {

		transactionKey := transactionInfo[counter].TransactionKey
		transactionId := transactionInfo[counter].TransactionId
		transactionType := transactionInfo[counter].TransactionType
		transactionStatus := transactionInfo[counter].TransactionStatus
		transactionTime := transactionInfo[counter].TransactionTimestamp
		timestamp := time.Unix(transactionTime, 0)

		sqlParams := make([]interface{}, 0)
		sqlParams = append(sqlParams, transactionKey)
		sqlParams = append(sqlParams, transactionId)
		sqlParams = append(sqlParams, transactionType)
		sqlParams = append(sqlParams, transactionStatus)
		sqlParams = append(sqlParams, timestamp)

		rowCount := ExecuteSqlStatement(transactionInfoSqlStatement, sqlParams)
		logger.LogInfo("Rows inserted: {1}", rowCount)

		players := transactionInfo[counter].PlayersInvolved

		for playerCounter := 0; playerCounter < len(players); playerCounter++ {
			playerKey := players[playerCounter].PlayerKey
			playerId := players[playerCounter].PlayerId
			playerTransactionType := players[playerCounter].DestinationType
			playerTransactionSource := players[playerCounter].TransactionSource
			playerTransactionDestination := players[playerCounter].DestinationType
			playerTransactionDestinationTeamId := players[playerCounter].DestinationTeamId

			sqlParams := make([]interface{}, 0)
			sqlParams = append(sqlParams, transactionKey)
			sqlParams = append(sqlParams, playerKey)
			sqlParams = append(sqlParams, playerId)
			sqlParams = append(sqlParams, playerTransactionType)
			sqlParams = append(sqlParams, playerTransactionSource)
			sqlParams = append(sqlParams, playerTransactionDestination)
			sqlParams = append(sqlParams, playerTransactionDestinationTeamId)

			playerInsertSqlStatement := "INSERT INTO transaction_player (transaction_key, player_key, player_id, transaction_type, transaction_source, destination_team, destination_team_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"

			rows := ExecuteSqlStatement(playerInsertSqlStatement, sqlParams)
			totalRowsAdded += rows
		}
	}
	return totalRowsAdded
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
