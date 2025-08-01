package common

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// =====================================================================================================================
// SortOrder
// =====================================================================================================================

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

// =====================================================================================================================
// SortQuery
// =====================================================================================================================

// SortableEntry defines the interface that sortable entries must implement
type SortableEntry interface {
	GetFieldCount() int
	HasAnyField() bool
	GetActiveFields() map[string]SortOrder
	GetSortPairs() []struct {
		Field string
		Order SortOrder
	}
	// SetFieldFromString sets a field by name from a string value
	// Returns the modified entry and error if field name is invalid
	SetFieldFromString(fieldName string, order SortOrder) (SortableEntry, error)
	// GetValidFieldNames returns a list of valid field names for this entry type
	GetValidFieldNames() []string
}

// SortMetadata defines the sorting-related response metadata
// Currently mirrors SortQuery structure but kept separate for potential changes
type SortMetadata[T SortableEntry] []T

// SortQuery represents an array of sort entries
type SortQuery[T SortableEntry] []T

// GetSortPairs returns field-order pairs for database queries
func (sq SortQuery[T]) GetSortPairs() []struct {
	Field string
	Order SortOrder
} {
	var pairs []struct {
		Field string
		Order SortOrder
	}

	for _, entry := range sq {
		pairs = append(pairs, entry.GetSortPairs()...)
	}

	return pairs
}

// ToSortMetadata converts SortQuery to SortMetadata
func (sq SortQuery[T]) ToSortMetadata() SortMetadata[T] {
	return SortMetadata[T](sq)
}

// ValidateWithBaseFields validates the sort query with base field restrictions
func (sq *SortQuery[T]) Validate() error {
	if err := sq.ValidateDuplicates(); err != nil {
		return err
	}

	// Validate that at least one field is set per entry
	for i, entry := range *sq {
		if !entry.HasAnyField() {
			return fmt.Errorf("sort entry at index %d has no fields set", i)
		}
	}

	return nil
}

// ValidateDuplicates validates the sort query for duplicate fields
func (sq SortQuery[T]) ValidateDuplicates() error {
	seen := make(map[string]bool)

	// Check for duplicate fields across entries
	for _, entry := range sq {
		for field, order := range entry.GetActiveFields() {
			if seen[field] {
				return fmt.Errorf("duplicate sort field: %s", field)
			}
			seen[field] = true
			if !order.IsValid() {
				return fmt.Errorf("invalid sort order for field %s: %s, must be one of: %v",
					field, order, ValidSortOrders())
			}
		}
	}

	return nil
}

// =====================================================================================================================
// Generic Query Parsing
// =====================================================================================================================

// ParseDeepNestedQuery handles deeply nested query parameters like sort[0][field]=value
// and converts them to a SortQuery structure for any SortableEntry type
func ParseDeepNestedQuery[T SortableEntry](queryString string, createEntry func() T) (SortQuery[T], error) {
	if queryString == "" {
		return SortQuery[T]{}, nil
	}

	// Parse the query string
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query string: %w", err)
	}

	// Map to collect sort entries by index
	entriesMap := make(map[int]T)

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
			entry, exists := entriesMap[index]
			if !exists {
				entry = createEntry()
				entriesMap[index] = entry
			}

			// Set the field using the interface method
			updatedEntry, err := entry.SetFieldFromString(field, order)
			if err != nil {
				return nil, fmt.Errorf("invalid sort field: %s", field)
			}

			// Update the entry in the map
			entriesMap[index] = updatedEntry.(T)
		}
	}

	// Convert map to slice, maintaining order
	var result SortQuery[T]
	for i := 0; i < len(entriesMap); i++ {
		if entry, exists := entriesMap[i]; exists {
			result = append(result, entry)
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
