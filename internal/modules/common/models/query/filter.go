package common

// FilterMetadata defines the filter-related response query parameters
type FilterMetadata struct {
	Title *string `schema:"title"`
}

// FilterQuery defines the filter-related request query paramaters
// filtered[title]=test
type FilterQuery struct {
	Title *string `schema:"title" json:"title,omitempty"`
}
