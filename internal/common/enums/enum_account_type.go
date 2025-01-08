package enums

import (
	"database/sql/driver"
	"fmt"
)

type AccountType int

const (
	AccountTypeReal AccountType = iota + 1
	AccountTypeDemo
)

// String - Creating common behavior - give the type a String function
func (at AccountType) String() string {
	return [...]string{"Real", "Demo"}[at-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (at AccountType) EnumIndex() int {
	return int(at)
}

// GormDataType defines the database type for GORM
func (at AccountType) GormDataType() string {
	return "int"
}

// Value to convert CustomDate back to time.Time for database insertion
func (at AccountType) Value() (driver.Value, error) {
	return at.EnumIndex(), nil
}

func (at *AccountType) Scan(value interface{}) error {
	// log.Println(value)

	// First, assert the correct type, which is likely int64 when scanning from SQL
	switch v := value.(type) {
	case int64:
		*at = AccountType(v)
	case float64:
		*at = AccountType(int(v))
	case int:
		*at = AccountType(v)
	default:
		return fmt.Errorf("failed to scan AccountType, unsupported type: %T", value)
	}

	return nil
}
