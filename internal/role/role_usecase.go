package role

type RoleUseCase interface {
	All() ([]Role, error)
	Create(role *Role) error
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
	return ru.roleRepository.Create(role)
}
