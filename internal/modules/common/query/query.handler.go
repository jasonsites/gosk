package query

import (
	"github.com/hetiansu5/urlquery"
	"github.com/jasonsites/gosk/internal/app"
)

type QueryConfig struct {
	Defaults *QueryDefaults `validate:"required"`
}

type QueryDefaults struct {
	Paging  QueryPaging   `validate:"required"`
	Sorting *QuerySorting `validate:"required"`
}

type QueryHandler struct {
	defaults *QueryDefaults
}

func NewQueryHandler(c *QueryConfig) (*QueryHandler, error) {
	if err := app.Validator.Validate.Struct(c); err != nil {
		return nil, err
	}

	if c.Defaults.Paging.Offset == nil {
		offset := 0
		c.Defaults.Paging.Offset = &offset
	}

	handler := &QueryHandler{
		defaults: c.Defaults,
	}

	return handler, nil
}

func (q *QueryHandler) ParseQuery(qs []byte) *QueryData {
	data := &QueryData{}

	// TODO: validate query
	urlquery.Unmarshal(qs, data)
	// if err := app.Validator.Validate.Struct(data); err != nil {
	// 	return nil, err
	// }
	data.Paging = q.pageSettings(data.Paging)
	data.Sorting = q.sortSettings(data.Sorting)

	return data
}

func (q *QueryHandler) pageSettings(p QueryPaging) QueryPaging {
	page := QueryPaging{
		Limit:  q.defaults.Paging.Limit,
		Offset: q.defaults.Paging.Offset,
	}

	if p.Limit != nil {
		page.Limit = p.Limit
	}
	if p.Offset != nil {
		page.Offset = p.Offset
	}

	return page
}

func (q *QueryHandler) sortSettings(s QuerySorting) QuerySorting {
	sort := QuerySorting{
		Order: q.defaults.Sorting.Order,
		Attr:  q.defaults.Sorting.Attr,
	}

	if s.Attr != nil {
		sort.Attr = s.Attr
	}
	if s.Order != nil {
		sort.Order = s.Order
	}

	return sort
}
