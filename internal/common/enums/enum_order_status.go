package enums

import (
	"database/sql/driver"
	"fmt"
)

type OrderStatus int

const (
	OrderStatusApproved OrderStatus = iota + 1
	OrderStatusRejected
	OrderStatusPending
	OrderStatusCancelled
	OrderStatusNew
	OrderStatusClosed
)

// String - Creating common behavior - give the type a String function
func (aas OrderStatus) String() string {
	return [...]string{"Approved", "Rejected", "Pending", "Cancelled", "New", "Closed"}[aas-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (aas OrderStatus) EnumIndex() int {
	return int(aas)
}

// GormDataType defines the database type for GORM
func (aas OrderStatus) GormDataType() string {
	return "int"
}

// Value to convert CustomDate back to time.Time for database insertion
func (aas OrderStatus) Value() (driver.Value, error) {
	return aas.EnumIndex(), nil
}

func (aas *OrderStatus) Scan(value interface{}) error {
	// log.Println(value)

	// First, assert the correct type, which is likely int64 when scanning from SQL
	switch v := value.(type) {
	case int64:
		*aas = OrderStatus(v)
	case float64:
		*aas = OrderStatus(int(v))
	case int:
		*aas = OrderStatus(v)
	default:
		return fmt.Errorf("failed to scan OrderStatus, unsupported type: %T", value)
	}

	return nil
}
