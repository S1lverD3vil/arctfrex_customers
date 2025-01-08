package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type RoleIdType int

const (
	AdminBackoffice RoleIdType = iota + 1
	SubIB
	HeadMarketing
	TeamLeader
	SalesManager
)

// // roleMap maps RoleIdType to string for reverse lookup
// var roleMap = map[RoleIdType]string{
// 	AdminBackoffice: "AdminBackoffice",
// 	SubIB:           "SubIB",
// 	HeadMarketing:   "HeadMarketing",
// 	TeamLeader:      "TeamLeader",
// 	SalesManager:    "SalesManager",
// }

// roleStringMap maps string to RoleIdType for conversion
var roleStringMap = map[string]RoleIdType{
	"AdminBackoffice": AdminBackoffice,
	"SubIB":           SubIB,
	"HeadMarketing":   HeadMarketing,
	"TeamLeader":      TeamLeader,
	"SalesManager":    SalesManager,
}

// String - Creating common behavior - give the type a String function
func (rit RoleIdType) String() string {
	return [...]string{"AdminBackoffice", "SubIB", "HeadMarketing", "TeamLeader", "SalesManager"}[rit-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex functio
func (rit RoleIdType) EnumIndex() int {
	return int(rit)
}

// GormDataType defines the database type for GORM
func (rit RoleIdType) GormDataType() string {
	return "int"
}

// Value to convert CustomDate back to time.Time for database insertion
func (rit RoleIdType) Value() (driver.Value, error) {
	return rit.EnumIndex(), nil
}

func (rit *RoleIdType) Scan(value interface{}) error {
	// log.Println(value)

	// First, assert the correct type, which is likely int64 when scanning from SQL
	switch v := value.(type) {
	case int64:
		*rit = RoleIdType(v)
	case float64:
		*rit = RoleIdType(int(v))
	case int:
		*rit = RoleIdType(v)
	case string:
		// Check if the string is a valid role name
		if role, exists := roleStringMap[v]; exists {
			*rit = role
			return nil
		}
		return fmt.Errorf("invalid RoleIdType string: %s", v)
	default:
		return fmt.Errorf("failed to scan JobPositionType, unsupported type: %T", value)
	}

	return nil
}

// Implement UnmarshalJSON to allow unmarshalling a string to RoleIdType
func (rit *RoleIdType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Look up the string in the roleStringMap to get the corresponding RoleIdType
	if role, exists := roleStringMap[str]; exists {
		*rit = role
		return nil
	}

	// If string is not valid, return an error
	return fmt.Errorf("invalid RoleIdType string: %s", str)
}
