package enums

import (
	"database/sql/driver"
	"fmt"
)

type OrderType int

const (
	OrderTypeBuy OrderType = iota + 1
	OrderTypeSell
)

// String - Creating common behavior - give the type a String function
func (at OrderType) String() string {
	return [...]string{"BuyNow", "Sell"}[at-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (at OrderType) EnumIndex() int {
	return int(at)
}

// GormDataType defines the database type for GORM
func (at OrderType) GormDataType() string {
	return "int"
}

// Value to convert CustomDate back to time.Time for database insertion
func (at OrderType) Value() (driver.Value, error) {
	return at.EnumIndex(), nil
}

func (at *OrderType) Scan(value interface{}) error {
	//log.Println(value)

	// First, assert the correct type, which is likely int64 when scanning from SQL
	switch v := value.(type) {
	case int64:
		*at = OrderType(v)
	case float64:
		*at = OrderType(int(v))
	case int:
		*at = OrderType(v)
	default:
		return fmt.Errorf("failed to scan AccountApprovalStatus, unsupported type: %T", value)
	}

	return nil
}
