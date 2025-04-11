package grouprole

import (
	"arctfrex-customers/internal/base"
	"strings"

	"gorm.io/gorm"
)

type GroupRole struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`

	base.BaseModel
}

type GroupRoleApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    string `json:"time"`
}

type CreateUserDTO struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// BeforeCreate memastikan ID diformat sebelum disimpan
func (r *GroupRole) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = formatRoleID(r.ID)
	return
}

// BeforeUpdate memastikan ID tidak bisa diubah setelah dibuat
func (r *GroupRole) BeforeUpdate(tx *gorm.DB) (err error) {
	var existing GroupRole
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

type GroupRoleRepository interface {
	GetActiveGroupRoles() ([]GroupRole, error)
	Create(role *GroupRole) error
	Update(role *GroupRole) error
	Delete(roleID string) error
	GetByID(roleID string) (*GroupRole, error)
}
