package withdrawal

import (
	"arctfrex-customers/internal/account"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/market"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type WithdrawalUsecase interface {
	Submit(withdrawal *Withdrawal) (string, error)
	Pending(userId, accountId string) error
	Detail(userId, withdrawalId string) (*Withdrawal, error)
	WithdrawalByAccountId(userId, accountId string) (*[]Withdrawals, error)
	BackOfficePending() (*[]BackOfficePendingWithdrawal, error)
	BackOfficePendingDetail(withdrawalId string) (*BackOfficePendingWithdrawalDetail, error)
	BackOfficePendingApproval(backOfficePendingApproval BackOfficePendingApprovalRequest) error
}

type withdrawalUsecase struct {
	withdrawalRepository WithdrawalRepository
	accountRepository    account.AccountRepository
	withdrawalApiclient  WithdrawalApiclient
	marketRepository     market.MarketRepository
}

func NewWithdrawalUsecase(
	dr WithdrawalRepository,
	ar account.AccountRepository,
	wa WithdrawalApiclient,
	mr market.MarketRepository,
) *withdrawalUsecase {
	return &withdrawalUsecase{
		withdrawalRepository: dr,
		accountRepository:    ar,
		withdrawalApiclient:  wa,
		marketRepository:     mr,
	}
}

func (wu *withdrawalUsecase) Submit(withdrawal *Withdrawal) (string, error) {
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

	return withdrawal.ID, wu.withdrawalRepository.Create(withdrawal)
}

func (wu *withdrawalUsecase) Pending(userId, accountId string) error {
	withdrawaldb, _ := wu.withdrawalRepository.GetPendingAccountByAccountIdUserId(accountId, userId)
	if withdrawaldb != nil && withdrawaldb.IsActive {
		return errors.New("withdrawal still in pending approval")
	}

	return nil
}

func (wu *withdrawalUsecase) WithdrawalByAccountId(userId, accountId string) (*[]Withdrawals, error) {
	withdrawals, err := wu.withdrawalRepository.GetWithdrawalsByUserIdAccountId(userId, accountId)
	if err != nil {
		return &[]Withdrawals{}, errors.New("record not found")
	}

	return withdrawals, nil
}

func (wu *withdrawalUsecase) Detail(userId, withdrawalId string) (*Withdrawal, error) {
	withdrawalDetail, err := wu.withdrawalRepository.GetWithdrawalByIdUserId(userId, withdrawalId)
	if withdrawalDetail == nil || err != nil {
		return nil, errors.New("not found")
	}

	return withdrawalDetail, nil
}

func (wu *withdrawalUsecase) BackOfficePending() (*[]BackOfficePendingWithdrawal, error) {
	withdrawals, err := wu.withdrawalRepository.GetBackOfficePendingWithdrawals()
	if err != nil {
		return &[]BackOfficePendingWithdrawal{}, errors.New("record not found")
	}

	return withdrawals, nil
}

func (wu *withdrawalUsecase) BackOfficePendingDetail(withdrawalId string) (*BackOfficePendingWithdrawalDetail, error) {
	withdrawalDetail, err := wu.withdrawalRepository.GetBackOfficePendingWithdrawalDetail(withdrawalId)
	if withdrawalDetail == nil || err != nil {
		return nil, errors.New("user not found")
	}

	return withdrawalDetail, nil
}

func (wu *withdrawalUsecase) BackOfficePendingApproval(backOfficePendingApproval BackOfficePendingApprovalRequest) error {
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

			tradeWithdrawal := TradeWithdrawal{
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
