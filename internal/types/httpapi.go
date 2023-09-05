package types

import "github.com/google/uuid"

type Map map[string]any

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

	// Token    string `in:"query=access_token;header=x-access-token"`
	// Page     int    `in:"query=page;default=1"`
	// PerPage  int    `in:"query=per_page;default=20"`
	// IsMember bool   `in:"query=is_member"`
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

// APIMetadata
type APIMetadata struct {
	Paging PageMetadata `json:"paging,omitempty"`
}

// ListMeta
type ListMetadata struct {
	Paging PageMetadata
}

// Resource
type ResponseResource struct {
	Type       string            `json:"type"`
	ID         uuid.UUID         `json:"id"`
	Meta       *ResourceMetadata `json:"meta,omitempty"`
	Attributes any               `json:"attributes"` // TODO
}

// ResourceMetadata
type ResourceMetadata struct{}

// JSONResponseSolo
type JSONResponseSolo struct {
	Meta *APIMetadata     `json:"meta,omitempty"`
	Data ResponseResource `json:"data"`
}

// JSONResponseMult
type JSONResponseMult struct {
	Meta *APIMetadata       `json:"meta"`
	Data []ResponseResource `json:"data"`
}
