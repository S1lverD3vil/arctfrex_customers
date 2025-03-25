package role

import (
	"arctfrex-customers/internal/base"

	"gorm.io/gorm"
)

type Role struct {
	ID             string  `json:"id" gorm:"primaryKey"`
	Name           string  `json:"name"`
	CommissionRate float64 `json:"commission_rate" gorm:"type:double"`
	ParentRoleID   *string `json:"parent_role_id"`

	ParentRole *Role `gorm:"foreignKey:ParentRoleID"`

	base.BaseModel
}

type CreateUserDTO struct {
	ID             string  `json:"id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	CommissionRate float64 `json:"commission_rate"`
	ParentRoleID   *string `json:"parent_role_id"`
}

func CreateRoleSeed(db *gorm.DB) {
	roles := []Role{
		{ID: "HM", Name: "Head of Marketing", CommissionRate: 0.01, BaseModel: base.BaseModel{
			IsActive: true,
		}},
		{ID: "SBM", Name: "Senior Business Manager", CommissionRate: 0.02, ParentRoleID: strPtr("HM"), BaseModel: base.BaseModel{
			IsActive: true,
		}},
		{ID: "BM", Name: "Business Manager", CommissionRate: 0.03, ParentRoleID: strPtr("SBM"), BaseModel: base.BaseModel{
			IsActive: true,
		}},
		{ID: "ABM", Name: "Assistant Business Manager", CommissionRate: 0.05, ParentRoleID: strPtr("BM"), BaseModel: base.BaseModel{
			IsActive: true,
		}},
		{ID: "MKT", Name: "Marketing", CommissionRate: 0.07, ParentRoleID: strPtr("ABM"), BaseModel: base.BaseModel{
			IsActive: true,
		}},
		{ID: "IB", Name: "Freelance", CommissionRate: 0.10, ParentRoleID: strPtr("MKT"), BaseModel: base.BaseModel{
			IsActive: true,
		}},
	}

	for _, role := range roles {
		db.FirstOrCreate(&role, Role{ID: role.ID})
	}
}

func strPtr(s string) *string {
	return &s
}

type RoleApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    string `json:"time"`
}

type RoleRepository interface {
	GetActiveRoles() ([]Role, error)
	Create(role *Role) error
	Update(role *Role) error
	Delete(roleID string) error
	GetByID(roleID string) (*Role, error)
}
