package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
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
	BackOfficeQuestions(ctx context.Context, userID string, accountID string) (response model.SurveyResponse, err error)
	BackOfficeAnswers(ctx context.Context, request model.SurveyAnswerRequest) (response model.SurveyAnswerResponse, err error)
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
	account.NoAggreement = common.GenerateShortID("PAN")

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

func (au *accountUsecase) BackOfficeQuestions(ctx context.Context, userID string, accountID string) (response model.SurveyResponse, err error) {
	response, err = au.GetQuestions()
	if err != nil {
		return response, err
	}

	userData, err := au.accountRepository.GetBackOfficeAccountQuestions(userID, accountID)
	if err != nil {
		return response, err
	}

	surveyData := model.SurveyDataTemplate{
		PhoneNumber:            userData.PhoneNumber,
		NoIdentity:             userData.NoIdentity,
		Address:                userData.Address,
		DomProvince:            userData.DomProvince,
		Profession:             userData.Profession,
		MotherName:             userData.MotherName,
		Email:                  userData.Email,
		Name:                   userData.Name,
		AccountType:            userData.AccountType,
		CurrencyRate:           userData.CurrencyRate,
		ProductServicePlatform: userData.ProductServicePlatform,
		AccountID:              userData.AccountID,
		Banks:                  userData.BankList,
		SurveyResult:           userData.SurveyResult,
	}

	for i, survey := range response.Data.SurveyChecklist {
		rendered := ""
		if survey.No == 19 {
			tpl := `Rekening {{.Index}}:<br/><b>{{.BankName}}<br/>No Rek {{.BankAccountNumber}}<br/>atas {{.BankBeneficiaryName}}</b><br/><br/>`

			var renderedBanks strings.Builder
			for index, bank := range surveyData.Banks {
				data := model.Banklist{
					Index:               index + 1,
					BankName:            bank.BankName,
					BankAccountNumber:   bank.BankAccountNumber,
					BankBeneficiaryName: bank.BankBeneficiaryName,
				}

				rendered, err := renderTemplateString(tpl, data)
				if err != nil {
					return response, err
				}
				renderedBanks.WriteString(rendered)
			}
			survey.Question = strings.ReplaceAll(survey.Question, "{{bank_list}}", renderedBanks.String())
		} else {
			rendered, err = renderTemplateString(survey.Question, surveyData)
			if err != nil {
				return response, err
			}
		}

		response.Data.SurveyChecklist[i].Question = rendered
		if i < len(surveyData.SurveyResult.SurveyChecklist) && survey.No == surveyData.SurveyResult.SurveyChecklist[i].No {
			response.Data.SurveyChecklist[i].Answer = surveyData.SurveyResult.SurveyChecklist[i].Answer
		}
	}

	for i, survey := range response.Data.PpatkChecklist {
		rendered, err := renderTemplateString(survey.Question, surveyData)
		if err != nil {
			return response, err
		}
		response.Data.PpatkChecklist[i].Question = rendered
		if i < len(surveyData.SurveyResult.PpatkChecklist) && survey.No == surveyData.SurveyResult.PpatkChecklist[i].No {
			response.Data.PpatkChecklist[i].Answer = surveyData.SurveyResult.PpatkChecklist[i].Answer
		}
	}

	return response, nil
}

func renderTemplateString(tpl string, data any) (string, error) {
	tmpl, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}

func (*accountUsecase) GetQuestions() (response model.SurveyResponse, err error) {
	path := os.Getenv("SURVEY_QUESTIONS_PATH")
	if path == "" {
		return response, errors.New("environment variable SURVEY_QUESTIONS_PATH is not set")
	}

	cwd, _ := os.Getwd()
	directoryTemplate := filepath.Join(cwd, path)
	if filepath.Base(cwd) == "cmd" {
		directoryTemplate = filepath.Join(cwd, "../", path)
	}

	jsonFile, err := os.Open(directoryTemplate)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	// Read file
	byteValue, _ := io.ReadAll(jsonFile)

	// Unmarshal to struct
	var result model.SurveyPayload
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		panic(err)
	}

	response.Data.SurveyChecklist = result.SurveyChecklist
	response.Data.PpatkChecklist = result.PpatkChecklist

	return response, err
}

func (au *accountUsecase) BackOfficeAnswers(ctx context.Context, request model.SurveyAnswerRequest) (response model.SurveyAnswerResponse, err error) {
	err = au.accountRepository.UpdateSurveyResult(&model.Account{
		ID:     request.AccountID,
		UserID: request.UserID,
		SurveyResult: model.SurveyResult{
			SurveyChecklist: request.SurveyChecklist,
			PpatkChecklist:  request.PpatkChecklist,
		},
	})
	if err != nil {
		return response, errors.New("failed to update survey result")
	}

	// Prepare response data
	response = model.SurveyAnswerResponse{
		SurveyChecklist: request.SurveyChecklist,
		PpatkChecklist:  request.PpatkChecklist,
	}

	return
}
