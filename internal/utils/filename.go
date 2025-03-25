package utils

import "fmt"

func GenerateReportFilename(reportCode, startDate, endDate string) string {
	switch {
	case startDate == "" && endDate == "":
		return fmt.Sprintf("report-%s-all", reportCode)
	case startDate == "":
		return fmt.Sprintf("report-%s-until-%s", reportCode, endDate)
	case endDate == "":
		return fmt.Sprintf("report-%s-from-%s", reportCode, startDate)
	default:
		return fmt.Sprintf("report-%s-%s-to-%s", reportCode, startDate, endDate)
	}
}
