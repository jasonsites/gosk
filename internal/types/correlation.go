package types

// CorrelationContextKey
const CorrelationContextKey string = "Trace"

// Trace
type Trace struct {
	Headers   map[string]string
	RequestID string // TODO: consider uuid.UUID
}
