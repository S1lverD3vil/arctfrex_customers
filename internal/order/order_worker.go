package order

import (
	"log"
	"time"
)

type OrderWorker struct {
	orderUsecase OrderUsecase
}

func NewOrderWorker(ou OrderUsecase) *OrderWorker {
	return &OrderWorker{orderUsecase: ou}
}

func (ow *OrderWorker) CloseAllExpiredOrder(interval time.Duration) {
	go func() {
		for {
			err := ow.orderUsecase.CloseAllExpiredOrder()
			if err != nil {
				log.Println("Error close all expired order:", err)
			}

			// log.Printf("News updated")
			time.Sleep(interval)
		}
	}()
}
