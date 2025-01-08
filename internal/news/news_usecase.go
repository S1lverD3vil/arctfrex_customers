package news

import (
	"arctfrex-customers/internal/base"
	"errors"
)

type NewsUsecase interface {
	GetNewsList() (*[]News, error)
	CreateNews(*News) error
	GetNewsByTitle(title string) (News, error)
	UpdateNews(title string, new *News) error
	DeactiveNews(title string) error
	GetNewsLatest() ([]News, error)
	GetNewsLatestBulletin() ([]NewsBulletin, error)
	NewsUpdates() error
	NewsBulletinUpdates() error
}

type newsUsecase struct {
	newsRepository NewsRepository
	newsApiclient  NewsApiclient
}

func NewNewsUsecase(nr NewsRepository, na NewsApiclient) NewsUsecase {
	return &newsUsecase{newsRepository: nr, newsApiclient: na}
}

func (nu *newsUsecase) GetNewsList() (*[]News, error) {
	news, err := nu.newsRepository.GetNewsList()
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (nu *newsUsecase) CreateNews(news *News) error {
	err := nu.newsRepository.SaveNews(news)
	if err != nil {
		return err
	}

	return nil
}

func (nu *newsUsecase) GetNewsByTitle(title string) (News, error) {
	news, err := nu.newsRepository.GetNewsByTitle(title)
	if err != nil {
		return news, err
	}
	return news, nil
}

func (nu *newsUsecase) GetNewsLatest() ([]News, error) {
	news, err := nu.newsRepository.GetActiveNews()
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (nu *newsUsecase) UpdateNews(title string, news *News) error {
	if _, err := nu.newsRepository.GetNewsByTitle(title); err != nil {
		return err
	}

	err := nu.newsRepository.UpdateNews(title, news)
	if err != nil {
		return err
	}

	return nil
}

func (nu *newsUsecase) DeactiveNews(title string) error {
	if _, err := nu.newsRepository.GetNewsByTitle(title); err != nil {
		return err
	}

	err := nu.newsRepository.DeactiveNews(title)
	if err != nil {
		return err
	}

	return nil
}

func (nu *newsUsecase) GetNewsLatestBulletin() ([]NewsBulletin, error) {
	newsBulletin, err := nu.newsRepository.GetActiveNewsBulletin()
	if err != nil {
		return nil, err
	}

	return newsBulletin, nil
}

func (nu *newsUsecase) NewsUpdates() error {
	latestNewsUpdates, err := nu.newsApiclient.GetLatestNews()
	if err != nil {
		return err
	}

	if len(latestNewsUpdates.Data) < 1 {
		return errors.New("not found")
	}

	for _, data := range latestNewsUpdates.Data {

		news := &News{
			NewsTitle:     data.Title,
			NewsLink:      data.NewsURL,
			NewsImageLink: data.ImageURL,
			NewsPreview:   data.Text,
			BaseModel: base.BaseModel{
				IsActive: true,
			},
		}

		if err := nu.newsRepository.SaveNews(news); err != nil {
			return err
		}
	}

	return nil
}

func (nu *newsUsecase) NewsBulletinUpdates() error {
	latestBulletinNewsUpdates, err := nu.newsApiclient.GetLatestNewsBulletin()
	if err != nil {
		return err
	}

	if len(latestBulletinNewsUpdates.Results) < 1 {
		return errors.New("not found")
	}

	for _, data := range latestBulletinNewsUpdates.Results {

		newsBulletin := &NewsBulletin{
			NewsBulletinTitle:     data.Title,
			NewsBulletinLink:      data.SourceURL,
			NewsBulletinImageLink: data.ImageURL,
			NewsBulletinPreview:   data.Description,
			BaseModel: base.BaseModel{
				IsActive: true,
			},
		}

		if err := nu.newsRepository.SaveNewsBulletin(newsBulletin); err != nil {
			return err
		}
	}

	return nil
}
