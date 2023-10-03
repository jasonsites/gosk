package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/jasonsites/gosk/internal/core/logger"
	"github.com/jasonsites/gosk/internal/core/trace"
)

// RequestLogContextKey
var RequestLogContextKey trace.ContextKey

// RequestLoggerConfig defines necessary components for the request logger middleware
type RequestLoggerConfig struct {
	ContextKey trace.ContextKey
	Logger     *logger.Logger
	Next       func(r *http.Request) bool
}

// RequestLogger returns the request logger middleware
func RequestLogger(config *RequestLoggerConfig) func(http.Handler) http.Handler {
	conf := setRequestLoggerConfig(config)
	logger := conf.Logger

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if conf.Next != nil && conf.Next(r) {
				next.ServeHTTP(w, r)
				return
			}

			data, err := logRequest(w, r, logger)
			if err != nil {
				logger.Log.Error(err.Error())
			}

			ctx := context.WithValue(r.Context(), RequestLogContextKey, data)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// setRequestLoggerConfig returns a valid request logger configuration and sets the context key for external use
func setRequestLoggerConfig(c *RequestLoggerConfig) *RequestLoggerConfig {
	if c.Logger == nil {
		log.Panicf("request logger middleware missing logger configuration")
	}
	conf := c

	// override defaults
	if c.ContextKey == "" {
		conf.ContextKey = "request_data"
	}
	// set middleware-scoped context key for use in other handlers
	RequestLogContextKey = conf.ContextKey

	return conf
}

// logRequest logs the captured request data
func logRequest(w http.ResponseWriter, r *http.Request, logger *logger.Logger) (*RequestLogData, error) {
	if logger.Enabled {
		traceID := trace.GetTraceIDFromContext(r.Context())
		log := logger.CreateContextLogger(traceID)

		var body map[string]any
		if logger.Level == "debug" {
			r.Body = http.MaxBytesReader(w, r.Body, int64(1048576))

			copy, err := io.ReadAll(r.Body)
			if err != nil {
				log.Error(err.Error())
				return nil, err
			}

			if len(copy) > 0 {
				if err := json.Unmarshal(copy, &body); err != nil {
					log.Error(err.Error())
					return nil, err
				}
			}

			r.Body = io.NopCloser(bytes.NewBuffer(copy))
		}

		data := newRequestLogData(r, body)
		attrs := newRequestLogEvent(data, logger.Level, log)
		log.With(attrs...).Info("request")

		return data, nil
	}

	return nil, nil
}

// RequestLogData defines the data captured for request logging
type RequestLogData struct {
	Body     map[string]any
	ClientIP string
	Headers  http.Header
	Method   string
	Path     string
}

// newRequestLogData captures relevant data from the request
func newRequestLogData(r *http.Request, body map[string]any) *RequestLogData {
	return &RequestLogData{
		Body:     body,
		ClientIP: r.RemoteAddr,
		Headers:  r.Header,
		Method:   r.Method,
		Path:     r.URL.Path,
	}
}

// newRequestLogEvent returns additional attributes for logging
func newRequestLogEvent(data *RequestLogData, level string, log *slog.Logger) []any {
	attrs := []any{
		slog.String("ip", data.ClientIP),
		slog.String("method", data.Method),
		slog.String("path", data.Path),
	}

	if level == "debug" {
		if data.Headers != nil {
			attrs = append(attrs, "headers", data.Headers)
		}
		if data.Body != nil {
			attrs = append(attrs, "body", data.Body)
		}
	}

	return attrs
}
