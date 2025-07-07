package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"arctfrex-customers/internal/api"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/dbtx"
	"arctfrex-customers/internal/dto"
	"arctfrex-customers/internal/model"
	"arctfrex-customers/internal/repository"
)

type WorkflowApproverUsecase interface {
	ApproveRejectWorkflow(approverReject dto.ApproveRejectRequest) (dto.ApproveRejectResponse, error)
}

type workflowApproverUsecase struct {
	workflowApproverRepository repository.WorkflowApproverRepository
	depositRepository          repository.DepositRepository
	withdrawalRepository       repository.WithdrawalRepository
	accountRepository          repository.AccountRepository
	marketRepository           repository.MarketRepository
	depositApiClient           api.DepositApiclient
	withdrawalApiClient        api.WithdrawalApiclient
	db                         *gorm.DB
}

func NewWorkflowApproverUsecase(
	workflowApproverRepository repository.WorkflowApproverRepository,
	depositRepository repository.DepositRepository,
	withdrawalRepository repository.WithdrawalRepository,
	accountRepository repository.AccountRepository,
	marketRepository repository.MarketRepository,
	depositApiClient api.DepositApiclient,
	withdrawalApiClient api.WithdrawalApiclient,
	db *gorm.DB,

) *workflowApproverUsecase {
	return &workflowApproverUsecase{
		workflowApproverRepository: workflowApproverRepository,
		depositRepository:          depositRepository,
		withdrawalRepository:       withdrawalRepository,
		depositApiClient:           depositApiClient,
		withdrawalApiClient:        withdrawalApiClient,
		accountRepository:          accountRepository,
		marketRepository:           marketRepository,
		db:                         db,
	}
}

func (workflowApproverUC workflowApproverUsecase) IsInitialMargin(userID string, accountID string) bool {
	var isInitialMargin bool

	initialMargin, err := workflowApproverUC.depositRepository.GetDepositsByUserIDAccountIDForIsInitialMargin(userID, accountID)
	if err != nil {
		log.Println("Error checking initial margin:", err)
		return isInitialMargin
	}

	if initialMargin != nil && initialMargin.AccountType == enums.AccountTypeReal && initialMargin.Total == 0 {
		isInitialMargin = true
	}

	return isInitialMargin
}

func (workflowApproverUC workflowApproverUsecase) ApproveRejectWorkflow(approverReject dto.ApproveRejectRequest) (response dto.ApproveRejectResponse, err error) {
	var (
		deposit    *model.Deposit
		withdrawal *model.Withdrawal
	)

	switch approverReject.WorkflowType {
	case common.WorkflowDepositApprover:
		deposit, err = workflowApproverUC.depositRepository.GetPendingDepositsById(approverReject.DocumentID)
		if err != nil {
			return response, err
		}

		if deposit.ID == "" {
			return response, fmt.Errorf("deposit not found")
		}

		if deposit.ApprovalStatus != enums.DepositApprovalStatusPending {
			return response, fmt.Errorf("deposit is not pending")
		}

		isInitialMargin := workflowApproverUC.IsInitialMargin(deposit.UserID, deposit.AccountID)
		if isInitialMargin {
			approverReject.DepositType = enums.DepositTypeInitialMargin
		}

		isApprovalComplete, workflowApproverUpdate, err := workflowApproverUC.populateWorkflowApprover(approverReject)
		if err != nil {
			return response, err
		}

		response, err = workflowApproverUC.updateAccountDeposit(workflowApproverUpdate, isApprovalComplete, approverReject, deposit)
		if err != nil {
			return response, err
		}
	case common.WorkflowWithdrawalApprover:
		withdrawal, err = workflowApproverUC.withdrawalRepository.GetPendingWithdrawalsById(approverReject.DocumentID)
		if err != nil {
			return response, err
		}

		if withdrawal.ID == "" {
			return response, fmt.Errorf("withdrawal not found")
		}

		if withdrawal.ApprovalStatus != enums.WithdrawalApprovalStatusPending {
			return response, fmt.Errorf("withdrawal is not pending")
		}

		isApprovalComplete, workflowApproverUpdate, err := workflowApproverUC.populateWorkflowApprover(approverReject)
		if err != nil {
			return response, err
		}

		response, err = workflowApproverUC.updateAccountWithdrawal(workflowApproverUpdate, isApprovalComplete, approverReject, withdrawal)
		if err != nil {
			return response, err
		}
	default:
		return response, fmt.Errorf("invalid approver type")
	}

	response.DocumentID = approverReject.DocumentID
	response.ApproveStatus = approverReject.Status.String()

	return response, err
}

func (workflowApproverUC workflowApproverUsecase) updateAccountDeposit(workflowApproverUpdate model.WorkflowApprover, isApprovalComplete bool, approverReject dto.ApproveRejectRequest, deposit *model.Deposit) (response dto.ApproveRejectResponse, err error) {
	response, err = dbtx.WithTransaction(context.Background(), workflowApproverUC.db, func(tx *dbtx.DepositWithdrawlRepositoryGroup) (dto.ApproveRejectResponse, error) {
		err = tx.WorkflowApproverRepository.UpdateApproverStatus(&workflowApproverUpdate)
		if err != nil {
			return response, err
		}

		isTradeCreditOut := workflowApproverUpdate.Level == 1 && deposit.CreditType == enums.TypeCreditIn
		isTradeCreditIn := workflowApproverUpdate.Level == 2 && deposit.CreditType == enums.TypeCreditIn
		isRejected := enums.DepositApprovalStatus(approverReject.Status) == enums.DepositApprovalStatusRejected

		if isApprovalComplete || enums.DepositApprovalStatus(approverReject.Status) == enums.DepositApprovalStatusRejected || isTradeCreditIn {
			deposit.ApprovedAt = time.Now()

			isDepositEnabled := (isApprovalComplete || isTradeCreditIn) && !isTradeCreditOut && !isRejected
			if isDepositEnabled {
				deposit.AmountUsd, _ = workflowApproverUC.ConvertPriceToUsd(tx, deposit.Amount)
				account, err := tx.AccountRepository.GetAccountsByIdUserId(deposit.UserID, deposit.AccountID)
				if err != nil || account == nil || account.ID == common.STRING_EMPTY {
					return response, errors.New("record not found")
				}

				tradeDepositRequest := model.TradeDeposit{
					Login:  account.MetaLoginId,
					Amount: deposit.AmountUsd,
				}

				demoAccountTopUpData, err := workflowApproverUC.depositApiClient.TradeDeposit(tradeDepositRequest)
				if err != nil {
					log.Println(demoAccountTopUpData.Result, "TradeDeposit failed with error:", err)
					return response, errors.New("failed to top up account due to trade deposit error")
				}

				account.Balance += deposit.AmountUsd
				account.Equity += deposit.AmountUsd
				account.ModifiedBy = deposit.UserID
				if err := tx.AccountRepository.UpdateAccount(account); err != nil {
					return response, errors.New("failed to top up account")
				}
			}
		}

		if enums.DepositApprovalStatus(approverReject.Status) == enums.DepositApprovalStatusRejected {
			deposit.ApprovalStatus = enums.DepositApprovalStatusRejected
		}

		isApproverWithoutCreditInOut := isApprovalComplete && !isTradeCreditIn && !isTradeCreditOut
		if isApproverWithoutCreditInOut {
			deposit.ApprovalStatus = enums.DepositApprovalStatusApproved
		}

		isApproverWithCreditInOut := (isTradeCreditOut && isApprovalComplete) || (isTradeCreditIn && isApprovalComplete)
		if isApproverWithCreditInOut {
			deposit.ApprovalStatus = enums.DepositApprovalStatusApproved
			deposit.CreditType = enums.TypeCreditOut
		}

		deposit.DepositType = approverReject.DepositType
		err = tx.DepositRepository.UpdateDepositApprovalStatus(deposit)
		if err != nil {
			return response, err
		}

		return response, nil
	})

	return response, err
}

func (workflowApproverUC workflowApproverUsecase) updateAccountWithdrawal(workflowApproverUpdate model.WorkflowApprover, isApprovalComplete bool, approverReject dto.ApproveRejectRequest, withdrawal *model.Withdrawal) (response dto.ApproveRejectResponse, err error) {
	response, err = dbtx.WithTransaction(context.Background(), workflowApproverUC.db, func(tx *dbtx.DepositWithdrawlRepositoryGroup) (dto.ApproveRejectResponse, error) {
		err = tx.WorkflowApproverRepository.UpdateApproverStatus(&workflowApproverUpdate)
		if err != nil {
			return response, err
		}

		if isApprovalComplete || enums.WithdrawalApprovalStatus(approverReject.Status) == enums.WithdrawalApprovalStatusRejected {
			withdrawal.ApprovedAt = time.Now()
			withdrawal.ApprovalStatus = enums.WithdrawalApprovalStatusRejected
			if isApprovalComplete {
				withdrawal.ApprovalStatus = enums.WithdrawalApprovalStatusApproved
				withdrawal.AmountUsd, _ = workflowApproverUC.ConvertPriceToUsd(tx, withdrawal.Amount)

				account, err := workflowApproverUC.accountRepository.GetAccountsByIdUserId(withdrawal.UserID, withdrawal.AccountID)
				if err != nil || account == nil || account.ID == common.STRING_EMPTY {
					return response, errors.New("record not found")
				}

				tradeWithdrawal := model.TradeWithdrawal{
					Login:  account.MetaLoginId,
					Amount: withdrawal.AmountUsd * -1,
				}

				demoAccountTopUpData, err := workflowApproverUC.withdrawalApiClient.TradeWithdrawal(tradeWithdrawal)
				if err != nil {
					log.Println(demoAccountTopUpData.Result)
				}

				account.Balance -= withdrawal.AmountUsd
				account.Equity -= withdrawal.AmountUsd
				account.ModifiedBy = withdrawal.UserID
				if err := tx.AccountRepository.UpdateAccount(account); err != nil {
					return response, errors.New("failed to top up account")
				}
			}

			err = tx.WithdrawalRepository.UpdateWithdrawalApprovalStatus(withdrawal)
			if err != nil {
				return response, err
			}
		}

		return response, nil
	})

	return response, err
}

func (workflowApproverUsecase workflowApproverUsecase) populateWorkflowApprover(approverReject dto.ApproveRejectRequest) (isApproverComplete bool, workflowApproverUpdate model.WorkflowApprover, err error) {
	approveCount := 0
	now := time.Now()

	workflowApprovers, err := workflowApproverUsecase.workflowApproverRepository.GetWorkflowApproverByDocumentId(approverReject.DocumentID)
	if err != nil || len(workflowApprovers) == 0 {
		return isApproverComplete, workflowApproverUpdate, fmt.Errorf("workflow approver not found")
	}

	for _, workflowApprover := range workflowApprovers {
		isStatusChangeValid := workflowApprover.Level == approverReject.Level && workflowApprover.Status == enums.AccountApprovalStatusPending
		isStatusInvalid := workflowApprover.Level == approverReject.Level && workflowApprover.Status != enums.AccountApprovalStatusPending

		if isStatusInvalid {
			return isApproverComplete, workflowApproverUpdate, fmt.Errorf("workflow approver status is not pending")
		}

		isRejectedDisallowed := workflowApprover.Level == 2 && workflowApprover.Status == enums.AccountApprovalStatusApproved && enums.DepositApprovalStatus(approverReject.Status) == enums.DepositApprovalStatusRejected && approverReject.WorkflowType == common.WorkflowDepositApprover && approverReject.Level == 1
		if isRejectedDisallowed {
			return isApproverComplete, workflowApproverUpdate, fmt.Errorf("cannot reject the settlement level has already been approved")
		}

		if isStatusChangeValid {
			workflowApprover.Status = approverReject.Status
			workflowApprover.ApprovedBy = approverReject.UserID
			workflowApprover.ApprovedAt = now
			workflowApproverUpdate = workflowApprover
		}

		if workflowApprover.Status == enums.AccountApprovalStatusApproved {
			approveCount++
		}
	}

	if approveCount == len(workflowApprovers) {
		isApproverComplete = true
	}

	return isApproverComplete, workflowApproverUpdate, nil
}

func (workflowApproverUC workflowApproverUsecase) ConvertPriceToUsd(tx *dbtx.DepositWithdrawlRepositoryGroup, amount float64) (float64, error) {
	marketCurrencyRate, err := tx.MarketRepository.GetMarketCurrencyRate("IDR", "USD")
	if marketCurrencyRate == nil || err != nil {
		return amount, errors.New("market rate not found")
	}
	amountUsd := common.RoundTo4DecimalPlaces(amount * marketCurrencyRate.Rate)
	return amountUsd, nil
}
