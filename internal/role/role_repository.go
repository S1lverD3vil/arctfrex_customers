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

func (rr *roleRepository) GetByID(roleID string) (*Role, error) {
	queryParams := Role{
		ID: roleID,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	var role Role
	err := rr.db.Where(queryParams).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (rr *roleRepository) Create(role *Role) error {
	role.IsActive = true

	if rr.db.Where("id = ?", role.ID).Updates(role).RowsAffected == 0 {
		return rr.db.Create(&role).Error
	}

	return nil
}

func (rr *roleRepository) Update(role *Role) error {
	return rr.db.Where("id = ?", role.ID).Updates(role).Error

}

func (rr *roleRepository) Delete(roleID string) error {
	role, err := rr.GetByID(roleID)
	if err != nil {
		return err
	}

	return rr.db.Model(&Role{}).Where("id = ?", role.ID).Update("is_active", false).Error

}
