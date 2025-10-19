package postgreshandlers

import (
	"curryware-kafka-go-processor/internal/fantasyclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"time"
)

func ProcessTransactionInfo(transactionJson fantasyclasses.TransactionInfoWithCount) int64 {

	leagueKey := transactionJson.LeagueKey
	databaseLastTransaction, lastTransactionDate := getLastTransactionFromDatabase(leagueKey)
	println(databaseLastTransaction)
	// No transactions in the database, so we need to insert everything.
	rowCount := insertTransactionInfo(transactionJson, lastTransactionDate)
	logger.LogInfo("Database Last Transaction: {1}", databaseLastTransaction)

	return rowCount
}

// Call to the database to see if any action is needed.
func getLastTransactionFromDatabase(leagueKey string) (int64, int64) {

	getLastTransactionStatement := "SELECT league_latest_transaction, last_transaction_date FROM latest_transaction_id WHERE league_transaction_id = $1"
	latestTransActionId, latestTransactionDate := ExecuteGetLatestTransactionSelectStatement(getLastTransactionStatement, leagueKey)

	return int64(latestTransActionId), int64(latestTransactionDate)
}

func updateLatestTransactions(transactionJson fantasyclasses.TransactionInfoWithCount, latestTransaction int, lastTransactionDate int) int64 {

	leagueKey := transactionJson.LeagueKey

	var rows int64 = 0
	for counter := 0; counter < len(transactionJson.Transactions); counter++ {
		transactionDate := int(transactionJson.Transactions[counter].TransactionTimestamp)
		if transactionDate > lastTransactionDate {
			transactionToInsert := transactionJson.Transactions[counter]
			rows := insertTransactionDetail(transactionToInsert)
			logger.LogInfo("Rows inserted: {1}", rows)
		}
	}

	updateLatestTransactionStatement := "UPDATE latest_transaction_id SET league_latest_transaction = $1 WHERE league_transaction_id = $2"
	sqlParams := make([]interface{}, 0)
	sqlParams = append(sqlParams, latestTransaction)
	sqlParams = append(sqlParams, leagueKey)
	rows = ExecuteSqlStatement(updateLatestTransactionStatement, sqlParams)
	return rows
}

func insertTransactionInfo(transactionJson fantasyclasses.TransactionInfoWithCount, lastTransactionDate int64) int64 {

	var totalRows int64 = 0
	allTransactions := transactionJson.Transactions
	for counter := 0; counter < len(allTransactions); counter++ {
		if allTransactions[counter].TransactionTimestamp > lastTransactionDate {
			rows := insertTransactionDetail(allTransactions[counter])
			totalRows += rows
		}
	}

	return totalRows
}

func insertTransactionDetail(transactionToInsert fantasyclasses.TransactionInfo) int64 {

	transactionInfoSqlStatement := "INSERT INTO transaction_info (transaction_key, transaction_id, transaction_type, transaction_status, transaction_time) VALUES ($1, $2, $3, $4, $5)"
	var totalRowsAdded int64 = 0

	transactionKey := transactionToInsert.TransactionKey
	transactionId := transactionToInsert.TransactionId
	transactionType := transactionToInsert.TransactionType
	transactionStatus := transactionToInsert.TransactionStatus
	transactionTime := transactionToInsert.TransactionTimestamp
	timestamp := time.Unix(transactionTime, 0)

	sqlParams := make([]interface{}, 0)
	sqlParams = append(sqlParams, transactionKey)
	sqlParams = append(sqlParams, transactionId)
	sqlParams = append(sqlParams, transactionType)
	sqlParams = append(sqlParams, transactionStatus)
	sqlParams = append(sqlParams, timestamp)

	rowCount := ExecuteSqlStatement(transactionInfoSqlStatement, sqlParams)
	logger.LogInfo("Rows inserted: {1}", rowCount)

	players := transactionToInsert.PlayersInvolved

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
	return totalRowsAdded
}
