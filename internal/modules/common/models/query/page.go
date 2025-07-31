package common

// PageMetadata defines the paging-related response metadata
type PageMetadata struct {
	Limit  uint32 `json:"limit"`
	Offset uint32 `json:"offset"`
	Total  uint32 `json:"total"`
}

// PageQuery defines the paging-related query paramaters
// p[limit]=20&p[offset]=10
type PageQuery struct {
	Limit  *int `schema:"limit" json:"limit,omitempty"`
	Offset *int `schema:"offset" json:"offset,omitempty"`
}
