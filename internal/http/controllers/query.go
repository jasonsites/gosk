package controllers

import (
	"github.com/hetiansu5/urlquery"
	"github.com/jasonsites/gosk-api/internal/core/query"
)

func parseQuery(qs []byte) *query.QueryData {
	data := &query.QueryData{}
	urlquery.Unmarshal(qs, data)
	data.Paging = pageSettings(data.Paging)
	data.Sorting = sortSettings(data.Sorting)
	return data
}

func pageSettings(p query.QueryPaging) query.QueryPaging {
	var (
		defaultLimit  = 20 // TODO: move to config
		defaultOffset = 0  // TODO: move to config
	)

	page := query.QueryPaging{
		Limit:  &defaultLimit,
		Offset: &defaultOffset,
	}

	if p.Limit != nil {
		page.Limit = p.Limit
	}
	if p.Offset != nil {
		page.Offset = p.Offset
	}

	return page
}

func sortSettings(s query.QuerySorting) query.QuerySorting {
	var (
		defaultOrder = "desc"       // TODO: move to config
		defaultProp  = "created_on" // TODO: move to config
	)

	sort := query.QuerySorting{
		Order: &defaultOrder,
		Prop:  &defaultProp,
	}

	if s.Order != nil {
		sort.Order = s.Order
	}
	if s.Prop != nil {
		sort.Prop = s.Prop
	}

	return sort
}
