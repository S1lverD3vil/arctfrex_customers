package enums

import (
	"database/sql/driver"
	"fmt"
)

type DepositApprovalStatus int

const (
	DepositApprovalStatusApproved DepositApprovalStatus = iota + 1
	DepositApprovalStatusRejected
	DepositApprovalStatusPending
	DepositApprovalStatusCancelled
	DepositApprovalStatusNew
)

// String - Creating common behavior - give the type a String function
func (das DepositApprovalStatus) String() string {
	return [...]string{"Approved", "Rejected", "Pending", "Cancelled", "New"}[das-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (das DepositApprovalStatus) EnumIndex() int {
	return int(das)
}

// GormDataType defines the database type for GORM
func (das DepositApprovalStatus) GormDataType() string {
	return "int"
}

// Value to convert CustomDate back to time.Time for database insertion
func (das DepositApprovalStatus) Value() (driver.Value, error) {
	return das.EnumIndex(), nil
}

func (das *DepositApprovalStatus) Scan(value interface{}) error {
	// log.Println(value)

	// First, assert the correct type, which is likely int64 when scanning from SQL
	switch v := value.(type) {
	case int64:
		*das = DepositApprovalStatus(v)
	case float64:
		*das = DepositApprovalStatus(int(v))
	case int:
		*das = DepositApprovalStatus(v)
	default:
		return fmt.Errorf("failed to scan DepositApprovalStatus, unsupported type: %T", value)
	}

	return nil
}
