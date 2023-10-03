package logger

import (
	"log/slog"

	"github.com/jasonsites/gosk/internal/core/trace"
)

type AttrKeys struct {
	Body         string
	BodySize     string
	Headers      string
	IP           string
	Method       string
	PID          string
	Path         string
	ResponseTime string
	Status       string
	Tags         string
}

var AttrKey = AttrKeys{
	Body:         "body",
	BodySize:     "body_size",
	Headers:      "headers",
	IP:           "ip",
	Method:       "method",
	PID:          "pid",
	Path:         "path",
	ResponseTime: "response_time",
	Status:       "status",
	Tags:         "tags",
}

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
