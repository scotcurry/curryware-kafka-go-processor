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

func ExecuteGetLatestTransactionSelectStatement(ctx context.Context, sqlStatement string, leagueId string) (int, int) {
	row, err := QueryRowStatement(ctx, sqlStatement, leagueId)
	if err != nil {
		logger.LogError(ctx, "Error getting database connection for select statement", "error", err.Error())
		return -1, -1
	}
	var lastTransactionNumber int
	var lastTransactionDate int
	err = row.Scan(&lastTransactionNumber, &lastTransactionDate)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0
	}
	if err != nil {
		logger.LogError(ctx, "Error executing sql statement", "error", err.Error())
		return -1, -1
	}
	return lastTransactionNumber, lastTransactionDate
}
