package common

import (
	"database/sql/driver"
	"fmt"
)

// RecordStatus represents the record_status enum from the database
type RecordStatus string

const (
	RecordStatusActive   RecordStatus = "active"
	RecordStatusArchived RecordStatus = "archived"
	RecordStatusDeleted  RecordStatus = "deleted"
)

// Value implements the driver.Valuer interface for database operations
func (rs RecordStatus) Value() (driver.Value, error) {
	return string(rs), nil
}

// Scan implements the sql.Scanner interface for database operations
func (rs *RecordStatus) Scan(value interface{}) error {
	if value == nil {
		*rs = ""
		return nil
	}
	if str, ok := value.(string); ok {
		*rs = RecordStatus(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into RecordStatus", value)
}
