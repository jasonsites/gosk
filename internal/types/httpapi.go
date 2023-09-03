package types

import "github.com/google/uuid"

// JSON Request -----------------------------------------------------------------------------------

// JSONRequestBody
type JSONRequestBody struct {
	Data *RequestResource `json:"data" validate:"required"`
}

// RequestResource
type RequestResource struct {
	Type       string `json:"type" validate:"required"`
	ID         string `json:"id" validate:"omitempty,uuid4"`
	Attributes any    `json:"attributes" validate:"required"`
}

// ExampleRequestData defines an Example resource for data attributes request binding
type ExampleRequestData struct {
	Deleted     bool    `json:"deleted" validate:"omitempty,boolean"`
	Description *string `json:"description" validate:"omitempty,min=3,max=999"`
	Enabled     bool    `json:"enabled"  validate:"omitempty,boolean"`
	Status      *uint32 `json:"status" validate:"omitempty,numeric"`
	Title       string  `json:"title" validate:"required,omitempty,min=2,max=255"`
}

// Query ------------------------------------------------------------------------------------------

// QueryData composes all query parameters into a single struct for use across the app
type QueryData struct {
	Filters QueryFilters `query:"f"`
	Options QueryOptions `query:"o"`
	Paging  QueryPaging  `query:"p"`
	Sorting QuerySorting `query:"s"`
}

// QueryFilters defines the filter-related query paramaters
// f[enabled]=true&f[name]=test&f[status]=4
type QueryFilters struct {
	Enabled *bool   `query:"enabled"`
	Name    *string `query:"name"`
	Status  *int    `query:"status"`
}

// QueryOptions defines the options-related query paramaters
// o[export]=true
type QueryOptions struct {
	Export *bool `query:"export"`
}

// QueryPaging defines the paging-related query paramaters
// p[limit]=20&p[offset]=10
type QueryPaging struct {
	Limit  *int `query:"limit"`
	Offset *int `query:"offset"`
}

// QuerySorting defines the sorting-related query paramaters
// s[order]=desc&s[prop]=name
type QuerySorting struct {
	Order *string `query:"order"`
	Prop  *string `query:"prop"`
}

// Response ---------------------------------------------------------------------------------------

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
