package repository

import (
	"fmt"
	"strings"
)

// BuildInsertFieldsAndValues builds the insert fields and values strings used in SQL query string builders
func BuildInsertFieldsAndValues(fields ...string) (string, string) {
	var insert, values strings.Builder

	length := len(fields)

	insert.WriteString("(")
	values.WriteString("(")

	for i, str := range fields {
		insert.WriteString(str)
		values.WriteString("$" + fmt.Sprint(i+1))
		if i < length-1 {
			insert.WriteString(",")
			values.WriteString(",")
		}
	}

	insert.WriteString(")")
	values.WriteString(")")

	return insert.String(), values.String()
}

// BuildReturnFields builds the return fields string used in SQL query string builders
func BuildReturnFields(fields ...string) string {
	var returnFields strings.Builder

	length := len(fields)

	for i, str := range fields {
		returnFields.WriteString(str)
		if i < length-1 {
			returnFields.WriteString(",")
		}
	}

	return returnFields.String()
}

// BuildUpdateValues builds the update values string used in SQL query string builders
func BuildUpdateValues(fields ...string) string {
	var values strings.Builder

	length := len(fields)

	for i, str := range fields {
		values.WriteString(str + "=$" + fmt.Sprint(i+1))
		if i < length-1 {
			values.WriteString(",")
		}
	}

	return values.String()
}
