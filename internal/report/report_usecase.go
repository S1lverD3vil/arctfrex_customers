package report

import (
	"arctfrex-customers/internal/account"
	"arctfrex-customers/internal/deposit"
	"fmt"
	"reflect"
)

type ReportUsecase interface {
	GetActiveReports() (*[]Report, error)
	GetActiveReportsByCode(reportCode string) (*ReportData, error)
}

type reportUsecase struct {
	reportRepository  ReportRepository
	accountRepository account.AccountRepository
	depositRepository deposit.DepositRepository
	reportApiClient   ReportApiClient
}

func NewReportUsecase(
	rr ReportRepository,
	ar account.AccountRepository,
	dr deposit.DepositRepository,
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

func (ru *reportUsecase) GetActiveReportsByCode(reportCode string) (*ReportData, error) {
	reportData := ReportData{
		Code: reportCode,
	}

	switch reportCode {
	case "R_DT":
		{
			deposits, err := ru.depositRepository.GetBackOfficePendingDeposits()
			// log.Println("data deposit")
			//log.Println(*deposits)
			// if deposits == nil || err != nil {
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

			var columns []string
			var items []AccountGetManifestResponse
			items = append(items, *accountGetManifestResponse)

			// Create an instance of the struct
			p := AccountGetManifestResponse{}

			// Get the type of the struct
			typ := reflect.TypeOf(p)

			// Iterate over the fields
			for i := 0; i < typ.NumField(); i++ {
				columns = append(columns, typ.Field(i).Name)
			}

			reportData.Column = columns
			reportData.Data = items
		}
	case "R_PL":
		{
			reportProfitLossData, err := ru.accountRepository.GetReportProfitLoss()
			if len(*reportProfitLossData) == 0 {
				reportData.Data = []interface{}{}
				return &reportData, err
			}

			var columns []string
			// var items []account.ReportProfitLossData
			// items = append(items, *reportProfitLossData...)

			// Create an instance of the struct
			p := account.ReportProfitLossData{}

			// Get the type of the struct
			typ := reflect.TypeOf(p)

			// Iterate over the fields
			for i := 0; i < typ.NumField(); i++ {
				columns = append(columns, typ.Field(i).Name)
			}
			// Map to JSON response struct
			response := make([]ReportProfitLossDataResponse, len(*reportProfitLossData))
			for i, data := range *reportProfitLossData {
				response[i] = mapToResponse(data)
			}

			reportData.Column = columns
			reportData.Data = response
		}
	default:
		{
			deposits, err := ru.depositRepository.GetBackOfficePendingDeposits()
			// log.Println("default data deposit")
			// log.Println(deposits)
			// if deposits == nil || err != nil {
			// 	reportData.Data = []string{}
			// 	return &reportData, err
			// }
			if len(*deposits) == 0 {
				reportData.Data = []interface{}{}
				return &reportData, err
			}

			reportData.Data = deposits
		}
	}
	//reports, err := ru.reportRepository.GetActiveReports()
	// if err != nil {
	// 	return &ReportData{}, err
	// }

	return &reportData, nil
}

func mapToResponse(data account.ReportProfitLossData) ReportProfitLossDataResponse {
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
