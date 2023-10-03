package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	cl "github.com/jasonsites/gosk/internal/core/logger"
	"github.com/jasonsites/gosk/internal/core/trace"
	"github.com/jasonsites/gosk/internal/core/validation"
	ctrl "github.com/jasonsites/gosk/internal/http/controllers"
)

// RequestLogContextKey
const RequestLogContextKey trace.ContextKey = "request_data"

// RequestLogData defines the data captured for request logging
type RequestLogData struct {
	Body     map[string]any
	ClientIP string
	Headers  http.Header
	Method   string
	Path     string
}

// RequestLoggerConfig defines necessary components for the request logger middleware
type RequestLoggerConfig struct {
	ContextKey trace.ContextKey
	Logger     *cl.CustomLogger `validate:"required"`
	Next       func(r *http.Request) bool
}

// RequestLogger returns the request logger middleware
func RequestLogger(c *RequestLoggerConfig) func(http.Handler) http.Handler {
	if err := validation.Validate.Struct(c); err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c.Next != nil && c.Next(r) {
				next.ServeHTTP(w, r)
				return
			}

			data, err := logRequest(w, r, c.Logger)
			if err != nil {
				c.Logger.Log.Error(err.Error())
				ctrl.EncodeError(w, r, err) // TODO: likely move http encode/decode from ctrl
			}

			ctx := context.WithValue(r.Context(), RequestLogContextKey, data)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// logRequest logs the captured request data
func logRequest(w http.ResponseWriter, r *http.Request, logger *cl.CustomLogger) (*RequestLogData, error) {
	if logger.Enabled {
		traceID := trace.GetTraceIDFromContext(r.Context())
		logger.Log = logger.CreateContextLogger(traceID)

		var body map[string]any
		if logger.Level == cl.LevelDebug {
			r.Body = http.MaxBytesReader(w, r.Body, int64(1048576))

			copy, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Log.Error(err.Error())
				return nil, err
			}

			if len(copy) > 0 {
				if err := json.Unmarshal(copy, &body); err != nil {
					logger.Log.Error(err.Error())
					return nil, err
				}
			}

			r.Body = io.NopCloser(bytes.NewBuffer(copy))
		}

		data := &RequestLogData{
			Body:     body,
			ClientIP: r.RemoteAddr,
			Headers:  r.Header,
			Method:   r.Method,
			Path:     r.URL.Path,
		}
		attrs := requestLogAttrs(data, logger.Level)
		logger.Log.With(attrs...).Info("request")

		return data, nil
	}

	return nil, nil
}

// requestLogAttrs returns additional request attributes for logging
func requestLogAttrs(data *RequestLogData, level string) []any {
	k := cl.AttrKey

	attrs := []any{
		slog.String(k.IP, data.ClientIP),
		slog.String(k.Method, data.Method),
		slog.String(k.Path, data.Path),
	}

	if level == cl.LevelDebug {
		if data.Headers != nil {
			attrs = append(attrs, k.Headers, data.Headers)
		}
		if data.Body != nil {
			attrs = append(attrs, k.Body, data.Body)
		}
	}

	return attrs
}
