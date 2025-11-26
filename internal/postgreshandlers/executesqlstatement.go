package postgreshandlers

import (
	logger "curryware-kafka-go-processor/internal/logging"
	"database/sql"
	"errors"
	"fmt"
)

func ExecuteSqlStatement(sqlStatement string, sqlParams []any) (int64, error) {

	db := GetDB()
	result, err := db.Exec(sqlStatement, sqlParams...)
	if err != nil {
		logger.LogError("Error executing sql statement: ", err)
		fmt.Println("Error executing sql statement: ", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.LogError("Error getting rows affected", "error", err.Error())
		return 0, err
	}
	return rowsAffected, nil
}

func ExecuteGetLatestTransactionSelectStatement(sqlStatement string, leagueId string) (int, int) {

	row := GetDB().QueryRow(sqlStatement, leagueId)
	var lastTransactionNumber int
	var lastTransactionDate int
	err := row.Scan(&lastTransactionNumber, &lastTransactionDate)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0
	}
	if err != nil {
		logger.LogError("Error executing sql statement: ", err)
		fmt.Println("Error executing sql statement: ", err)
		return -1, -1
	} else {
		return lastTransactionNumber, lastTransactionDate
	}
}
