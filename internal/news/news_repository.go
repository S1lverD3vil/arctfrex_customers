package news

import (
	"arctfrex-customers/internal/base"

	"gorm.io/gorm"
)

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (nr *newsRepository) GetNewsList() (*[]News, error) {
	var news []News
	if err := nr.db.Find(&news).Error; err != nil {
		return &[]News{}, err
	}

	return &news, nil
}

func (nr *newsRepository) GetNewsByTitle(title string) (News, error) {
	var news News
	if err := nr.db.Where(&News{NewsTitle: title}).First(&news).Error; err != nil {
		return news, err
	}

	return news, nil

}

func (nr *newsRepository) UpdateNews(title string, news *News) error {
	return nr.db.Model(&News{}).Where("news_title = ?", title).Updates(&news).Error
}

func (nr *newsRepository) DeactiveNews(title string) error {
	return nr.db.Model(&News{}).Where("news_title = ?", title).Update("is_active", false).Error
}

func (nr *newsRepository) SaveNews(news *News) error {
	return nr.db.Save(news).Error
}

func (nr *newsRepository) SaveNewsBulletin(newsBulletin *NewsBulletin) error {
	return nr.db.Save(newsBulletin).Error
}

func (nr *newsRepository) GetActiveNews() ([]News, error) {
	var news []News
	queryParams := News{
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := nr.db.Find(&news, &queryParams).Error
	if err != nil || news == nil {
		return nil, err
	}

	return news, nil
}

func (nr *newsRepository) GetActiveNewsBulletin() ([]NewsBulletin, error) {
	var newsBulletin []NewsBulletin
	queryParams := NewsBulletin{
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := nr.db.Find(&newsBulletin, &queryParams).Error
	if err != nil || newsBulletin == nil {
		return nil, err
	}

	return newsBulletin, nil
}
