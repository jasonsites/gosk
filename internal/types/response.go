package types

import (
	"github.com/google/uuid"
)

// JSONResponse defines the interface for a JSON serialized response
type JSONResponse interface {
	Discoverable
}

// APIMetadata
type APIMetadata struct {
	Paging ListPaging `json:"paging,omitempty"`
}

// ListMeta
type ListMeta struct {
	Paging ListPaging
}

// ListPaging
type ListPaging struct {
	Limit  uint32
	Offset uint32
	Total  uint32
}

// ResourceMetadata
type ResourceMetadata struct{}

// Resource
type ResponseResource struct {
	Type       string            `json:"type"`
	ID         uuid.UUID         `json:"id"`
	Meta       *ResourceMetadata `json:"meta,omitempty"`
	Attributes any               `json:"attributes"` // TODO
}

// JSONResponseSolo
type JSONResponseSolo struct {
	Meta *APIMetadata     `json:"meta,omitempty"`
	Data ResponseResource `json:"data"`
}

// Discover
func (r *JSONResponseSolo) Discover() Discoverable {
	return r
}

// JSONResponseMult
type JSONResponseMult struct {
	Meta *APIMetadata       `json:"meta"`
	Data []ResponseResource `json:"data"`
}

// Discover
func (r *JSONResponseMult) Discover() Discoverable {
	return r
}
