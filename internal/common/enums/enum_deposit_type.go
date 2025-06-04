package enums

import (
	"database/sql/driver"
	"fmt"
)

type DepositType int

const (
	DepositTypeInitialMargin DepositType = iota + 1
	DepositTypeNormalDeposit
)

// String - Creating common behavior - give the type a String function
func (dt DepositType) String() string {
	return [...]string{"InitialMargin", "NormalDeposit", "CreditIn"}[dt-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (dt DepositType) EnumIndex() int {
	return int(dt)
}

// GormDataType defines the database type for GORM
func (dt DepositType) GormDataType() string {
	return "int"
}

// Value to convert CustomDate back to time.Time for database insertion
func (dt DepositType) Value() (driver.Value, error) {
	return dt.EnumIndex(), nil
}

func (dt *DepositType) Scan(value interface{}) error {
	// log.Println(value)

	// First, assert the correct type, which is likely int64 when scanning from SQL
	switch v := value.(type) {
	case int64:
		*dt = DepositType(v)
	case float64:
		*dt = DepositType(int(v))
	case int:
		*dt = DepositType(v)
	default:
		return fmt.Errorf("failed to scan DepositType, unsupported type: %T", value)
	}

	return nil
}
