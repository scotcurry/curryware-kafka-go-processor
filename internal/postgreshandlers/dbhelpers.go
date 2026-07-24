package postgreshandlers

import (
	"context"
	logger "curryware-kafka-go-processor/internal/logging"
	"database/sql"
)

// ExecStatement executes a SQL statement and returns the number of rows affected.
// It uses the singleton DB connection and logs errors, but does not panic.
func ExecStatement(ctx context.Context, sqlStatement string, params ...any) (int64, error) {
	db, err := GetDB(ctx)
	if err != nil {
		logger.LogError(ctx, "Error getting database connection", "error", err.Error())
		return 0, err
	}
	result, err := db.ExecContext(ctx, sqlStatement, params...)
	if err != nil {
		logger.LogError(ctx, "Error executing sql statement", "error", err.Error())
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.LogError(ctx, "Error getting rows affected", "error", err.Error())
		return 0, err
	}
	return rowsAffected, nil
}

// QueryRowStatement executes a query that returns a single row using the singleton DB connection.
func QueryRowStatement(ctx context.Context, sqlStatement string, params ...any) (*sql.Row, error) {
	db, err := GetDB(ctx)
	if err != nil {
		logger.LogError(ctx, "Error getting database connection", "error", err.Error())
		return nil, err
	}
	return db.QueryRowContext(ctx, sqlStatement, params...), nil
}
