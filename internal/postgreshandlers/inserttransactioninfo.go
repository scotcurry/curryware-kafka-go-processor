package postgreshandlers

import (
	"context"
	"curryware-kafka-go-processor/internal/fantasyclasses/transactionclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"time"
)

func ProcessTransactionInfo(ctx context.Context, transactionJson transactionclasses.TransactionInfoWithCount) int64 {

	leagueKey := transactionJson.LeagueKey
	databaseLastTransaction, lastTransactionDate := getLastTransactionFromDatabase(ctx, leagueKey)
	logger.LogDebug(ctx, "Database last transaction", "transaction", databaseLastTransaction, "date", lastTransactionDate)
	rowCount := updateLatestTransactions(ctx, transactionJson, lastTransactionDate)
	logger.LogInfo(ctx, "Database Last Transaction", "transaction", databaseLastTransaction)

	return rowCount
}

// Call the database to see if any action is needed.
func getLastTransactionFromDatabase(ctx context.Context, leagueKey string) (int64, int64) {

	getLastTransactionStatement := "SELECT league_latest_transaction, last_transaction_date FROM latest_transaction_id WHERE league_transaction_id = $1"
	latestTransActionId, latestTransactionDate := ExecuteGetLatestTransactionSelectStatement(ctx, getLastTransactionStatement, leagueKey)

	return int64(latestTransActionId), int64(latestTransactionDate)
}

// This is to set the pointer so the next time only new transactions are inserted.
func updateLatestTransactions(ctx context.Context, transactionJson transactionclasses.TransactionInfoWithCount, lastTransactionDate int64) int64 {

	leagueKey := transactionJson.LeagueKey

	var totalRows int64 = 0
	latestTransaction := 0
	latestTransactionDate := lastTransactionDate
	for counter := 0; counter < len(transactionJson.Transactions); counter++ {
		transactionToInsert := transactionJson.Transactions[counter]
		transactionDate := transactionToInsert.TransactionTimestamp

		if transactionToInsert.TransactionId > latestTransaction {
			latestTransaction = transactionToInsert.TransactionId
		}
		if transactionDate > latestTransactionDate {
			latestTransactionDate = transactionDate
		}

		if transactionDate > lastTransactionDate {
			rows, err := insertTransactionDetail(ctx, transactionToInsert)
			if err != nil {
				logger.LogError(ctx, "Error inserting transaction info", "error", err)
				continue
			}
			logger.LogInfo(ctx, "Rows inserted", "rowCount", rows)
			totalRows += rows
		}
	}

	if latestTransaction == 0 {
		return totalRows
	}

	updateLatestTransactionStatement := "UPDATE latest_transaction_id SET league_latest_transaction = $1, last_transaction_date = $2 WHERE league_transaction_id = $3"
	sqlParams := make([]interface{}, 0)
	sqlParams = append(sqlParams, latestTransaction)
	sqlParams = append(sqlParams, latestTransactionDate)
	sqlParams = append(sqlParams, leagueKey)
	_, err := ExecuteSqlStatement(ctx, updateLatestTransactionStatement, sqlParams)
	if err != nil {
		logger.LogError(ctx, "Error updating latest transaction id", "error", err)
	}
	return totalRows
}

func insertTransactionDetail(ctx context.Context, transactionToInsert transactionclasses.TransactionInfo) (int64, error) {

	transactionInfoSqlStatement := "INSERT INTO transaction_info (game_id, league_id, transaction_key, transaction_id, transaction_type, transaction_status, transaction_time) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (transaction_key) DO NOTHING;"
	var totalRowsAdded int64 = 0

	gameId := transactionToInsert.GameID
	leagueId := transactionToInsert.LeagueID
	transactionKey := transactionToInsert.TransactionKey
	transactionId := transactionToInsert.TransactionId
	transactionType := transactionToInsert.TransactionType
	transactionStatus := transactionToInsert.TransactionStatus
	transactionTime := transactionToInsert.TransactionTimestamp
	timestamp := time.Unix(transactionTime, 0)

	sqlParams := make([]interface{}, 0)
	sqlParams = append(sqlParams, gameId)
	sqlParams = append(sqlParams, leagueId)
	sqlParams = append(sqlParams, transactionKey)
	sqlParams = append(sqlParams, transactionId)
	sqlParams = append(sqlParams, transactionType)
	sqlParams = append(sqlParams, transactionStatus)
	sqlParams = append(sqlParams, timestamp)

	if gameId == 0 || leagueId == 0 {
		logger.LogError(ctx, "GameId or LeagueId is 0")
	}

	rowCount, err := ExecuteSqlStatement(ctx, transactionInfoSqlStatement, sqlParams)
	if err != nil {
		logger.LogError(ctx, "Error inserting transaction info", "error", err)
		return 0, err
	}

	if rowCount == 0 {
		logger.LogInfo(ctx, "No rows inserted, record exists", "transactionKey", transactionKey)
		return rowCount, nil
	}
	logger.LogInfo(ctx, "Rows inserted", "rowCount", rowCount)

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

		rows, err := ExecuteSqlStatement(ctx, playerInsertSqlStatement, sqlParams)
		if err != nil {
			logger.LogError(ctx, "Error inserting transaction player info", "error", err)
			return 0, err
		}
		totalRowsAdded += rows
	}
	return totalRowsAdded, nil
}
