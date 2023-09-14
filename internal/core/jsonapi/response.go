package jsonapi

import (
	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/core/pagination"
)

// Response
type Response struct {
	Meta *ResponseMetadata `json:"meta"`
	Data any               `json:"data"`
}

// ListMeta
type ResponseMetadata struct {
	Paging pagination.PageMetadata `json:"paging,omitempty"`
}

// Resource
type ResponseResource struct {
	Type       string            `json:"type"`
	ID         uuid.UUID         `json:"id"`
	Meta       *ResourceMetadata `json:"meta,omitempty"`
	Attributes any               `json:"attributes"`
}

// ResourceMetadata
type ResourceMetadata struct{}
