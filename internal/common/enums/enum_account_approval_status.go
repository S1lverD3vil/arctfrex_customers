package enums

import (
	"database/sql/driver"
	"fmt"
)

type AccountApprovalStatus int

const (
	AccountApprovalStatusApproved AccountApprovalStatus = iota + 1
	AccountApprovalStatusRejected
	AccountApprovalStatusPending
	AccountApprovalStatusCancelled
)

var ApprovalStatusMap = map[string]AccountApprovalStatus{
	"approved":  AccountApprovalStatusApproved,
	"rejected":  AccountApprovalStatusRejected,
	"pending":   AccountApprovalStatusPending,
	"cancelled": AccountApprovalStatusCancelled,
}

// String - Creating common behavior - give the type a String function
func (aas AccountApprovalStatus) String() string {
	return [...]string{"Approved", "Rejected", "Pending", "Cancelled"}[aas-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (aas AccountApprovalStatus) EnumIndex() int {
	return int(aas)
}

// GormDataType defines the database type for GORM
func (aas AccountApprovalStatus) GormDataType() string {
	return "int"
}

// Value to convert CustomDate back to time.Time for database insertion
func (aas AccountApprovalStatus) Value() (driver.Value, error) {
	return aas.EnumIndex(), nil
}

func (aas *AccountApprovalStatus) Scan(value interface{}) error {
	// log.Println(value)

	// First, assert the correct type, which is likely int64 when scanning from SQL
	switch v := value.(type) {
	case int64:
		*aas = AccountApprovalStatus(v)
	case float64:
		*aas = AccountApprovalStatus(int(v))
	case int:
		*aas = AccountApprovalStatus(v)
	default:
		return fmt.Errorf("failed to scan AccountApprovalStatus, unsupported type: %T", value)
	}

	return nil
}
