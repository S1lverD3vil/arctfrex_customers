package usecase

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"arctfrex-customers/internal/api"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/model"
	"arctfrex-customers/internal/repository"
)

type AccountUsecase interface {
	Submit(account *model.Account) error
	Pending(userId string) (*[]model.Account, error)
	PendingCheck(userId string) (*[]model.Account, error)
	GetAccounts(userId string) (*[]model.Account, error)
	TopUpAccount(topUpAccount model.TopUpAccount) error
	BackOfficeAll(request model.BackOfficeAllAccountRequest) (response model.BackOfficeAllAccountResponse, err error)
	BackOfficeAccountByMenuType(request model.BackOfficeAccountByMenuTypeRequest) (response model.BackOfficeAccountByMenuTypeResponse, err error)
	BackOfficePending(request model.BackOfficePendingAccountRequest) (response model.BackOfficePendingAccountResponse, err error)
	BackOfficePendingApproval(backOfficePendingApproval model.BackOfficePendingAccountApprovalRequest) error
}

type accountUsecase struct {
	accountRepository repository.AccountRepository
	accountApiclient  api.AccountApiclient
}

func NewAccountUsecase(
	ar repository.AccountRepository,
	aa api.AccountApiclient,
) *accountUsecase {
	return &accountUsecase{
		accountRepository: ar,
		accountApiclient:  aa,
	}
}

func (au *accountUsecase) Submit(account *model.Account) error {
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

func (au *accountUsecase) Pending(userId string) (*[]model.Account, error) {
	accounts, err := au.accountRepository.GetPendingAccountsByUserdId(userId)
	if err != nil {
		return &[]model.Account{}, errors.New("record not found")
	}

	return accounts, nil
}

func (au *accountUsecase) PendingCheck(userId string) (*[]model.Account, error) {
	accounts, err := au.accountRepository.GetPendingAccountsByUserdId(userId)
	if accounts != nil {
		return accounts, errors.New("account still in pending approval")
	}

	if err != nil {
		return nil, errors.New("record not found")
	}

	return accounts, nil
}

func (au *accountUsecase) GetAccounts(userId string) (*[]model.Account, error) {
	accounts, err := au.accountRepository.GetAccountsByUserdId(userId)
	if err != nil {
		return &[]model.Account{}, errors.New("record not found")
	}

	// Sort the slice based on the Accounts field descending
	sort.Slice(*accounts, func(i, j int) bool {
		return (*accounts)[i].Type > (*accounts)[j].Type
	})

	return accounts, nil
}

func (au *accountUsecase) TopUpAccount(topUpAccount model.TopUpAccount) error {
	account, err := au.accountRepository.GetAccountsByIdUserId(topUpAccount.UserID, topUpAccount.ID)
	if err != nil || account == nil || account.ID == common.STRING_EMPTY {
		return errors.New("record not found")
	}

	demoAccountTopUp := model.DemoAccountTopUp{
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

func (au *accountUsecase) BackOfficeAll(request model.BackOfficeAllAccountRequest) (response model.BackOfficeAllAccountResponse, err error) {
	response.Pagination = request.Pagination

	accounts, err := au.accountRepository.GetBackOfficeAllAccounts(request)
	if err != nil {
		return response, errors.New("record not found")
	}

	response.Data = accounts

	return response, nil
}

func (au *accountUsecase) BackOfficePending(request model.BackOfficePendingAccountRequest) (response model.BackOfficePendingAccountResponse, err error) {
	response.Pagination = request.Pagination

	accounts, err := au.accountRepository.GetBackOfficePendingAccounts(request)
	if err != nil {
		return response, errors.New("record not found")
	}

	response.Data = accounts

	return response, nil
}

func (au *accountUsecase) BackOfficePendingApproval(backOfficePendingApproval model.BackOfficePendingAccountApprovalRequest) error {
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
			clientAdd := model.ClientAdd{
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

func (au *accountUsecase) BackOfficeAccountByMenuType(request model.BackOfficeAccountByMenuTypeRequest) (response model.BackOfficeAccountByMenuTypeResponse, err error) {
	switch request.MenuType {
	case common.SPA, common.Multi:
		response.Pagination = request.Pagination
		accounts, err := au.accountRepository.GetBackOfficeAccountByFilterParams(model.BackOfficeAccountByFilterParams{
			Type:           request.Type,
			ApprovalStatus: request.ApprovalStatus,
			Pagination:     request.Pagination,
		})
		if err != nil {
			return response, errors.New("record not found")
		}
		response.Data = accounts
	default:
		return response, errors.New("invalid menu type")
	}

	return response, nil
}
