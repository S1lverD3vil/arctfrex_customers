package grouprole

import (
	"arctfrex-customers/internal/base"

	"gorm.io/gorm"
)

type groupRoleRepository struct {
	db *gorm.DB
}

func NewGroupRoleRepository(db *gorm.DB) GroupRoleRepository {
	return &groupRoleRepository{db: db}
}

func (rr *groupRoleRepository) GetActiveGroupRoles() ([]GroupRole, error) {
	var roles []GroupRole

	queryParams := GroupRole{
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

func (rr *groupRoleRepository) GetByID(groupRoleID string) (*GroupRole, error) {
	queryParams := GroupRole{
		ID: groupRoleID,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	var groupRole GroupRole
	err := rr.db.Where(queryParams).First(&groupRole).Error
	if err != nil {
		return nil, err
	}
	return &groupRole, nil
}

func (rr *groupRoleRepository) Create(groupRole *GroupRole) error {
	groupRole.IsActive = true

	if rr.db.Where("id = ?", groupRole.ID).Updates(groupRole).RowsAffected == 0 {
		return rr.db.Create(&groupRole).Error
	}

	return nil
}

func (rr *groupRoleRepository) Update(groupRole *GroupRole) error {
	return rr.db.Where("id = ?", groupRole.ID).Updates(groupRole).Error

}

func (rr *groupRoleRepository) Delete(groupRoleID string) error {
	groupRole, err := rr.GetByID(groupRoleID)
	if err != nil {
		return err
	}

	return rr.db.Model(&GroupRole{}).Where("id = ?", groupRole.ID).Update("is_active", false).Error

}
