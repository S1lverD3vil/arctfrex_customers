package enums

import (
	"database/sql/driver"
	"fmt"
)

type WithdrawalApprovalStatus int

const (
	WithdrawalApprovalStatusApproved WithdrawalApprovalStatus = iota + 1
	WithdrawalApprovalStatusRejected
	WithdrawalApprovalStatusPending
	WithdrawalApprovalStatusCancelled
)

// String - Creating common behavior - give the type a String function
func (was WithdrawalApprovalStatus) String() string {
	return [...]string{"Approved", "Rejected", "Pending", "Cancelled"}[was-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (was WithdrawalApprovalStatus) EnumIndex() int {
	return int(was)
}

// GormDataType defines the database type for GORM
func (was WithdrawalApprovalStatus) GormDataType() string {
	return "int"
}

// Value to convert CustomDate back to time.Time for database insertion
func (was WithdrawalApprovalStatus) Value() (driver.Value, error) {
	return was.EnumIndex(), nil
}

func (was *WithdrawalApprovalStatus) Scan(value interface{}) error {
	// log.Println(value)

	// First, assert the correct type, which is likely int64 when scanning from SQL
	switch v := value.(type) {
	case int64:
		*was = WithdrawalApprovalStatus(v)
	case float64:
		*was = WithdrawalApprovalStatus(int(v))
	case int:
		*was = WithdrawalApprovalStatus(v)
	default:
		return fmt.Errorf("failed to scan WithdrawalApprovalStatus, unsupported type: %T", value)
	}

	return nil
}
