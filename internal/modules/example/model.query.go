package example

import (
	q "github.com/jasonsites/gosk/internal/modules/common/models/query"
)

type ExampleQueryData q.QueryData[SortEntry]

type ExampleQueryHandler q.QueryHandler[SortEntry]
