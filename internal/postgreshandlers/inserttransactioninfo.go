package postgreshandlers

import (
	"context"
	"curryware-kafka-go-processor/internal/fantasyclasses/transactionclasses"
	logger "curryware-kafka-go-processor/internal/logging"
	"time"
)

func ProcessTransactionInfo(ctx context.Context, transactionJson transactionclasses.TransactionInfoWithCount) int64 {

	if len(transactionJson.Transactions) == 0 {
		logger.LogInfo(ctx, "No transactions in payload, nothing to process")
		return 0
	}

	gameId := transactionJson.Transactions[0].GameID
	leagueId := transactionJson.Transactions[0].LeagueID

	databaseLastTransaction, lastTransactionDate, found := getLastTransactionFromDatabase(ctx, gameId, leagueId)
	logger.LogDebug(ctx, "Database last transaction", "transaction", databaseLastTransaction, "date", lastTransactionDate, "found", found)
	rowCount := updateLatestTransactions(ctx, transactionJson, gameId, leagueId, databaseLastTransaction, lastTransactionDate, found)

	return rowCount
}

// Call the database to see if any action is needed. found is false when the game_id/league_id
// combination has no row yet in latest_transaction_id, meaning every transaction in the incoming
// payload should be treated as new.
func getLastTransactionFromDatabase(ctx context.Context, gameId int64, leagueId int64) (int, int64, bool) {

	getLastTransactionStatement := "SELECT league_latest_transaction, last_transaction_date FROM latest_transaction_id WHERE game_id = $1 AND league_id = $2"
	latestTransActionId, latestTransactionDate, found := ExecuteGetLatestTransactionSelectStatement(ctx, getLastTransactionStatement, gameId, leagueId)

	return latestTransActionId, int64(latestTransactionDate), found
}

// updateLatestTransactions inserts only the transactions newer than the database's last known transaction
// (or all of them if the game_id/league_id combination has no row yet), then upserts the pointer row so
// the next run only picks up transactions past this point.
func updateLatestTransactions(ctx context.Context, transactionJson transactionclasses.TransactionInfoWithCount, gameId int64, leagueId int64, databaseLastTransaction int, lastTransactionDate int64, found bool) int64 {

	leagueKey := transactionJson.LeagueKey

	var totalRows int64 = 0
	latestTransaction := databaseLastTransaction
	latestTransactionDate := lastTransactionDate

	for counter := 0; counter < len(transactionJson.Transactions); counter++ {
		transactionToInsert := transactionJson.Transactions[counter]

		if found && transactionToInsert.TransactionId <= databaseLastTransaction {
			continue
		}

		rows, err := insertTransactionDetail(ctx, transactionToInsert)
		if err != nil {
			logger.LogError(ctx, "Error inserting transaction info", "error", err)
			continue
		}
		logger.LogInfo(ctx, "Rows inserted", "rowCount", rows)
		totalRows += rows

		if transactionToInsert.TransactionId > latestTransaction {
			latestTransaction = transactionToInsert.TransactionId
		}
		if transactionToInsert.TransactionTimestamp > latestTransactionDate {
			latestTransactionDate = transactionToInsert.TransactionTimestamp
		}
	}

	if totalRows == 0 {
		return totalRows
	}

	upsertLatestTransactionStatement := `INSERT INTO latest_transaction_id (game_id, league_id, league_transaction_id, league_latest_transaction, last_transaction_date)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (game_id, league_id) DO UPDATE
		SET league_transaction_id = EXCLUDED.league_transaction_id,
		    league_latest_transaction = EXCLUDED.league_latest_transaction,
		    last_transaction_date = EXCLUDED.last_transaction_date`
	sqlParams := make([]any, 0)
	sqlParams = append(sqlParams, gameId)
	sqlParams = append(sqlParams, leagueId)
	sqlParams = append(sqlParams, leagueKey)
	sqlParams = append(sqlParams, latestTransaction)
	sqlParams = append(sqlParams, latestTransactionDate)
	_, err := ExecuteSqlStatement(ctx, upsertLatestTransactionStatement, sqlParams)
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

	sqlParams := make([]any, 0)
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

	for playerCounter := range players {
		playerKey := players[playerCounter].PlayerKey
		playerId := players[playerCounter].PlayerId
		playerTransactionType := players[playerCounter].DestinationType
		playerTransactionSource := players[playerCounter].TransactionSource
		playerTransactionDestination := players[playerCounter].DestinationType
		playerTransactionDestinationTeamId := players[playerCounter].DestinationTeamId

		sqlParams := make([]any, 0)
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
