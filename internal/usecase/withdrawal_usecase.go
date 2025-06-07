package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"arctfrex-customers/internal/api"
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/model"
	"arctfrex-customers/internal/repository"
)

type WithdrawalUsecase interface {
	Submit(withdrawal *model.Withdrawal) (string, error)
	Pending(userId, accountId string) error
	Detail(userId, withdrawalId string) (*model.Withdrawal, error)
	WithdrawalByAccountId(userId, accountId string) (*[]model.Withdrawals, error)
	BackOfficePending() (*[]model.BackOfficePendingWithdrawal, error)
	BackOfficePendingDetail(withdrawalId string) (*model.BackOfficePendingWithdrawalDetail, error)
	BackOfficePendingApproval(backOfficePendingApproval model.BackOfficePendingWithdrawalApprovalRequest) error
	BackOfficePendingSPA(ctx context.Context, request model.WithdrawalBackOfficeParam) (model.BackOfficePendingWithdrawalPagination, error)
	BackOfficePendingMulti(ctx context.Context, request model.WithdrawalBackOfficeParam) (model.BackOfficePendingWithdrawalPagination, error)
}

type withdrawalUsecase struct {
	withdrawalRepository       repository.WithdrawalRepository
	accountRepository          repository.AccountRepository
	withdrawalApiclient        api.WithdrawalApiclient
	marketRepository           repository.MarketRepository
	workflowSettingRepository  repository.WorkflowSettingRepository
	workflowApproverRepository repository.WorkflowApproverRepository
}

func NewWithdrawalUsecase(
	dr repository.WithdrawalRepository,
	ar repository.AccountRepository,
	wa api.WithdrawalApiclient,
	mr repository.MarketRepository,
	workflowSettingRepository repository.WorkflowSettingRepository,
	workflowApproverRepository repository.WorkflowApproverRepository,
) *withdrawalUsecase {
	return &withdrawalUsecase{
		withdrawalRepository:       dr,
		accountRepository:          ar,
		withdrawalApiclient:        wa,
		marketRepository:           mr,
		workflowSettingRepository:  workflowSettingRepository,
		workflowApproverRepository: workflowApproverRepository,
	}
}

func (wu *withdrawalUsecase) Submit(withdrawal *model.Withdrawal) (string, error) {
	withdrawaldb, _ := wu.withdrawalRepository.GetPendingAccountByAccountIdUserId(withdrawal.AccountID, withdrawal.UserID)
	if withdrawaldb != nil && withdrawaldb.IsActive {
		return withdrawaldb.ID, errors.New("withdrawal still in pending approval")
	}

	account, err := wu.accountRepository.GetAccountsByIdUserId(withdrawal.UserID, withdrawal.AccountID)
	if err != nil || account == nil || account.ID == common.STRING_EMPTY {
		return common.STRING_EMPTY, errors.New("record not found")
	}
	withdrawal.AmountUsd, _ = wu.ConvertPriceToUsd(withdrawal.Amount)

	if account.Balance < withdrawal.AmountUsd {
		return common.STRING_EMPTY, errors.New("insufficient balance")
	}

	withdrawalID, err := uuid.NewUUID()
	if err != nil {
		return common.STRING_EMPTY, err
	}

	withdrawal.ID = common.UUIDNormalizer(withdrawalID)
	withdrawal.IsActive = true
	withdrawal.ApprovalStatus = enums.WithdrawalApprovalStatusPending

	err = wu.withdrawalRepository.Create(withdrawal)
	if err != nil {
		return common.STRING_EMPTY, err
	}

	err = wu.AddWorkflowApprover(withdrawal, withdrawal.UserID)
	if err != nil {
		return common.STRING_EMPTY, err
	}

	return withdrawal.ID, err
}

func (wu *withdrawalUsecase) AddWorkflowApprover(withdrawal *model.Withdrawal, userID string) error {
	workflowApprover, err := wu.workflowApproverRepository.GetWorkflowApproverByDocumentId(withdrawal.ID)
	if err != nil {
		return err
	}

	if len(workflowApprover) > 0 && withdrawal.ApprovalStatus == enums.WithdrawalApprovalStatusPending {
		return nil
	}

	workflowSetting, err := wu.workflowSettingRepository.GetWorkflowSettingByWorkflowType(common.WorkflowWithdrawalApprover)
	if err != nil {
		return err
	}

	var config model.WorkflowConfig
	err = json.Unmarshal([]byte(workflowSetting.Config), &config)
	if err != nil {
		return err
	}

	var approvers []model.WorkflowApprover
	for _, approver := range config.Approvers {
		data := model.WorkflowApprover{
			ID:                common.UUIDNormalizer(uuid.New()),
			WorkflowSettingID: &workflowSetting.ID,
			Level:             approver.Level,
			Status:            enums.AccountApprovalStatusPending,
			DocumentID:        withdrawal.ID,
			BaseModel: base.BaseModel{
				CreatedBy: userID,
				IsActive:  true,
			},
		}

		approvers = append(approvers, data)
	}
	err = wu.workflowApproverRepository.CreateBulk(approvers)
	if err != nil {
		return err
	}

	return nil
}

func (wu *withdrawalUsecase) Pending(userId, accountId string) error {
	withdrawaldb, _ := wu.withdrawalRepository.GetPendingAccountByAccountIdUserId(accountId, userId)
	if withdrawaldb != nil && withdrawaldb.IsActive {
		return errors.New("withdrawal still in pending approval")
	}

	return nil
}

func (wu *withdrawalUsecase) WithdrawalByAccountId(userId, accountId string) (*[]model.Withdrawals, error) {
	withdrawals, err := wu.withdrawalRepository.GetWithdrawalsByUserIdAccountId(userId, accountId)
	if err != nil {
		return &[]model.Withdrawals{}, errors.New("record not found")
	}

	return withdrawals, nil
}

func (wu *withdrawalUsecase) Detail(userId, withdrawalId string) (*model.Withdrawal, error) {
	withdrawalDetail, err := wu.withdrawalRepository.GetWithdrawalByIdUserId(userId, withdrawalId)
	if withdrawalDetail == nil || err != nil {
		return nil, errors.New("not found")
	}

	return withdrawalDetail, nil
}

func (wu *withdrawalUsecase) BackOfficePending() (*[]model.BackOfficePendingWithdrawal, error) {
	withdrawals, err := wu.withdrawalRepository.GetBackOfficePendingWithdrawals()
	if err != nil {
		return &[]model.BackOfficePendingWithdrawal{}, errors.New("record not found")
	}

	return withdrawals, nil
}

func (wu *withdrawalUsecase) BackOfficePendingDetail(withdrawalId string) (*model.BackOfficePendingWithdrawalDetail, error) {
	withdrawalDetail, err := wu.withdrawalRepository.GetBackOfficePendingWithdrawalDetail(withdrawalId)
	if withdrawalDetail == nil || err != nil {
		return nil, errors.New("user not found")
	}

	return withdrawalDetail, nil
}

func (wu *withdrawalUsecase) BackOfficePendingApproval(backOfficePendingApproval model.BackOfficePendingWithdrawalApprovalRequest) error {
	pendingWithdrawal, err := wu.withdrawalRepository.GetPendingWithdrawalsById(backOfficePendingApproval.Withdrawalid)
	if err != nil {
		return errors.New("record not found")
	}

	if pendingWithdrawal == nil || pendingWithdrawal.ID == common.STRING_EMPTY {
		return errors.New("record not found")
	}

	switch strings.ToLower(backOfficePendingApproval.Decision) {
	case "approved":
		{
			pendingWithdrawal.ApprovalStatus = enums.WithdrawalApprovalStatusApproved
			pendingWithdrawal.AmountUsd, _ = wu.ConvertPriceToUsd(pendingWithdrawal.Amount)

			account, err := wu.accountRepository.GetAccountsByIdUserId(pendingWithdrawal.UserID, pendingWithdrawal.AccountID)
			if err != nil || account == nil || account.ID == common.STRING_EMPTY {
				return errors.New("record not found")
			}

			tradeWithdrawal := model.TradeWithdrawal{
				Login:  account.MetaLoginId,
				Amount: pendingWithdrawal.AmountUsd * -1,
			}

			demoAccountTopUpData, err := wu.withdrawalApiclient.TradeWithdrawal(tradeWithdrawal)
			if err != nil {
				log.Println(demoAccountTopUpData.Result)
			}

			account.Balance -= pendingWithdrawal.AmountUsd
			account.Equity -= pendingWithdrawal.AmountUsd
			account.ModifiedBy = pendingWithdrawal.UserID
			if err := wu.accountRepository.UpdateAccount(account); err != nil {
				return errors.New("failed to top up account")
			}
		}
	case "rejected":
		{
			pendingWithdrawal.ApprovalStatus = enums.WithdrawalApprovalStatusRejected
		}
	default:
		{
			pendingWithdrawal.ApprovalStatus = enums.WithdrawalApprovalStatusCancelled
		}
	}

	pendingWithdrawal.ApprovedAt = time.Now()
	pendingWithdrawal.ApprovedBy = backOfficePendingApproval.UserLogin

	return wu.withdrawalRepository.UpdateWithdrawalApprovalStatus(pendingWithdrawal)
}

func (du *withdrawalUsecase) ConvertPriceToUsd(amount float64) (float64, error) {
	marketCurrencyRate, err := du.marketRepository.GetMarketCurrencyRate("IDR", "USD")
	if marketCurrencyRate == nil || err != nil {
		return amount, errors.New("market rate not found")
	}
	amountUsd := common.RoundTo4DecimalPlaces(amount * marketCurrencyRate.Rate)
	return amountUsd, nil
}

func (wu *withdrawalUsecase) BackOfficePendingSPA(ctx context.Context, request model.WithdrawalBackOfficeParam) (withdrawal model.BackOfficePendingWithdrawalPagination, err error) {
	withdrawal.Pagination = request.Pagination
	withdrawals, err := wu.withdrawalRepository.GetBackOfficePendingWithdrawalSPA(request)
	if err != nil {
		return withdrawal, err
	}

	withdrawal.Data = withdrawals

	return withdrawal, nil
}

func (wu *withdrawalUsecase) BackOfficePendingMulti(ctx context.Context, request model.WithdrawalBackOfficeParam) (withdrawal model.BackOfficePendingWithdrawalPagination, err error) {
	withdrawal.Pagination = request.Pagination
	withdrawals, err := wu.withdrawalRepository.GetBackOfficePendingWithdrawalMulti(request)
	if err != nil {
		return withdrawal, err
	}

	withdrawal.Data = withdrawals

	return withdrawal, nil
}
