package dbtx

import (
	"gorm.io/gorm"

	"arctfrex-customers/internal/repository"
)

type DepositWithdrawlRepositoryGroup struct {
	WorkflowApproverRepository repository.WorkflowApproverRepository
	DepositRepository          repository.DepositRepository
	WithdrawalRepository       repository.WithdrawalRepository
	AccountRepository          repository.AccountRepository
	MarketRepository           repository.MarketRepository
}

func NewDepositWithdrawlRepositoryGroup(tx *gorm.DB) *DepositWithdrawlRepositoryGroup {
	return &DepositWithdrawlRepositoryGroup{
		WorkflowApproverRepository: repository.NewWorkflowApproverRepository(tx),
		DepositRepository:          repository.NewDepositRepository(tx),
		WithdrawalRepository:       repository.NewWithdrawalRepository(tx),
		AccountRepository:          repository.NewAccountRepository(tx),
		MarketRepository:           repository.NewMarketRepository(tx),
	}
}
