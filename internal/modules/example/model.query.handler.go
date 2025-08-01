package example

import (
	q "github.com/jasonsites/gosk/internal/modules/common/models/query"
)

// NewExampleQueryHandler creates a new query handler for the Example module
func NewExampleQueryHandler(c *q.QueryConfig[SortEntry]) (*ExampleQueryHandler, error) {
	handler, err := q.NewQueryHandler(c)
	if err != nil {
		return nil, err
	}
	return (*ExampleQueryHandler)(handler), nil
}

// CreateSortEntry is a factory function for creating new SortEntry instances
func CreateSortEntry() SortEntry {
	return SortEntry{}
}

// DefaultExampleSortQuery returns the default sort configuration for Example
func DefaultExampleSortQuery() q.SortQuery[SortEntry] {
	desc := q.SortOrderDesc
	return q.SortQuery[SortEntry]{
		SortEntry{ModifiedOn: &desc},
	}
}

// ParseQuery parses query parameters for the Example module
func (h *ExampleQueryHandler) ParseQuery(qs []byte) *ExampleQueryData {
	result := (*q.QueryHandler[SortEntry])(h).ParseQuery(qs)
	return (*ExampleQueryData)(result)
}
