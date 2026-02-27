package postgreshandlers

import (
	logger "curryware-kafka-go-processor/internal/logging"
	"database/sql"
	"errors"
)

func ExecuteSqlStatement(sqlStatement string, sqlParams []any) (int64, error) {
	return ExecStatement(sqlStatement, sqlParams...)
}

func ExecuteGetLatestTransactionSelectStatement(sqlStatement string, leagueId string) (int, int) {
	row := QueryRowStatement(sqlStatement, leagueId)
	var lastTransactionNumber int
	var lastTransactionDate int
	err := row.Scan(&lastTransactionNumber, &lastTransactionDate)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0
	}
	if err != nil {
		logger.LogError("Error executing sql statement", "error", err.Error())
		return -1, -1
	}
	return lastTransactionNumber, lastTransactionDate
}
