package postgreshandlers

import (
	logger "curryware-kafka-go-processor/internal/logging"
	"fmt"
)

func ExecuteSqlStatement(sqlStatement string, sqlParams []interface{}) int64 {

	db := GetDB()
	result, err := db.Exec(sqlStatement, sqlParams...)
	if err != nil {
		logger.LogError("Error executing sql statement: ", err)
		fmt.Println("Error executing sql statement: ", err)
	}

	rowsAffected, err := result.RowsAffected()
	return rowsAffected
}
