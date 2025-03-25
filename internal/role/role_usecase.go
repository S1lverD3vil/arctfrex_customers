package role

import (
	"errors"
)

type RoleUseCase interface {
	All() ([]Role, error)
	Create(role *Role) error
	Update(role *Role) error
	Delete(roleID string) error
	GetByID(roleID string) (*Role, error)
}

type roleUseCase struct {
	roleRepository RoleRepository
}

func NewRoleUseCase(rr RoleRepository) *roleUseCase {
	return &roleUseCase{
		roleRepository: rr,
	}
}

func (ru *roleUseCase) All() ([]Role, error) {
	roles, err := ru.roleRepository.GetActiveRoles()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (ru *roleUseCase) Create(role *Role) error {
	existingRole, _ := ru.roleRepository.GetByID(role.ID)

	if existingRole != nil {
		return errors.New("Role already exists")
	}

	return ru.roleRepository.Create(role)
}

func (ru *roleUseCase) Update(role *Role) error {
	return ru.roleRepository.Update(role)
}

func (ru *roleUseCase) Delete(roleID string) error {
	return ru.roleRepository.Delete(roleID)
}

func (ru *roleUseCase) GetByID(roleID string) (*Role, error) {
	return ru.roleRepository.GetByID(roleID)
}
