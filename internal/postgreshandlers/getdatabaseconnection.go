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
// If the connection is not yet established, it retries with exponential backoff.
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

	// Retry with exponential backoff so the pod survives transient postgres
	// unavailability (e.g. pod starting before postgres is ready).
	maxRetries := 5
	backoff := 2 * time.Second

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			logger.LogInfo("Retrying postgres connection",
				"attempt", attempt+1,
				"backoff", backoff.String())
			time.Sleep(backoff)
			backoff *= 2
		}

		conn, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			logger.LogError("Error opening postgres connection", "error", err.Error())
			lastErr = err
			continue
		}

		// Configure connection pool settings
		conn.SetMaxOpenConns(25)
		conn.SetMaxIdleConns(25)
		conn.SetConnMaxLifetime(5 * time.Minute)

		// Test the connection
		if pingErr := conn.Ping(); pingErr != nil {
			logger.LogError("Error pinging postgres connection", "error", pingErr.Error())
			_ = conn.Close()
			lastErr = pingErr
			continue
		}

		db = conn
		return db, nil
	}

	logger.LogError("All postgres connection attempts failed", "attempts", maxRetries)
	return nil, lastErr
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
