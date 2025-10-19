package postgreshandlers

import (
	logger "curryware-kafka-go-processor/internal/logging"
	"database/sql"
	"sync"
	"time"
)

var (
	db   *sql.DB
	once sync.Once
)

// GetDB returns a singleton database connection pool
func GetDB() *sql.DB {
	once.Do(func() {
		psqlInfo, variableError := GetDatabaseInformation()
		if variableError != nil {
			logger.LogError("Error getting database information")
			panic(variableError)
		}

		var err error
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			logger.LogError("Error opening postgres connection")
			panic(err)
		}

		// Configure connection pool settings
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		// Test the connection
		if err = db.Ping(); err != nil {
			logger.LogError("Error pinging postgres connection")
			println(err)
			panic(err)
		}
	})
	return db
}
