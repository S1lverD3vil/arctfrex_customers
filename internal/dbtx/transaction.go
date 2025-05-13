package dbtx

import (
	"context"

	"gorm.io/gorm"
)

func WithTransaction[T any](ctx context.Context, db *gorm.DB, txFunc func(*DepositWithdrawlRepositoryGroup) (T, error)) (T, error) {
	var result T

	// Start a new transaction
	tx := db.Begin()

	// If transaction begins, run the provided function
	if tx.Error != nil {
		return result, tx.Error
	}

	// Create a repository group and execute the logic within the transaction
	repoGroup := NewDepositWithdrawlRepositoryGroup(tx)

	// Call the provided transaction function
	result, err := txFunc(repoGroup)
	if err != nil {
		// Rollback transaction if error occurs
		tx.Rollback()
		return result, err
	}

	// Commit if everything goes well
	if err := tx.Commit().Error; err != nil {
		return result, err
	}

	return result, nil
}
