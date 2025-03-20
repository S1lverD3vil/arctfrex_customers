package user

import (
	"arctfrex-customers/internal/base"

	"gorm.io/gorm"
)

type backofficeUserRepository struct {
	db *gorm.DB
}

func NewBackofficeUserRepository(db *gorm.DB) BackofficeUserRepository {
	return &backofficeUserRepository{db: db}
}

// Create inserts a new user into the database
func (bur *backofficeUserRepository) Create(backofficeUser *BackofficeUsers) error {
	return bur.db.Create(backofficeUser).Error
}

func (bur *backofficeUserRepository) GetUserByEmail(email string) (*BackofficeUsers, error) {
	var backofficeUser BackofficeUsers
	if err := bur.db.Where(&BackofficeUsers{Email: email}).First(&backofficeUser).Error; err != nil {
		return nil, err
	}

	return &backofficeUser, nil
}

func (bur *backofficeUserRepository) GetActiveUsers() (*[]BackofficeUsers, error) {
	var backofficeUsers []BackofficeUsers
	queryParams := BackofficeUsers{
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := bur.db.Find(&backofficeUsers, &queryParams).Error
	if err != nil || backofficeUsers == nil {
		return nil, err
	}

	return &backofficeUsers, nil
}

func (bur *backofficeUserRepository) Update(user *BackofficeUsers) error {
	return bur.db.Updates(user).Error
}

func (bur *backofficeUserRepository) GetActiveUsersByRoleId(roleId string) ([]BackofficeUsers, error) {
	var backofficeUsers []BackofficeUsers

	queryParams := BackofficeUsers{
		RoleId: roleId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := bur.db.Find(&backofficeUsers, &queryParams).Error
	if err != nil || backofficeUsers == nil {
		return nil, err
	}

	return backofficeUsers, nil
}

func (bur *backofficeUserRepository) GetActiveSubordinate(userId string) (*[]BackofficeUsers, error) {
	var backofficeUsers []BackofficeUsers

	queryParams := BackofficeUsers{
		SuperiorId: userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := bur.db.Find(&backofficeUsers, &queryParams).Error
	if err != nil || backofficeUsers == nil {
		return nil, err
	}

	return &backofficeUsers, nil
}
