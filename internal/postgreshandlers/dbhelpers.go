package postgreshandlers

import (
	logger "curryware-kafka-go-processor/internal/logging"
	"database/sql"
)

// ExecStatement executes a SQL statement and returns the number of rows affected.
// It uses the singleton DB connection and logs errors, but does not panic.
func ExecStatement(sqlStatement string, params ...any) (int64, error) {
	db := GetDB()
	result, err := db.Exec(sqlStatement, params...)
	if err != nil {
		logger.LogError("Error executing sql statement", "error", err.Error())
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.LogError("Error getting rows affected", "error", err.Error())
		return 0, err
	}
	return rowsAffected, nil
}

// QueryRowStatement executes a query that returns a single row using the singleton DB connection.
func QueryRowStatement(sqlStatement string, params ...any) *sql.Row {
	db := GetDB()
	return db.QueryRow(sqlStatement, params...)
}
