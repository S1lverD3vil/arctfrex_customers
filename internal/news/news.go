package news

import "arctfrex-customers/internal/base"

// Models
type News struct {
	NewsTitle      string `gorm:"primary_key" json:"news_title"`
	NewsLink       string `json:"news_link"`
	NewsImageLink  string `json:"news_image_link"`
	NewsImpact     string `json:"news_impact"`
	NewsSource     string `json:"news_source"`
	NewsSourceLink string `json:"news_source_link"`
	NewsTimeAgo    string `json:"news_time_ago"`
	NewsPreview    string `json:"news_preview"`
	NewsType       string `json:"news_type"`
	NewsCategory   string `json:"news_category"`
	NewsTicker     string `json:"news_ticker"`
	NewsDate       string `json:"news_date"`

	base.BaseModel
}
type NewsBulletin struct {
	NewsBulletinTitle      string `gorm:"primary_key" json:"news_bulletin_title"`
	NewsBulletinLink       string `json:"news_bulletin_link"`
	NewsBulletinImageLink  string `json:"news_bulletin_image_link"`
	NewsBulletinImpact     string `json:"news_bulletin_impact"`
	NewsBulletinSource     string `json:"news_bulletin_source"`
	NewsBulletinSourceLink string `json:"news_bulletin_source_link"`
	NewsBulletinTimeAgo    string `json:"news_bulletin_time_ago"`
	NewsBulletinPreview    string `json:"news_bulletin_preview"`
	NewsBulletinType       string `json:"news_bulletin_type"`
	NewsBulletinCategory   string `json:"news_bulletin_category"`
	NewsBulletinTicker     string `json:"news_bulletin_ticker"`
	NewsBulletinDate       string `json:"news_bulletin_date"`

	base.BaseModel
}

// Dtos
type NewsResponse struct {
	Data []NewsData `json:"data"`
}

type NewsData struct {
	NewsURL   string   `json:"news_url"`
	ImageURL  string   `json:"image_url"`
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	Source    string   `json:"source_name"`
	Date      string   `json:"date"`
	Topics    []string `json:"topics"`
	Sentiment string   `json:"sentiment"`
	Type      string   `json:"type"`
	Currency  []string `json:"currency"`
}

type NewsBulletinResponse struct {
	Results []Article `json:"results"`
}

type Article struct {
	ArticleID      string   `json:"article_id"`
	Title          string   `json:"title"`
	Link           string   `json:"link"`
	Keywords       []string `json:"keywords"`
	Creator        []string `json:"creator"`
	VideoURL       *string  `json:"video_url"` // Using pointer to handle null values
	Description    string   `json:"description"`
	Content        string   `json:"content"`
	PubDate        string   `json:"pubDate"`
	PubDateTZ      string   `json:"pubDateTZ"`
	ImageURL       string   `json:"image_url"`
	SourceID       string   `json:"source_id"`
	SourcePriority int      `json:"source_priority"`
	SourceName     string   `json:"source_name"`
	SourceURL      string   `json:"source_url"`
	SourceIcon     string   `json:"source_icon"`
	Language       string   `json:"language"`
	Country        []string `json:"country"`
	Category       []string `json:"category"`
	AITag          string   `json:"ai_tag"`
	Sentiment      string   `json:"sentiment"`
	SentimentStats string   `json:"sentiment_stats"`
	AIRegion       string   `json:"ai_region"`
	AIOrg          string   `json:"ai_org"`
	Duplicate      bool     `json:"duplicate"`
}

type NewsRepository interface {
	GetNewsList() (*[]News, error)
	SaveNews(news *News) error
	SaveNewsBulletin(newsBulletin *NewsBulletin) error
	UpdateNews(title string, news *News) error
	DeactiveNews(title string) error
	GetNewsByTitle(title string) (News, error)
	GetActiveNews() ([]News, error)
	GetActiveNewsBulletin() ([]NewsBulletin, error)
}
