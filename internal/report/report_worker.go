package report

import (
	"log"
	"time"
)

type ReportWorker struct {
	reportUsecase ReportUsecase
}

func NewReportWorker(ru ReportUsecase) *ReportWorker {
	return &ReportWorker{reportUsecase: ru}
}

func (rw *ReportWorker) GroupUserLoginsUpdates(interval time.Duration) {
	go func() {
		for {
			err := rw.reportUsecase.GroupUserLoginsUpdates()
			if err != nil {
				log.Println("Error generating report:", err)
			}

			time.Sleep(interval)
		}
	}()
}
