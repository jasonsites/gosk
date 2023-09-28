package logger

import (
	"log/slog"

	"github.com/jasonsites/gosk/internal/core/trace"
)

// Logger encapsulates a logger with an associated log level and toggle
type Logger struct {
	Enabled bool
	Level   string
	Log     *slog.Logger
}

// CreateContextLogger returns a new child logger with attached trace ID
func (l *Logger) CreateContextLogger(traceID string) *slog.Logger {
	return l.Log.With(slog.String(string(trace.TraceIDContextKey), traceID))
}
