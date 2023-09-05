package types

import (
	"context"
)

// Correlation
type Correlation struct {
	Headers map[string][]string
	TraceID ContextKey // TODO: consider uuid.UUID
}

// ContextKey wraps scalar string for use in context (to avoid naming collisions)
type ContextKey string

// String implements the Stringer interface for print statements
func (c ContextKey) String() string {
	return string(c)
}

// TraceIDContextKey defines the context key used for tracking operation trace ID
const TraceIDContextKey ContextKey = "trace_id"

// CreateOpContext creates an operation context with correlation data
func CreateOpContext(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ContextKey(string(TraceIDContextKey)), traceID)
}

// GetTraceIDFromContext retrieves the trace ID from the operation context
func GetTraceIDFromContext(ctx context.Context) string {
	val := ctx.Value(TraceIDContextKey)
	traceID, ok := val.(string)
	if !ok {
		traceID = "unknown"
	}
	return traceID
}
