package role

import (
	"arctfrex-customers/internal/base"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (rr *roleRepository) GetActiveRoles() ([]Role, error) {
	var roles []Role

	queryParams := Role{
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := rr.db.Find(&roles, &queryParams).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (rr *roleRepository) Create(role *Role) error {
	role.IsActive = true

	if rr.db.Where("id = ?", role.ID).Updates(role).RowsAffected == 0 {
		return rr.db.Create(&role).Error
	}

	return nil
}
