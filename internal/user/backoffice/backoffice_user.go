package user

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/grouprole"
	"arctfrex-customers/internal/role"

	"gorm.io/gorm"
)

type BackofficeUsers struct {
	ID           string                `gorm:"primary_key" json:"userid"`
	Name         string                `json:"customer_name"`
	Email        string                `json:"email" binding:"required,email"`
	MobilePhone  string                `json:"mobile_phone"`
	Password     string                `json:"password" binding:"required,min=8"`
	DeviceId     string                `json:"device_id"`
	SessionId    string                `json:"session_id"`
	RoleIdType   enums.RoleIdType      `json:"role_id_type"`
	GroupRoleId  string                `json:"group_role_id"`
	RoleId       string                `json:"role_id"`
	SuperiorId   string                `json:"superior_id"`
	ReferralCode string                `gorm:"unique" json:"referral_code"`
	JobPosition  enums.JobPositionType `json:"job_position"`

	GroupRole grouprole.GroupRole `gorm:"foreignKey:GroupRoleId"`
	Role      role.Role           `gorm:"foreignKey:RoleId"`
	Superior  *BackofficeUsers    `gorm:"foreignKey:SuperiorId"`

	base.BaseModel
}

func (bou *BackofficeUsers) BeforeCreate(db *gorm.DB) (err error) {
	if bou.SuperiorId == "" {
		bou.SuperiorId = "SYSTEM"
	}

	return
}

type BackofficeUserLoginSessionResponse struct {
	ID           string `json:"userid"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	Expiration   int64  `json:"expiration"`
	RoleId       string `json:"role_id"`
	ReferralCode string `json:"referral_code"`
}

type BackofficeUserApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    string `json:"time"`
}

type BackofficeUserRepository interface {
	Create(backofficeUser *BackofficeUsers) error
	GetUserByEmail(email string) (*BackofficeUsers, error)
	GetActiveUsers() (*[]BackofficeUsers, error)
	Update(user *BackofficeUsers) error
	GetActiveSubordinate(userId string) (*[]BackofficeUsers, error)
	GetActiveUsersByRoleId(roleId string) ([]BackofficeUsers, error)
}
