package logging

import (
	"log/slog"
	"os"
	"sync"
)

var (
	loggerInstance *slog.Logger
	once           sync.Once
)

// GetLogger initializes the logger only once and makes it reusable.
func GetLogger() *slog.Logger {
	once.Do(func() {
		loggerInstance = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	})
	return loggerInstance
}

// LogError logs error level messages.
func LogError(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}

// LogInfo logs informational level messages.
func LogInfo(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

// LogDebug logs debug level messages.
func LogDebug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}
