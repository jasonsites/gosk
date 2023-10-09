package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	cl "github.com/jasonsites/gosk/internal/core/logger"
	"github.com/jasonsites/gosk/internal/core/trace"
	"github.com/jasonsites/gosk/internal/core/validation"
)

// ExtendedResponseWriter extends http.ResponseWriter with a bytes.Buffer to capture the response body
type ExtendedResponseWriter struct {
	http.ResponseWriter
	BodyLogBuffer *bytes.Buffer
}

// Write extends ResponseWriter.Write by first capturing response body to the ExtendedResponseWriter.BodyLogBuffer
func (erw *ExtendedResponseWriter) Write(b []byte) (int, error) {
	erw.BodyLogBuffer.Write(b)
	return erw.ResponseWriter.Write(b)
}

// ResponseLogData defines the data captured for response logging
type ResponseLogData struct {
	Body         map[string]any
	BodySize     *int
	Headers      http.Header
	ResponseTime string
	Status       int
}

// ResponseLoggerConfig defines necessary components for the response logger middleware
type ResponseLoggerConfig struct {
	Logger *cl.CustomLogger
	Next   func(r *http.Request) bool
}

// ResponseLogger returns the response logger middleware
func ResponseLogger(c *ResponseLoggerConfig) func(http.Handler) http.Handler {
	if err := validation.Validate.Struct(c); err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c.Next != nil && c.Next(r) {
				next.ServeHTTP(w, r)
				return
			}

			// mark response time start
			start := time.Now()

			// extend default response writer
			erw := &ExtendedResponseWriter{
				BodyLogBuffer:  bytes.NewBufferString(""),
				ResponseWriter: w,
			}
			w = erw

			// call next middleware
			next.ServeHTTP(w, r)

			// calc response time and set header
			elapsed := time.Since(start).Milliseconds()
			rt := fmt.Sprintf("%sms", strconv.FormatInt(int64(elapsed), 10))
			erw.Header().Set("X-Response-Time", rt)

			if err := logResponse(erw, r, rt, c.Logger); err != nil {
				c.Logger.Log.Error(err.Error())
			}
		})
	}
}

// logResponse logs the captured response data
func logResponse(erw *ExtendedResponseWriter, r *http.Request, rt string, logger *cl.CustomLogger) error {
	if logger.Enabled {
		traceID := trace.GetTraceIDFromContext(r.Context())
		logger.Log = logger.CreateContextLogger(traceID)

		bodyBytes := erw.BodyLogBuffer.Bytes()
		bodySize := len(bodyBytes)

		data := &ResponseLogData{
			Headers:      erw.Header(),
			ResponseTime: rt,
			Status:       200, // TODO: get from response
		}

		var body map[string]any
		if bodySize > 0 {
			if err := json.Unmarshal(bodyBytes, &body); err != nil {
				return err
			}

			data.Body = body
			data.BodySize = &bodySize
		}

		attrs := responseLogAttrs(data, logger.Level)
		logger.Log.With(attrs...).Info("response")
	}

	return nil
}

// responseLogAttrs returns additional response attributes for logging
func responseLogAttrs(data *ResponseLogData, level string) []any {
	k := cl.AttrKey

	attrs := []any{
		slog.Int(k.HTTP.Status, data.Status),
		slog.String(k.ResponseTime, data.ResponseTime),
	}

	if data.BodySize != nil {
		attrs = append(attrs, slog.Int(k.HTTP.BodySize, *data.BodySize))
	}

	if level == cl.LevelDebug {
		if data.Body != nil {
			attrs = append(attrs, k.HTTP.Body, data.Body)
		}
		if data.Headers != nil {
			attrs = append(attrs, k.HTTP.Headers, data.Headers)
		}
	}

	return attrs
}
