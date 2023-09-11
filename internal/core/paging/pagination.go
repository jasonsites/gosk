package paging

type PageMetadata struct {
	Limit  uint32
	Offset uint32
	Total  uint32
}

type SortMetadata struct {
	Order string // "asc" || "desc"
	Prop  string
}
