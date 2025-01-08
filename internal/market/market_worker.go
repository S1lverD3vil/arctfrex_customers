package market

import (
	"log"
	"time"
)

type MarketWorker struct {
	usecase MarketUsecase
}

func NewMarketWorker(uc MarketUsecase) *MarketWorker {
	return &MarketWorker{usecase: uc}
}

func (mw *MarketWorker) PriceUpdates(interval time.Duration) {
	go func() {
		for {
			err := mw.usecase.PriceUpdates()
			if err != nil {
				log.Println("Error fetching market price:", err)
			}

			//log.Printf("Latest market price updated")
			time.Sleep(interval)
		}
	}()
}

func (mw *MarketWorker) LiveMarketUpdates(interval time.Duration) {
	go func() {
		for {
			err := mw.usecase.LiveMarketUpdates()
			if err != nil {
				log.Println("Error fetching market price:", err)
			}

			//log.Printf("Live market price updated")
			time.Sleep(interval)
		}
	}()
}
