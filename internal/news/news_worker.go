package news

import (
	"log"
	"time"
)

type NewsWorker struct {
	newsUsecase NewsUsecase
}

func NewNewsWorker(nu NewsUsecase) *NewsWorker {
	return &NewsWorker{newsUsecase: nu}
}

func (mw *NewsWorker) NewsUpdates(interval time.Duration) {
	go func() {
		for {
			err := mw.newsUsecase.NewsUpdates()
			if err != nil {
				log.Println("Error fetching news latest updates:", err)
			}

			//log.Printf("News updated")
			time.Sleep(interval)
		}
	}()
}

func (mw *NewsWorker) NewsBulletinUpdates(interval time.Duration) {
	go func() {
		for {
			err := mw.newsUsecase.NewsBulletinUpdates()
			if err != nil {
				log.Println("Error fetching news bulletin updates:", err)
			}

			//log.Printf("News Bulletin updated")
			time.Sleep(interval)
		}
	}()
}
