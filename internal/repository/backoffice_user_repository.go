package repository

import (
	"gorm.io/gorm"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/model"
)

type BackofficeUserRepository interface {
	Create(backofficeUser *model.BackofficeUsers) error
	GetUserByEmail(email string) (*model.BackofficeUsers, error)
	GetActiveUsers() (*[]model.BackofficeUsers, error)
	Update(user *model.BackofficeUsers) error
	GetActiveSubordinate(userId string) (*[]model.BackofficeUsers, error)
	GetActiveUsersByRoleId(roleId string) ([]model.BackofficeUsers, error)
	GetUserByUserId(userID string) (*model.BackofficeUsers, error)
}

type backofficeUserRepository struct {
	db *gorm.DB
}

func NewBackofficeUserRepository(db *gorm.DB) BackofficeUserRepository {
	return &backofficeUserRepository{db: db}
}

// Create inserts a new user into the database
func (bur *backofficeUserRepository) Create(backofficeUser *model.BackofficeUsers) error {
	return bur.db.Create(backofficeUser).Error
}

func (bur *backofficeUserRepository) GetUserByEmail(email string) (*model.BackofficeUsers, error) {
	var backofficeUser model.BackofficeUsers
	if err := bur.db.Where(&model.BackofficeUsers{Email: email}).First(&backofficeUser).Error; err != nil {
		return nil, err
	}

	return &backofficeUser, nil
}

func (bur *backofficeUserRepository) GetActiveUsers() (*[]model.BackofficeUsers, error) {
	var backofficeUsers []model.BackofficeUsers
	queryParams := model.BackofficeUsers{
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

func (bur *backofficeUserRepository) Update(user *model.BackofficeUsers) error {
	return bur.db.Updates(user).Error
}

func (bur *backofficeUserRepository) GetActiveUsersByRoleId(roleId string) ([]model.BackofficeUsers, error) {
	var backofficeUsers []model.BackofficeUsers

	queryParams := model.BackofficeUsers{
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

func (bur *backofficeUserRepository) GetActiveSubordinate(userId string) (*[]model.BackofficeUsers, error) {
	var backofficeUsers []model.BackofficeUsers

	queryParams := model.BackofficeUsers{
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

func (bur *backofficeUserRepository) GetUserByUserId(userID string) (*model.BackofficeUsers, error) {
	var backofficeUser model.BackofficeUsers
	if err := bur.db.Where(&model.BackofficeUsers{ID: userID}).First(&backofficeUser).Error; err != nil {
		return nil, err
	}

	return &backofficeUser, nil
}
