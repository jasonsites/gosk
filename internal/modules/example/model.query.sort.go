package example

import (
	"fmt"

	q "github.com/jasonsites/gosk/internal/modules/common/models/query"
)

// ExampleSortMetadata defines the Example sorting-related response metadata
// Currently mirrors SortQuery structure but kept separate for potential changes
type ExampleSortMetadata q.SortMetadata[SortEntry]

// ExampleSortQuery represents an array of sort entries
type ExampleSortQuery q.SortQuery[SortEntry]

// SortEntry represents a single SortQuery element with defined sortable fields
type SortEntry struct {
	CreatedOn  *q.SortOrder `schema:"created_on" json:"created_on,omitempty"`
	ModifiedOn *q.SortOrder `schema:"modified_on" json:"modified_on,omitempty"`
	Title      *q.SortOrder `schema:"title" json:"title,omitempty"`
}

// GetFieldCount returns the number of non-nil fields in the entry
func (se SortEntry) GetFieldCount() int {
	count := 0
	if se.CreatedOn != nil {
		count++
	}
	if se.ModifiedOn != nil {
		count++
	}
	if se.Title != nil {
		count++
	}
	return count
}

// HasAnyField returns true if at least one field is set
func (se SortEntry) HasAnyField() bool {
	return se.CreatedOn != nil || se.ModifiedOn != nil || se.Title != nil
}

// GetActiveFields returns a map of active field names and their orders
func (se SortEntry) GetActiveFields() map[string]q.SortOrder {
	fields := make(map[string]q.SortOrder)
	if se.CreatedOn != nil {
		fields["created_on"] = *se.CreatedOn
	}
	if se.ModifiedOn != nil {
		fields["modified_on"] = *se.ModifiedOn
	}
	if se.Title != nil {
		fields["title"] = *se.Title
	}
	return fields
}

// GetSortPairs returns field-order pairs for database queries
func (se SortEntry) GetSortPairs() []struct {
	Field string
	Order q.SortOrder
} {
	var pairs []struct {
		Field string
		Order q.SortOrder
	}

	if se.CreatedOn != nil {
		pairs = append(pairs, struct {
			Field string
			Order q.SortOrder
		}{Field: "created_on", Order: *se.CreatedOn})
	}
	if se.ModifiedOn != nil {
		pairs = append(pairs, struct {
			Field string
			Order q.SortOrder
		}{Field: "modified_on", Order: *se.ModifiedOn})
	}
	if se.Title != nil {
		pairs = append(pairs, struct {
			Field string
			Order q.SortOrder
		}{Field: "title", Order: *se.Title})
	}

	return pairs
}

// SetFieldFromString sets a field by name from a string value
func (se SortEntry) SetFieldFromString(fieldName string, order q.SortOrder) (q.SortableEntry, error) {
	switch fieldName {
	case "created_on":
		se.CreatedOn = &order
	case "modified_on":
		se.ModifiedOn = &order
	case "title":
		se.Title = &order
	default:
		return se, fmt.Errorf("invalid field name: %s", fieldName)
	}
	return se, nil
}

// GetValidFieldNames returns a list of valid field names for this entry type
func (se SortEntry) GetValidFieldNames() []string {
	return []string{"created_on", "modified_on", "title"}
}
