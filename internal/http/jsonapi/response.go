package jsonapi

import (
	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/modules/common/pagination"
)

// Response
type Response struct {
	Meta *ResponseMetadata `json:"meta"`
	Data any               `json:"data"`
}

// ResponseMetadata
type ResponseMetadata struct {
	Page pagination.PageMetadata `json:"page,omitempty"`
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
