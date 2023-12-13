package jsonapi

import (
	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/core/pagination"
)

// Response
type Response[T any] struct {
	Meta *ResponseMetadata `json:"meta"`
	Data T                 `json:"data"`
}

// ResponseMetadata
type ResponseMetadata struct {
	Paging pagination.PageMetadata `json:"paging,omitempty"`
}

// Resource
type ResponseResource[T any] struct {
	Type       string            `json:"type"`
	ID         uuid.UUID         `json:"id"`
	Meta       *ResourceMetadata `json:"meta,omitempty"`
	Attributes T                 `json:"attributes"`
}

// ResourceMetadata
type ResourceMetadata struct{}
