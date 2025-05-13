package report

import (
	"fmt"
	"reflect"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/model"
	"arctfrex-customers/internal/repository"
)

type ReportUsecase interface {
	GetActiveReports() (*[]Report, error)
	GetActiveReportsByCode(reportCode, startDate, endDate string) (*ReportData, error)
	GroupUserLoginsUpdates() error
}

type reportUsecase struct {
	reportRepository  ReportRepository
	accountRepository repository.AccountRepository
	depositRepository repository.DepositRepository
	reportApiClient   ReportApiClient
}

func NewReportUsecase(
	rr ReportRepository,
	ar repository.AccountRepository,
	dr repository.DepositRepository,
	rac ReportApiClient,
) *reportUsecase {
	return &reportUsecase{
		reportRepository:  rr,
		accountRepository: ar,
		depositRepository: dr,
		reportApiClient:   rac,
	}
}

func (ru *reportUsecase) GetActiveReports() (*[]Report, error) {
	reports, err := ru.reportRepository.GetActiveReports()
	if err != nil {
		return &[]Report{}, err
	}

	return reports, nil
}

func (ru *reportUsecase) GetActiveReportsByCode(reportCode, startDate, endDate string) (*ReportData, error) {
	reportData := ReportData{
		Code: reportCode,
	}

	switch reportCode {
	case "R_DT":
		{
			deposits, err := ru.depositRepository.GetBackOfficePendingDeposits()
			if len(*deposits) == 0 {
				reportData.Data = []interface{}{}
				return &reportData, err
			}

			reportData.Data = deposits
		}
	case "R_MANIFEST":
		{
			accountGetManifestResponse, err := ru.reportApiClient.GetAccountsManifest()
			if err != nil {
				fmt.Println("Error getting account manifest:", err)
				return &ReportData{}, err
			}

			reportData.Data = []interface{}{*accountGetManifestResponse}

			// Get the type of the struct
			typ := reflect.TypeOf(AccountGetManifestResponse{})

			// Iterate over the fields
			var columns []string
			for i := 0; i < typ.NumField(); i++ {
				columns = append(columns, typ.Field(i).Name)
			}
			reportData.Column = columns
		}
	case "R_PL":
		{
			reportProfitLossData, err := ru.accountRepository.GetReportProfitLoss(startDate, endDate)
			if len(*reportProfitLossData) == 0 {
				reportData.Data = []interface{}{}
				return &reportData, err
			}

			var columns []string
			typ := reflect.TypeOf(model.ReportProfitLossData{})
			for i := 0; i < typ.NumField(); i++ {
				columns = append(columns, typ.Field(i).Name)
			}

			// Map to JSON response struct
			response := make([]ReportProfitLossDataResponse, len(*reportProfitLossData))
			for i, data := range *reportProfitLossData {
				response[i] = mapToResponse(data)
			}

			// Convert `response` to `[]interface{}`
			dataInterfaces := make([]interface{}, len(response))
			for i, v := range response {
				dataInterfaces[i] = v
			}

			reportData.Column = columns
			reportData.Data = dataInterfaces
		}
	default:
		{
			deposits, err := ru.depositRepository.GetBackOfficePendingDeposits()

			if len(*deposits) == 0 {
				reportData.Data = []interface{}{}
				return &reportData, err
			}

			reportData.Data = deposits
		}
	}

	return &reportData, nil
}

func (ru *reportUsecase) GroupUserLoginsUpdates() error {

	groupUserLoginsUpdates, err := ru.reportApiClient.GetGroupUserLogins(GroupUserLoginsRequest{GroupName: "demo\\PKB\\B-USD-SFL15-MAR-C50-SWAP"})

	if err != nil {
		return err
	}

	if err := ru.reportRepository.SaveGroupUserLogins(ConvertResponseToReport(*groupUserLoginsUpdates)); err != nil {
		return err
	}

	return nil
}

func mapToResponse(data model.ReportProfitLossData) ReportProfitLossDataResponse {
	return ReportProfitLossDataResponse{
		MetaLoginID:           data.MetaLoginID,
		Name:                  data.Name,
		DomCity:               data.DomCity,
		Currency:              data.Currency,
		CurrencyRate:          data.CurrencyRate,
		TotalDepositAmount:    data.TotalDepositAmount,
		TotalWithdrawalAmount: data.TotalWithdrawalAmount,
		PrevEquity:            data.PrevEquity,
		Nmii:                  data.Nmii,
		LastEquity:            data.LastEquity,
		GrossProfit:           data.GrossProfit,
		GrossProfitUSD:        data.GrossProfitUSD,
		SingleSideLot:         data.SingleSideLot,
		Commission:            data.Commission,
		Rebate:                data.Rebate,
		PrevBadDebt:           data.PrevBadDebt,
		LastBadDebt:           data.LastBadDebt,
		NetProfit:             data.NetProfit,
		NetProfitUSD:          data.NetProfitUSD,
		AccountID:             data.AccountID,
		UserID:                data.UserID,
	}
}

func ConvertResponseToReport(response GroupUserLoginsResponse) []ReportGroupUserLogins {
	// Use slices.Map to transform data slice to []ReportGroupUserLogins
	var result []ReportGroupUserLogins
	for _, id := range response.Data {
		result = append(result, ReportGroupUserLogins{Login: id, BaseModel: base.BaseModel{IsActive: true}})
	}
	return result
}
