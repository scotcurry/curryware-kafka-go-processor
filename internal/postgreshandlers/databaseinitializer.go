package postgreshandlers

import (
	logger "curryware-kafka-go-processor/internal/logging"
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
)

// Database This is the shared *sql.DB connection.
type Database struct {
	db   *sql.DB
	once sync.Once
}

func GetDb() *sql.DB {
	psqlInfo, err := GetDatabaseInformation()
	if err != nil {
		logger.LogError("Error getting database information")
		return nil
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil
	}

	return db
}

func CloseDB() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
