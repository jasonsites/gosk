package common

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// SortOrder represents valid sort order values
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// IsValid checks if the sort order is valid
func (s SortOrder) IsValid() bool {
	return s == SortOrderAsc || s == SortOrderDesc
}

// String returns the string representation
func (s SortOrder) String() string {
	return string(s)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (s *SortOrder) UnmarshalText(text []byte) error {
	order := SortOrder(string(text))
	if !order.IsValid() {
		return fmt.Errorf("invalid sort order: %s, must be one of: %v", string(text), ValidSortOrders())
	}
	*s = order
	return nil
}

// ValidSortOrders returns all valid sort order options
func ValidSortOrders() []SortOrder {
	return []SortOrder{SortOrderAsc, SortOrderDesc}
}

// SortEntry represents a single SortQuery element with defined sortable fields
type SortEntry struct {
	CreatedOn  *SortOrder `schema:"created_on" json:"created_on,omitempty"`
	ModifiedOn *SortOrder `schema:"modified_on" json:"modified_on,omitempty"`
	Title      *SortOrder `schema:"title" json:"title,omitempty"`
}

// SortMetadata defines the sorting-related response metadata
// Currently mirrors SortQuery structure but kept separate for potential changes
type SortMetadata []SortEntry

// SortQuery represents an array of sort entries
type SortQuery []SortEntry

// GetSortPairs returns field-order pairs for database queries
func (sq SortQuery) GetSortPairs() []struct {
	Field string
	Order SortOrder
} {
	var pairs []struct {
		Field string
		Order SortOrder
	}

	for _, entry := range sq {
		// Check each field in the struct
		if entry.CreatedOn != nil {
			pairs = append(pairs, struct {
				Field string
				Order SortOrder
			}{Field: "created_on", Order: *entry.CreatedOn})
		}
		if entry.ModifiedOn != nil {
			pairs = append(pairs, struct {
				Field string
				Order SortOrder
			}{Field: "modified_on", Order: *entry.ModifiedOn})
		}
		if entry.Title != nil {
			pairs = append(pairs, struct {
				Field string
				Order SortOrder
			}{Field: "title", Order: *entry.Title})
		}
	}

	return pairs
}

// ToSortMetadata
func (sq SortQuery) ToSortMetadata() SortMetadata {
	return SortMetadata(sq)
}

// ValidateWithBaseFields validates the sort query with base field restrictions
func (sq *SortQuery) Validate() error {
	if err := sq.ValidateDuplicates(); err != nil {
		return err
	}

	// With struct-based SortEntry, validation is now compile-time enforced
	// Only need to validate that at least one field is set per entry
	for i, entry := range *sq {
		hasField := entry.CreatedOn != nil || entry.ModifiedOn != nil || entry.Title != nil
		if !hasField {
			return fmt.Errorf("sort entry at index %d has no fields set", i)
		}
	}

	return nil
}

// Validate validates the sort query
func (sq SortQuery) ValidateDuplicates() error {
	seen := make(map[string]bool)

	// Check for duplicate fields across entries
	for _, entry := range sq {
		// Check each field in the struct
		if entry.CreatedOn != nil {
			if seen["created_on"] {
				return fmt.Errorf("duplicate sort field: created_on")
			}
			seen["created_on"] = true
			if !entry.CreatedOn.IsValid() {
				return fmt.Errorf("invalid sort order for field created_on: %s, must be one of: %v",
					*entry.CreatedOn, ValidSortOrders())
			}
		}
		if entry.ModifiedOn != nil {
			if seen["modified_on"] {
				return fmt.Errorf("duplicate sort field: modified_on")
			}
			seen["modified_on"] = true
			if !entry.ModifiedOn.IsValid() {
				return fmt.Errorf("invalid sort order for field modified_on: %s, must be one of: %v",
					*entry.ModifiedOn, ValidSortOrders())
			}
		}
		if entry.Title != nil {
			if seen["title"] {
				return fmt.Errorf("duplicate sort field: title")
			}
			seen["title"] = true
			if !entry.Title.IsValid() {
				return fmt.Errorf("invalid sort order for field title: %s, must be one of: %v",
					*entry.Title, ValidSortOrders())
			}
		}
	}

	return nil
}

// DefaultSortQuery returns the default sort configuration
func DefaultSortQuery() SortQuery {
	return SortQuery{
		{ModifiedOn: &[]SortOrder{SortOrderDesc}[0]},
	}
}

// ParseSortQuery parses a sort query from string format
// Example: "modified_on:desc,created_on:asc"
func ParseSortQuery(sortStr string) (SortQuery, error) {
	if sortStr == "" {
		return DefaultSortQuery(), nil
	}

	var query SortQuery
	parts := strings.Split(sortStr, ",")

	for _, part := range parts {
		fieldOrder := strings.Split(strings.TrimSpace(part), ":")
		if len(fieldOrder) != 2 {
			return nil, fmt.Errorf("invalid sort format: %s", part)
		}

		field := strings.TrimSpace(fieldOrder[0])
		orderStr := strings.TrimSpace(fieldOrder[1])
		order := SortOrder(orderStr)

		if !order.IsValid() {
			return nil, fmt.Errorf("invalid sort order: %s", orderStr)
		}

		var entry SortEntry
		switch field {
		case "created_on":
			entry.CreatedOn = &order
		case "modified_on":
			entry.ModifiedOn = &order
		case "title":
			entry.Title = &order
		default:
			return nil, fmt.Errorf("invalid sort field: %s", field)
		}

		query = append(query, entry)
	}

	return query, nil
}

// ParseDeepNestedQuery handles deeply nested query parameters like sort[0][field]=value
// and converts them to a SortQuery structure
func ParseDeepNestedQuery(queryString string) (SortQuery, error) {
	if queryString == "" {
		return DefaultSortQuery(), nil
	}

	// Parse the query string
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query string: %w", err)
	}

	// Map to collect sort entries by index
	entriesMap := make(map[int]*SortEntry)

	// Process each key-value pair
	for key, valueSlice := range values {
		if len(valueSlice) == 0 {
			continue
		}
		value := valueSlice[0] // Take the first value

		// Check if it's a sort parameter with bracket notation
		if strings.HasPrefix(key, "sort[") {
			index, field, err := parseNestedKey(key)
			if err != nil {
				continue // Skip invalid keys
			}

			// Parse the sort order
			order := SortOrder(value)
			if !order.IsValid() {
				return nil, fmt.Errorf("invalid sort order: %s", value)
			}

			// Get or create the entry for this index
			if entriesMap[index] == nil {
				entriesMap[index] = &SortEntry{}
			}

			// Set the appropriate field
			switch field {
			case "created_on":
				entriesMap[index].CreatedOn = &order
			case "modified_on":
				entriesMap[index].ModifiedOn = &order
			case "title":
				entriesMap[index].Title = &order
			default:
				return nil, fmt.Errorf("invalid sort field: %s", field)
			}
		}
	}

	// Convert map to slice, maintaining order
	var result SortQuery
	for i := 0; i < len(entriesMap); i++ {
		if entry, exists := entriesMap[i]; exists {
			result = append(result, *entry)
		}
	}

	return result, nil
}

// parseNestedKey parses a key like "sort[0][created_on]" and returns index and field
func parseNestedKey(key string) (int, string, error) {
	// Remove "sort[" prefix
	if !strings.HasPrefix(key, "sort[") {
		return 0, "", fmt.Errorf("invalid sort key format")
	}

	remaining := key[5:] // Remove "sort["

	// Find the first closing bracket
	firstClose := strings.Index(remaining, "]")
	if firstClose == -1 {
		return 0, "", fmt.Errorf("invalid sort key format: missing closing bracket")
	}

	// Extract index
	indexStr := remaining[:firstClose]
	index := 0
	if indexStr != "" {
		var err error
		index, err = strconv.Atoi(indexStr)
		if err != nil {
			return 0, "", fmt.Errorf("invalid index: %s", indexStr)
		}
	}

	// Extract field name
	remaining = remaining[firstClose+1:] // Skip "]"
	if !strings.HasPrefix(remaining, "[") || !strings.HasSuffix(remaining, "]") {
		return 0, "", fmt.Errorf("invalid sort key format: invalid field brackets")
	}

	field := remaining[1 : len(remaining)-1] // Remove "[" and "]"

	return index, field, nil
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
func (se SortEntry) GetActiveFields() map[string]SortOrder {
	fields := make(map[string]SortOrder)
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
