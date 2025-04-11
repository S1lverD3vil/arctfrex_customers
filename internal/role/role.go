package role

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/grouprole"
	"strings"

	"gorm.io/gorm"
)

type Role struct {
	ID             string  `json:"id" gorm:"primaryKey"`
	Name           string  `json:"name"`
	CommissionRate float64 `json:"commission_rate" gorm:"type:double precision"`
	ParentRoleID   *string `json:"parent_role_id"`
	GroupRoleID    string  `json:"role_group_id"`

	ParentRole *Role               `gorm:"foreignKey:ParentRoleID"`
	GroupRole  grouprole.GroupRole `gorm:"foreignKey:GroupRoleID"`

	base.BaseModel
}

type CreateUserDTO struct {
	ID             string  `json:"id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	CommissionRate float64 `json:"commission_rate"`
	ParentRoleID   *string `json:"parent_role_id"`
}

type RoleApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    string `json:"time"`
}

// BeforeCreate memastikan ID diformat sebelum disimpan
func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = formatRoleID(r.ID)
	return
}

// BeforeUpdate memastikan ID tidak bisa diubah setelah dibuat
func (r *Role) BeforeUpdate(tx *gorm.DB) (err error) {
	var existing Role
	if err := tx.First(&existing, "id = ?", r.ID).Error; err == nil {
		// Jika ID sudah ada, gunakan ID lama
		r.ID = existing.ID
	}
	return
}

// formatRoleID mengubah string menjadi lowercase dan mengganti spasi dengan underscore
func formatRoleID(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

type RoleRepository interface {
	GetActiveRoles() ([]Role, error)
	Create(role *Role) error
	Update(role *Role) error
	Delete(roleID string) error
	GetByID(roleID string) (*Role, error)
}
