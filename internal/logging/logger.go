package logging

import (
	"context"
	"log/slog"
	"os"
	"sync"

	ddslog "github.com/DataDog/dd-trace-go/contrib/log/slog/v2"
)

var (
	loggerInstance *slog.Logger
	once           sync.Once
)

// GetLogger initializes the logger only once and makes it reusable.
func GetLogger() *slog.Logger {
	once.Do(func() {
		jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
		loggerInstance = slog.New(ddslog.WrapHandler(jsonHandler))
	})
	return loggerInstance
}

// LogError logs error level messages.
func LogError(ctx context.Context, msg string, args ...any) {
	GetLogger().ErrorContext(ctx, msg, args...)
}

// LogInfo logs informational level messages.
func LogInfo(ctx context.Context, msg string, args ...any) {
	GetLogger().InfoContext(ctx, msg, args...)
}

// LogDebug logs debug level messages.
func LogDebug(ctx context.Context, msg string, args ...any) {
	GetLogger().DebugContext(ctx, msg, args...)
}
