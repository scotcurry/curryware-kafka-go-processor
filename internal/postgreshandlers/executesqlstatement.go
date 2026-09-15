package postgreshandlers

import (
	"context"
	logger "curryware-kafka-go-processor/internal/logging"
	"database/sql"
	"errors"
)

func ExecuteSqlStatement(ctx context.Context, sqlStatement string, sqlParams []any) (int64, error) {
	return ExecStatement(ctx, sqlStatement, sqlParams...)
}

// ExecuteGetLatestTransactionSelectStatement returns the last known transaction number and date for a
// game_id/league_id combination, along with whether a row was found. found is false when there is no
// row yet for that combination in latest_transaction_id (or the query failed), meaning every
// transaction in the incoming payload is new.
func ExecuteGetLatestTransactionSelectStatement(ctx context.Context, sqlStatement string, gameId int64, leagueId int64) (int, int, bool) {
	row, err := QueryRowStatement(ctx, sqlStatement, gameId, leagueId)
	if err != nil {
		logger.LogError(ctx, "Error getting database connection for select statement", "error", err.Error())
		return 0, 0, false
	}
	var lastTransactionNumber int
	var lastTransactionDate int
	err = row.Scan(&lastTransactionNumber, &lastTransactionDate)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0, false
	}
	if err != nil {
		logger.LogError(ctx, "Error executing sql statement", "error", err.Error())
		return 0, 0, false
	}
	return lastTransactionNumber, lastTransactionDate, true
}
