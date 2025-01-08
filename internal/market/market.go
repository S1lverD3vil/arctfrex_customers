package market

import (
	"arctfrex-customers/internal/base"
)

type Market struct {
	Code              string  `gorm:"primary_key" json:"code"`
	BaseCurrency      string  `json:"base_currency"`
	QuoteCurrency     string  `json:"quote_currency"`
	Price             float64 `json:"price"`
	Ask               float64 `json:"ask"`
	Bid               float64 `json:"bid"`
	Mid               float64 `json:"mid"`
	Change            float64 `json:"change"`
	ChangePercentage  float64 `json:"change_percentage"`
	DayLow            float64 `json:"day_low"`
	DayHigh           float64 `json:"day_high"`
	YearLow           float64 `json:"year_low"`
	YearHigh          float64 `json:"year_high"`
	Avg50Price        float64 `json:"avg50_price"`
	Avg200Price       float64 `json:"avg200_price"`
	IsWatchList       bool    `json:"is_watch_list"`
	IsForeignExchange bool    `json:"is_foreign_exchange"`
	IsMetals          bool    `json:"is_metals"`
	IsUsIndex         bool    `json:"is_us_index"`
	IsAsiaIndex       bool    `json:"is_asia_index"`
	IsOil             bool    `json:"is_oil"`

	base.BaseModel
}

type MarketCountry struct {
	CountryCode  string  `json:"country_code"`
	CurrencyCode string  `gorm:"primary_key" json:"currency_code"`
	CurrencyRate float64 `json:"currency_rate"`
	Description  string  `json:"description"`
	FlagImage    string  `json:"flag_image"`

	base.BaseModel
}

type MarketCurrencyRate struct {
	BaseCurrency  string  `json:"base_currency"`
	QuoteCurrency string  `json:"quote_currency"`
	Rate          float64 `json:"rate"`

	base.BaseModel
}

type ConvertPrice struct {
	AmountBase    float64 `json:"amount_base"`
	AmountQuote   float64 `json:"amount_quote"`
	CurrencyBase  string  `json:"currency_base"`
	CurrencyQuote string  `json:"currency_quote"`
}

type MarketPriceResponse struct {
	Code               string  `json:"code"`
	Description        string  `json:"description"`
	Ask                float64 `json:"ask"`
	Bid                float64 `json:"bid"`
	Low                float64 `json:"low"`
	High               float64 `json:"high"`
	Price              float64 `json:"price"`
	Change             float64 `json:"change"`
	ChangePercentage   float64 `json:"change_percentage"`
	BaseCurrencyImage  string  `json:"base_currency_image"`
	QuoteCurrencyImage string  `json:"quote_currency_image"`
	IsWatchList        bool    `json:"is_watch_list"`
}

type WebSocketRequest struct {
	FilterBy   string `json:"filter_by"`
	MarketCode string `json:"market_code"`
	SessionId  string `json:"session_id"`
}

type WebSocketResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    string `json:"time"`
}

type ForexPriceResponse struct {
	Ticker           string  `json:"ticker"`
	BaseCurrency     string  `json:"baseCurrency"`
	QuoteCurrency    string  `json:"quoteCurrency"`
	Price            float64 `json:"price"`
	Change           float64 `json:"change"`
	ChangePercentage float64 `json:"changePercentage"`
	DayLow           float64 `json:"dayLow"`
	DayHigh          float64 `json:"dayHigh"`
	YearLow          float64 `json:"yearLow"`
	YearHigh         float64 `json:"yearHigh"`
	Avg50Price       float64 `json:"avg50Price"`
	Avg200Price      float64 `json:"avg200Price"`
}

// Define the struct for the API response
type LiveMarketUpdatesResponse struct {
	Endpoint      string                   `json:"endpoint"`
	Quotes        []LiveMarketUpdatesQuote `json:"quotes"`
	RequestedTime string                   `json:"requested_time"`
	Timestamp     int64                    `json:"timestamp"`
}

// Define the struct for each quote
type LiveMarketUpdatesQuote struct {
	Ask           float64 `json:"ask"`
	BaseCurrency  string  `json:"base_currency"`
	Bid           float64 `json:"bid"`
	Mid           float64 `json:"mid"`
	QuoteCurrency string  `json:"quote_currency"`
}

// Define the struct for the entire API response
type HistoricalMarketUpdatesResponse struct {
	Date        string                         `json:"date"`
	Endpoint    string                         `json:"endpoint"`
	Quotes      []HistoricalMarketUpdatesQuote `json:"quotes"`
	RequestTime string                         `json:"request_time"`
}

// Define the struct for each quote
type HistoricalMarketUpdatesQuote struct {
	BaseCurrency  string  `json:"base_currency"`
	Close         float64 `json:"close"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Open          float64 `json:"open"`
	QuoteCurrency string  `json:"quote_currency"`
}

type ArcMetaIntegratorPrice struct {
	Data  []ArcMetaIntegratorPriceData `json:"data"`
	Error string                       `json:"error"`
}

type ArcMetaIntegratorPriceData struct {
	Symbol      string  `json:"Symbol"`
	Digits      int     `json:"Digits"`
	Bid         float64 `json:"Bid"`
	Ask         float64 `json:"Ask"`
	Last        float64 `json:"Last"`
	Volume      int     `json:"Volume"`
	VolumeReal  float64 `json:"VolumeReal"`
	Datetime    int64   `json:"Datetime"`
	DatetimeMsc int64   `json:"DatetimeMsc"`
}

// MarketRepository defines the interface for updating Forex prices
type MarketRepository interface {
	CreateMarket(market *Market) error
	SaveMarket(market Market) error
	SaveMarkets(markets []Market) error
	UpdateMarket(market *Market) error
	UpdateMarketIsWatchList(market *Market) error
	UpdateMarkets(market []Market) error
	GetMarketByCode(code string) (*Market, error)
	GetActiveMarketByCode(code string) (*Market, error)
	GetActiveMarketIsWatchList() (*Market, error)
	//GetActiveMarketsIsWatchList() ([]Market, error)
	GetActiveMarketsIsWatchList(sessionId string) ([]Market, error)
	GetActiveMarketIsForeignExchange() (*Market, error)
	GetActiveMarketsIsForeignExchange() ([]Market, error)
	GetActiveMarketIsMetals() (*Market, error)
	GetActiveMarketsIsMetals() ([]Market, error)
	GetActiveMarketIsUsIndex() (*Market, error)
	GetActiveMarketsIsUsIndex() ([]Market, error)
	GetActiveMarketIsAsiaIndex() (*Market, error)
	GetActiveMarketsIsAsiaIndex() ([]Market, error)
	GetActiveMarketIsOil() (*Market, error)
	GetActiveMarketsIsOil() ([]Market, error)
	GetActiveMarkets() ([]Market, error)
	GetActiveMarketCountries() ([]MarketCountry, error)
	GetMarketCurrencyRate(currencyBase, currencyQuote string) (*MarketCurrencyRate, error)
}
