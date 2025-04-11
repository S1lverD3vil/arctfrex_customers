package grouprole

import (
	"errors"
)

type GroupRoleUseCase interface {
	All() ([]GroupRole, error)
	Create(role *GroupRole) error
	Update(role *GroupRole) error
	Delete(roleID string) error
	GetByID(roleID string) (*GroupRole, error)
}

type groupRoleUseCase struct {
	groupRoleRespository GroupRoleRepository
}

func NewGroupRoleUseCase(grr GroupRoleRepository) *groupRoleUseCase {
	return &groupRoleUseCase{
		groupRoleRespository: grr,
	}
}

func (ru *groupRoleUseCase) All() ([]GroupRole, error) {
	roles, err := ru.groupRoleRespository.GetActiveGroupRoles()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (ru *groupRoleUseCase) Create(groupRole *GroupRole) error {
	existingRole, _ := ru.groupRoleRespository.GetByID(groupRole.ID)

	if existingRole != nil {
		return errors.New("Group Role already exists")
	}

	return ru.groupRoleRespository.Create(groupRole)
}

func (ru *groupRoleUseCase) Update(groupRole *GroupRole) error {
	return ru.groupRoleRespository.Update(groupRole)
}

func (ru *groupRoleUseCase) Delete(roleID string) error {
	return ru.groupRoleRespository.Delete(roleID)
}

func (ru *groupRoleUseCase) GetByID(roleID string) (*GroupRole, error) {
	return ru.groupRoleRespository.GetByID(roleID)
}
