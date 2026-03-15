package postgreshandlers

import (
	logger "curryware-kafka-go-processor/internal/logging"
	"database/sql"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// db is a global variable that can be used by all the database function calls.
var (
	db   *sql.DB
	dbMu sync.Mutex
)

// GetDB returns a singleton database connection pool.
// Returns (*sql.DB, error) so callers can handle connection failures gracefully.
func GetDB() (*sql.DB, error) {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db != nil {
		return db, nil
	}

	psqlInfo, variableError := GetDatabaseInformation()
	if variableError != nil {
		logger.LogError("Error getting database information")
		return nil, variableError
	}

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.LogError("Error opening postgres connection")
		return nil, err
	}

	// Configure connection pool settings
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err = conn.Ping(); err != nil {
		logger.LogError("Error pinging postgres connection", "error", err.Error())
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	db = conn
	return db, nil
}

// CloseDB closes the singleton database connection pool.
func CloseDB() error {
	dbMu.Lock()
	defer dbMu.Unlock()
	if db == nil {
		return nil
	}
	err := db.Close()
	db = nil
	return err
}
