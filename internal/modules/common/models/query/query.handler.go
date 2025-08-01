package common

import (
	"net/url"
	"strings"

	"github.com/gorilla/schema"
	"github.com/jasonsites/gosk/internal/app"
)

// QueryData composes all query parameters into a single struct for use across the app
type QueryData[T SortableEntry] struct {
	Filter *FilterQuery `schema:"filter" json:"filter,omitempty"`
	Page   PageQuery    `schema:"page" json:"page,omitempty"`
	Sort   SortQuery[T] `schema:"sort" json:"sort,omitempty"`
}

type QueryConfig[T SortableEntry] struct {
	Defaults     *QueryDefaults[T] `validate:"required"`
	EntryFactory func() T          `validate:"required"`
}

type QueryDefaults[T SortableEntry] struct {
	Page PageQuery    `validate:"required"`
	Sort SortQuery[T] `validate:"required"`
}

type QueryHandler[T SortableEntry] struct {
	defaults     *QueryDefaults[T]
	entryFactory func() T
}

func NewQueryHandler[T SortableEntry](c *QueryConfig[T]) (*QueryHandler[T], error) {
	if err := app.Validator.Validate.Struct(c); err != nil {
		return nil, err
	}

	if c.Defaults.Page.Offset == nil {
		offset := 0
		c.Defaults.Page.Offset = &offset
	}

	handler := &QueryHandler[T]{
		defaults:     c.Defaults,
		entryFactory: c.EntryFactory,
	}

	return handler, nil
}

func (q *QueryHandler[T]) ParseQuery(qs []byte) *QueryData[T] {
	data := &QueryData[T]{}
	queryString := string(qs)

	// Check if we have the deeply nested bracket notation for sort
	if strings.Contains(queryString, "sort[") && strings.Contains(queryString, "][") {
		// Use custom parser for bracket notation
		sortQuery, err := ParseDeepNestedQuery(queryString, q.entryFactory)
		if err == nil {
			data.Sort = sortQuery
		} else {
			// Fall back to defaults if parsing fails
			data.Sort = q.defaults.Sort
		}
	} else {
		// Parse the query string into url.Values for standard parsing
		values, err := url.ParseQuery(queryString)
		if err != nil {
			// If parsing fails, return data with defaults
			data.Page = q.normalizePage(data.Page)
			data.Sort = q.normalizeSort(data.Sort)
			return data
		}

		// Create a new decoder
		decoder := schema.NewDecoder()

		// Decode the values into our struct
		if err := decoder.Decode(data, values); err != nil {
			// If decoding fails, return data with defaults
			data.Page = q.normalizePage(data.Page)
			data.Sort = q.normalizeSort(data.Sort)
			return data
		}
	}

	// TODO: validate query
	// if err := app.Validator.Validate.Struct(data); err != nil {
	// 	return nil, err
	// }
	data.Page = q.normalizePage(data.Page)
	data.Sort = q.normalizeSort(data.Sort)

	return data
}

func (q *QueryHandler[T]) normalizePage(p PageQuery) PageQuery {
	page := PageQuery{
		Limit:  q.defaults.Page.Limit,
		Offset: q.defaults.Page.Offset,
	}

	if p.Limit != nil {
		page.Limit = p.Limit
	}
	if p.Offset != nil {
		page.Offset = p.Offset
	}

	return page
}

func (q *QueryHandler[T]) normalizeSort(s SortQuery[T]) SortQuery[T] {
	// If no sort query provided, use defaults
	if len(s) == 0 {
		return q.defaults.Sort
	}

	// Return the provided sort query as-is since validation happens elsewhere
	return s
}
