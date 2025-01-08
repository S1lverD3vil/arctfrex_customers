package enums

import (
	"database/sql/driver"
	"fmt"
)

type JobPositionType int

const (
	JobPositionTypeAdmin JobPositionType = iota + 1
	JobPositionTypeMarketing
	JobPositionTypeTeamLeader
	JobPositionTypeSalesManager
	JobPositionTypeMarketingGroup
)

// String - Creating common behavior - give the type a String function
func (jpt JobPositionType) String() string {
	return [...]string{"Admin", "Marketing", "TeamLeader", "SalesManager", "MarketingGroup"}[jpt-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (jpt JobPositionType) EnumIndex() int {
	return int(jpt)
}

// GormDataType defines the database type for GORM
func (jpt JobPositionType) GormDataType() string {
	return "int"
}

// Value to convert CustomDate back to time.Time for database insertion
func (jpt JobPositionType) Value() (driver.Value, error) {
	return jpt.EnumIndex(), nil
}

func (jpt *JobPositionType) Scan(value interface{}) error {
	// log.Println(value)

	// First, assert the correct type, which is likely int64 when scanning from SQL
	switch v := value.(type) {
	case int64:
		*jpt = JobPositionType(v)
	case float64:
		*jpt = JobPositionType(int(v))
	case int:
		*jpt = JobPositionType(v)
	default:
		return fmt.Errorf("failed to scan JobPositionType, unsupported type: %T", value)
	}

	return nil
}
