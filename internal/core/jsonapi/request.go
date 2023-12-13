package jsonapi

// Envelope
type Envelope map[string]any

// RequestBody
type RequestBody[T any] struct {
	Data *RequestResource[T] `json:"data" validate:"required"`
}

// RequestResource
type RequestResource[T any] struct {
	Type       string `json:"type" validate:"required"`
	ID         string `json:"id" validate:"omitempty,uuid4"`
	Attributes *T     `json:"attributes" validate:"required"`
}
