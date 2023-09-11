package query

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
