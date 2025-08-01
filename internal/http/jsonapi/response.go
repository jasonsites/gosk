package jsonapi

import (
	"github.com/google/uuid"
	query "github.com/jasonsites/gosk/internal/modules/common/models/query"
)

// Response
type Response struct {
	Meta *ResponseMetadata `json:"meta"`
	Data any               `json:"data"`
}

// ResponseMetadata
type ResponseMetadata struct {
	Filter *query.FilterMetadata `json:"filter,omitempty"`
	Page   query.PageMetadata    `json:"page,omitempty"`
	Sort   any                   `json:"sort,omitempty"`
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
