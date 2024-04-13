package logger

import (
	"log/slog"

	"github.com/jasonsites/gosk/internal/http/trace"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// CustomLogger encapsulates a logger with an associated log level and toggle
type CustomLogger struct {
	Enabled bool
	Level   string
	Log     *slog.Logger
}

// CreateContextLogger returns a new child logger with attached trace ID
func (l *CustomLogger) CreateContextLogger(traceID string) *slog.Logger {
	return l.Log.With(slog.String(string(trace.TraceIDContextKey), traceID))
}
