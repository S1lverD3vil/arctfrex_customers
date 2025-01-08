package account

import (
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AccountUsecase interface {
	Submit(account *Account) error
	Pending(userId string) (*[]Account, error)
	PendingCheck(userId string) (*[]Account, error)
	GetAccounts(userId string) (*[]Account, error)
	TopUpAccount(topUpAccount TopUpAccount) error
	BackOfficeAll() (*[]BackOfficeAllAccount, error)
	BackOfficePending() (*[]BackOfficePendingAccount, error)
	BackOfficePendingApproval(backOfficePendingApproval BackOfficePendingApprovalRequest) error
}

type accountUsecase struct {
	accountRepository AccountRepository
	accountApiclient  AccountApiclient
}

func NewAccountUsecase(
	ar AccountRepository,
	aa AccountApiclient,
) *accountUsecase {
	return &accountUsecase{
		accountRepository: ar,
		accountApiclient:  aa,
	}
}

func (au *accountUsecase) Submit(account *Account) error {
	accountdb, _ := au.accountRepository.GetPendingAccountByUserId(account.UserID)
	if accountdb != nil && accountdb.IsActive {
		return errors.New("account still in pending approval")
	}

	accountID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	account.ID = common.UUIDNormalizer(accountID)
	account.IsActive = true
	account.Type = enums.AccountTypeReal
	account.ApprovalStatus = enums.AccountApprovalStatusPending

	return au.accountRepository.Create(account)
}

func (au *accountUsecase) Pending(userId string) (*[]Account, error) {
	accounts, err := au.accountRepository.GetPendingAccountsByUserdId(userId)
	if err != nil {
		return &[]Account{}, errors.New("record not found")
	}

	return accounts, nil
}

func (au *accountUsecase) PendingCheck(userId string) (*[]Account, error) {
	accounts, err := au.accountRepository.GetPendingAccountsByUserdId(userId)
	if accounts != nil {
		return accounts, errors.New("account still in pending approval")
	}

	if err != nil {
		return nil, errors.New("record not found")
	}

	return accounts, nil
}

func (au *accountUsecase) GetAccounts(userId string) (*[]Account, error) {
	accounts, err := au.accountRepository.GetAccountsByUserdId(userId)
	if err != nil {
		return &[]Account{}, errors.New("record not found")
	}

	// Sort the slice based on the Accounts field descending
	sort.Slice(*accounts, func(i, j int) bool {
		return (*accounts)[i].Type > (*accounts)[j].Type
	})

	return accounts, nil
}

func (au *accountUsecase) TopUpAccount(topUpAccount TopUpAccount) error {
	account, err := au.accountRepository.GetAccountsByIdUserId(topUpAccount.UserID, topUpAccount.ID)
	if err != nil || account == nil || account.ID == common.STRING_EMPTY {
		return errors.New("record not found")
	}

	demoAccountTopUp := DemoAccountTopUp{
		Login:  account.MetaLoginId,
		Amount: 1000,
	}
	demoAccountTopUpData, err := au.accountApiclient.DemoAccountTopUp(demoAccountTopUp)
	if err != nil {
		log.Println(demoAccountTopUpData.Result)
	}

	// account.Balance += topUpAccount.Amount
	// account.Equity += topUpAccount.Amount
	account.Balance += demoAccountTopUp.Amount
	account.Equity += demoAccountTopUp.Amount
	account.ModifiedBy = topUpAccount.UserID

	return au.accountRepository.UpdateAccount(account)
}

func (au *accountUsecase) BackOfficeAll() (*[]BackOfficeAllAccount, error) {
	accounts, err := au.accountRepository.GetBackOfficeAllAccounts()

	fmt.Printf("Pending Accounts: %+v\n", accounts)

	if err != nil {
		return &[]BackOfficeAllAccount{}, errors.New("record not found")
	}

	return accounts, nil
}

func (au *accountUsecase) BackOfficePending() (*[]BackOfficePendingAccount, error) {
	accounts, err := au.accountRepository.GetBackOfficePendingAccounts()

	fmt.Printf("Pending Accounts: %+v\n", accounts)

	if err != nil {
		return &[]BackOfficePendingAccount{}, errors.New("record not found")
	}

	return accounts, nil
}

func (au *accountUsecase) BackOfficePendingApproval(backOfficePendingApproval BackOfficePendingApprovalRequest) error {
	fmt.Printf("Pending Accoun requestt: %+v\n", backOfficePendingApproval)

	pendingAccount, err := au.accountRepository.GetPendingAccountsById(backOfficePendingApproval.Accountid)
	if err != nil {
		return errors.New("record not found")
	}

	if pendingAccount == nil || pendingAccount.ID == common.STRING_EMPTY {
		return errors.New("record not found")
	}

	fmt.Printf("Pending Account: %+v\n", pendingAccount)

	switch strings.ToLower(backOfficePendingApproval.Decision) {
	case "approved":
		{
			pendingAccount.ApprovalStatus = enums.AccountApprovalStatusApproved

			accountUserData, err := au.accountRepository.GetBackOfficePendingAccountUserData(pendingAccount.UserID)
			if err != nil {
				return err
			}

			securedPassword, err := common.GenerateSecurePassword()
			if err != nil {
				return err
			}
			clientAdd := ClientAdd{
				Name:     accountUserData.Name,
				Password: securedPassword,
				Group:    "demo\\PKB\\B-USD-SFL-MAR-C5-SWAP",
				Leverage: 100,
				Email:    accountUserData.Email,
			}

			clientAddData, err := au.accountApiclient.ClientAdd(clientAdd)
			if err != nil {
				return err
			}

			//fmt.Printf("Client Add Data: %+v\n", clientAddData)
			pendingAccount.MetaLoginId = clientAddData.Login
			pendingAccount.MetaLoginPassword = clientAdd.Password
		}
	case "rejected":
		{
			pendingAccount.ApprovalStatus = enums.AccountApprovalStatusRejected
		}
	default:
		{
			pendingAccount.ApprovalStatus = enums.AccountApprovalStatusCancelled
		}
	}

	pendingAccount.ApprovedAt = time.Now()
	pendingAccount.ApprovedBy = backOfficePendingApproval.UserLogin

	return au.accountRepository.UpdateAccountApprovalStatus(pendingAccount)
}
