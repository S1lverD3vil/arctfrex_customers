package common

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type CustomDate time.Time

// UnmarshalJSON for parsing the custom date format (YYYY-MM-DD)
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		return nil // Handle empty strings
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*cd = CustomDate(t)
	return nil
}

// MarshalJSON to convert CustomDate back to JSON format
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	t := time.Time(cd)
	if t.IsZero() {
		return []byte(`null`), nil // Handle zero time properly in JSON
	}
	return []byte(`"` + t.Format("2006-01-02") + `"`), nil
}

// GormDataType defines the database type for GORM
func (cd CustomDate) GormDataType() string {
	return "date"
}

// Value to convert CustomDate back to time.Time for database insertion
func (cd CustomDate) Value() (driver.Value, error) {
	t := time.Time(cd)
	if t.IsZero() {
		return nil, nil // Return nil to avoid saving zero time
	}
	return t, nil
}

// Scan to read from the database into CustomDate
func (cd *CustomDate) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan CustomDate")
	}
	*cd = CustomDate(t)
	return nil
}
