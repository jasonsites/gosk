package jsonapi

import (
	"github.com/google/uuid"
	"github.com/jasonsites/gosk-api/internal/core/paging"
)

// ListMeta
type ListMetadata struct {
	Paging paging.PageMetadata `json:"paging,omitempty"`
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

// // ResponseSolo
// type ResponseSolo struct {
// 	Meta *ListMetadata    `json:"meta,omitempty"`
// 	Data ResponseResource `json:"data"`
// }

// // ResponseMult
// type ResponseMult struct {
// 	Meta *ListMetadata      `json:"meta"`
// 	Data []ResponseResource `json:"data"`
// }

type Response struct {
	Meta *ListMetadata `json:"meta"`
	Data any           `json:"data"`
}
