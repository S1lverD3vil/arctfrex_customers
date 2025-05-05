package deposit

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"arctfrex-customers/internal/account"
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/market"
	"arctfrex-customers/internal/workflowapprover"
	"arctfrex-customers/internal/workflowsetting"
)

type DepositUsecase interface {
	Submit(deposit *Deposit) (string, error)
	Pending(userId, accountId string) error
	DepositByAccountId(userId, accountId string) (*[]Deposits, error)
	Detail(userId, depositId string) (*Deposit, error)
	BackOfficePending() (*[]BackOfficePendingDeposit, error)
	BackOfficePendingDetail(depositId string) (*BackOfficePendingDepositDetail, error)
	BackOfficePendingApproval(backOfficePendingApproval BackOfficePendingApprovalRequest) error
	BackOfficePendingSPA(ctx context.Context, request DepositBackOfficeParam) (BackOfficePendingDepositPagination, error)
	BackOfficePendingMulti(ctx context.Context, request DepositBackOfficeParam) (BackOfficePendingDepositPagination, error)
	BackOfficeCreditSPA(ctx context.Context, request CreditBackOfficeParam) (BackOfficeCreditPagination, error)
	BackOfficeCreditMulti(ctx context.Context, request CreditBackOfficeParam) (BackOfficeCreditPagination, error)
}

type depositUsecase struct {
	depositRepository          DepositRepository
	accountRepository          account.AccountRepository
	depositApiclient           DepositApiclient
	marketRepository           market.MarketRepository
	workflowSettingRepository  workflowsetting.WorkflowSettingRepository
	workflowApproverRepository workflowapprover.WorkflowApproverRepository
}

func NewDepositUsecase(
	dr DepositRepository,
	ar account.AccountRepository,
	da DepositApiclient,
	mr market.MarketRepository,
	workflowSettingRepository workflowsetting.WorkflowSettingRepository,
	workflowApproverRepository workflowapprover.WorkflowApproverRepository,
) *depositUsecase {
	return &depositUsecase{
		depositRepository:          dr,
		accountRepository:          ar,
		depositApiclient:           da,
		marketRepository:           mr,
		workflowSettingRepository:  workflowSettingRepository,
		workflowApproverRepository: workflowApproverRepository,
	}
}

func (du *depositUsecase) Submit(deposit *Deposit) (string, error) {
	depositPending, _ := du.depositRepository.GetPendingAccountByAccountIdUserId(deposit.AccountID, deposit.UserID)
	if depositPending != nil && depositPending.IsActive {
		return depositPending.ID, errors.New("deposit still in pending approval")
	}

	depositdb, _ := du.depositRepository.GetNewDepositByAccountIdUserId(deposit.AccountID, deposit.UserID)
	if depositdb != nil && depositdb.IsActive {
		depositdb.Amount = deposit.Amount
		depositdb.AmountUsd, _ = du.ConvertPriceToUsd(deposit.Amount)
		// pendingDeposit.AmountUsd, _ = du.ConvertPriceToUsd(pendingDeposit.Amount)
		if depositdb.DepositPhoto != common.STRING_EMPTY {
			depositdb.ApprovalStatus = enums.DepositApprovalStatusPending
		}
		depositdb.BankName = deposit.BankName
		depositdb.BankAccountNumber = deposit.BankAccountNumber
		depositdb.BankBeneficiaryName = deposit.BankBeneficiaryName
		depositdb.ModifiedBy = deposit.UserID

		err := du.depositRepository.Update(depositdb)
		if err != nil {
			return common.STRING_EMPTY, err
		}

		err = du.AddWorkflowApprover(deposit, deposit.UserID)
		if err != nil {
			return common.STRING_EMPTY, err
		}

		return depositdb.ID, nil
	}

	depositID, err := uuid.NewUUID()
	if err != nil {
		return common.STRING_EMPTY, err
	}

	deposit.ID = common.UUIDNormalizer(depositID)
	deposit.IsActive = true
	deposit.ApprovalStatus = enums.DepositApprovalStatusNew
	deposit.AmountUsd, _ = du.ConvertPriceToUsd(deposit.Amount)

	err = du.depositRepository.Create(deposit)
	if err != nil {
		return common.STRING_EMPTY, err
	}

	return deposit.ID, du.depositRepository.Create(deposit)
}

func (du *depositUsecase) AddWorkflowApprover(depositdb *Deposit, userID string) error {
	workflowSetting, err := du.workflowSettingRepository.GetWorkflowSettingByWorkflowType(common.WorkflowDepositApprover)
	if err != nil {
		return err
	}

	var config workflowsetting.WorkflowConfig
	err = json.Unmarshal([]byte(workflowSetting.Config), &config)
	if err != nil {
		return err
	}

	var approvers []workflowapprover.WorkflowApprover
	for _, approver := range config.Approvers {
		data := workflowapprover.WorkflowApprover{
			ID:                common.UUIDNormalizer(uuid.New()),
			WorkflowSettingID: &workflowSetting.ID,
			Level:             approver.Level,
			Status:            enums.AccountApprovalStatusPending,
			DocumentID:        depositdb.ID,
			BaseModel: base.BaseModel{
				CreatedBy: userID,
				IsActive:  true,
			},
		}

		approvers = append(approvers, data)
	}
	err = du.workflowApproverRepository.CreateBulk(approvers)
	if err != nil {
		return err
	}

	return nil
}

func (du *depositUsecase) Pending(userId, accountId string) error {
	depositPending, _ := du.depositRepository.GetPendingAccountByAccountIdUserId(accountId, userId)
	if depositPending != nil && depositPending.IsActive {
		return errors.New("deposit still in pending approval")
	}

	return nil
}

func (du *depositUsecase) DepositByAccountId(userId, accountId string) (*[]Deposits, error) {
	deposits, err := du.depositRepository.GetDepositsByUserIdAccountId(userId, accountId)
	if err != nil {
		return &[]Deposits{}, errors.New("record not found")
	}

	return deposits, nil
}

func (du *depositUsecase) Detail(userId, depositId string) (*Deposit, error) {
	depositDetail, err := du.depositRepository.GetDepositByIdUserId(userId, depositId)
	if depositDetail == nil || err != nil {
		return nil, errors.New("not found")
	}

	return depositDetail, nil
}

func (du *depositUsecase) BackOfficePending() (*[]BackOfficePendingDeposit, error) {
	deposits, err := du.depositRepository.GetBackOfficePendingDeposits()
	if err != nil {
		return &[]BackOfficePendingDeposit{}, errors.New("record not found")
	}

	return deposits, nil
}

func (du *depositUsecase) BackOfficePendingDetail(depositId string) (*BackOfficePendingDepositDetail, error) {
	depositDetail, err := du.depositRepository.GetBackOfficePendingDepositDetail(depositId)
	if depositDetail == nil || err != nil {
		return nil, errors.New("user not found")
	}

	return depositDetail, nil
}

func (du *depositUsecase) BackOfficePendingApproval(backOfficePendingApproval BackOfficePendingApprovalRequest) error {
	pendingDeposit, err := du.depositRepository.GetPendingDepositsById(backOfficePendingApproval.Depositid)
	if err != nil {
		return errors.New("record not found")
	}

	if pendingDeposit == nil || pendingDeposit.ID == common.STRING_EMPTY {
		return errors.New("record not found")
	}

	switch strings.ToLower(backOfficePendingApproval.Decision) {
	case "approved":
		{
			pendingDeposit.ApprovalStatus = enums.DepositApprovalStatusApproved
			pendingDeposit.AmountUsd, _ = du.ConvertPriceToUsd(pendingDeposit.Amount)

			account, err := du.accountRepository.GetAccountsByIdUserId(pendingDeposit.UserID, pendingDeposit.AccountID)
			if err != nil || account == nil || account.ID == common.STRING_EMPTY {
				return errors.New("record not found")
			}

			tradeDeposit := TradeDeposit{
				Login:  account.MetaLoginId,
				Amount: pendingDeposit.AmountUsd,
			}

			demoAccountTopUpData, err := du.depositApiclient.TradeDeposit(tradeDeposit)
			if err != nil {
				log.Println(demoAccountTopUpData.Result)
			}

			account.Balance += pendingDeposit.AmountUsd
			account.Equity += pendingDeposit.AmountUsd
			account.ModifiedBy = pendingDeposit.UserID
			if err := du.accountRepository.UpdateAccount(account); err != nil {
				return errors.New("failed to top up account")
			}
		}
	case "rejected":
		{
			pendingDeposit.ApprovalStatus = enums.DepositApprovalStatusRejected
		}
	default:
		{
			pendingDeposit.ApprovalStatus = enums.DepositApprovalStatusCancelled
		}
	}

	pendingDeposit.DepositType = backOfficePendingApproval.DepositType
	pendingDeposit.ApprovedAt = time.Now()
	pendingDeposit.ApprovedBy = backOfficePendingApproval.UserLogin

	return du.depositRepository.UpdateDepositApprovalStatus(pendingDeposit)
}

func (du *depositUsecase) ConvertPriceToUsd(amount float64) (float64, error) {
	marketCurrencyRate, err := du.marketRepository.GetMarketCurrencyRate("IDR", "USD")
	if marketCurrencyRate == nil || err != nil {
		return amount, errors.New("market rate not found")
	}
	amountUsd := common.RoundTo4DecimalPlaces(amount * marketCurrencyRate.Rate)
	return amountUsd, nil
}

func (du *depositUsecase) BackOfficePendingSPA(ctx context.Context, request DepositBackOfficeParam) (deposits BackOfficePendingDepositPagination, err error) {
	deposits.Pagination = request.Pagination
	deposit, err := du.depositRepository.GetBackOfficePendingDepositSPA(request)
	if err != nil {
		return deposits, err
	}

	deposits.Data = deposit

	return deposits, nil
}

func (du *depositUsecase) BackOfficePendingMulti(ctx context.Context, request DepositBackOfficeParam) (deposits BackOfficePendingDepositPagination, err error) {
	deposits.Pagination = request.Pagination
	deposit, err := du.depositRepository.GetBackOfficePendingDepositMulti(request)
	if err != nil {
		return deposits, err
	}

	deposits.Data = deposit

	return deposits, nil
}

func (du *depositUsecase) BackOfficeCreditSPA(ctx context.Context, request CreditBackOfficeParam) (credit BackOfficeCreditPagination, err error) {
	credit.Pagination = request.Pagination
	credits, err := du.depositRepository.GetBackOfficeCreditSPA(request)
	if err != nil {
		return credit, err
	}

	credit.Data = credits

	return credit, nil
}

func (du *depositUsecase) BackOfficeCreditMulti(ctx context.Context, request CreditBackOfficeParam) (credit BackOfficeCreditPagination, err error) {
	credit.Pagination = request.Pagination
	credits, err := du.depositRepository.GetBackOfficeCreditMulti(request)
	if err != nil {
		return credit, err
	}

	credit.Data = credits

	return credit, nil
}
