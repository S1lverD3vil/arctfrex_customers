package repository

import (
	"gorm.io/gorm"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/model"
)

type MarketRepository interface {
	CreateMarket(market *model.Market) error
	SaveMarket(market model.Market) error
	SaveMarkets(markets []model.Market) error
	UpdateMarket(market *model.Market) error
	UpdateMarketIsWatchList(market *model.Market) error
	UpdateMarkets(market []model.Market) error
	GetMarketByCode(code string) (*model.Market, error)
	GetActiveMarketByCode(code string) (*model.Market, error)
	GetActiveMarketIsWatchList() (*model.Market, error)
	//GetActiveMarketsIsWatchList() ([]Market, error)
	GetActiveMarketsIsWatchList(sessionId string) ([]model.Market, error)
	GetActiveMarketIsForeignExchange() (*model.Market, error)
	GetActiveMarketsIsForeignExchange() ([]model.Market, error)
	GetActiveMarketIsMetals() (*model.Market, error)
	GetActiveMarketsIsMetals() ([]model.Market, error)
	GetActiveMarketIsUsIndex() (*model.Market, error)
	GetActiveMarketsIsUsIndex() ([]model.Market, error)
	GetActiveMarketIsAsiaIndex() (*model.Market, error)
	GetActiveMarketsIsAsiaIndex() ([]model.Market, error)
	GetActiveMarketIsOil() (*model.Market, error)
	GetActiveMarketsIsOil() ([]model.Market, error)
	GetActiveMarkets() ([]model.Market, error)
	GetActiveMarketCountries() ([]model.MarketCountry, error)
	GetMarketCurrencyRate(currencyBase, currencyQuote string) (*model.MarketCurrencyRate, error)
}

type marketRepository struct {
	db *gorm.DB
}

func NewMarketRepository(db *gorm.DB) MarketRepository {
	return &marketRepository{db: db}
}

func (mr *marketRepository) CreateMarket(market *model.Market) error {
	return mr.db.Create(market).Error
}

func (mr *marketRepository) SaveMarket(market model.Market) error {
	return mr.db.Save(market).Error
}

func (mr *marketRepository) SaveMarkets(markets []model.Market) error {
	return mr.db.Save(markets).Error
}

func (mr *marketRepository) UpdateMarket(market *model.Market) error {
	return mr.db.Updates(market).Error
}

func (mr *marketRepository) UpdateMarketIsWatchList(market *model.Market) error {
	return mr.db.Select("IsWatchList").Updates(market).Error
}

func (mr *marketRepository) UpdateMarkets(market []model.Market) error {
	return mr.db.Updates(market).Error
}

func (mr *marketRepository) GetMarketByCode(code string) (*model.Market, error) {
	var market model.Market
	if err := mr.db.Where(&model.Market{Code: code}).First(&market).Error; err != nil {
		return nil, err
	}

	return &market, nil
}

func (mr *marketRepository) GetActiveMarketByCode(code string) (*model.Market, error) {
	var market model.Market
	queryParams := model.Market{
		Code: code,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market.Code == common.STRING_EMPTY {
		return nil, err
	}

	return &market, nil
}

func (mr *marketRepository) GetActiveMarketIsWatchList() (*model.Market, error) {
	var market model.Market
	queryParams := model.Market{
		IsWatchList: true,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market.Code == common.STRING_EMPTY {
		return nil, err
	}

	return &market, nil
}

func (mr *marketRepository) GetActiveMarketsIsWatchList(sessionId string) ([]model.Market, error) {
	var market []model.Market
	err := mr.db.Model(&model.Market{}).
		Joins("JOIN users ON markets.code = ANY(users.watchlist)").
		Where("users.session_id = ? AND users.is_active = ?", sessionId, true).
		Scan(&market).Error

	if err != nil || market == nil {
		return nil, err
	}

	return market, nil
}

func (mr *marketRepository) GetActiveMarketIsForeignExchange() (*model.Market, error) {
	var market model.Market
	queryParams := model.Market{
		IsForeignExchange: true,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market.Code == common.STRING_EMPTY {
		return nil, err
	}

	return &market, nil
}

func (mr *marketRepository) GetActiveMarketsIsForeignExchange() ([]model.Market, error) {
	var market []model.Market
	queryParams := model.Market{
		IsForeignExchange: true,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market == nil {
		return nil, err
	}

	return market, nil
}

func (mr *marketRepository) GetActiveMarketIsMetals() (*model.Market, error) {
	var market model.Market
	queryParams := model.Market{
		IsMetals: true,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market.Code == common.STRING_EMPTY {
		return nil, err
	}

	return &market, nil
}

func (mr *marketRepository) GetActiveMarketsIsMetals() ([]model.Market, error) {
	var market []model.Market
	queryParams := model.Market{
		IsMetals: true,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market == nil {
		return nil, err
	}

	return market, nil
}

func (mr *marketRepository) GetActiveMarketIsUsIndex() (*model.Market, error) {
	var market model.Market
	queryParams := model.Market{
		IsUsIndex: true,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market.Code == common.STRING_EMPTY {
		return nil, err
	}

	return &market, nil
}

func (mr *marketRepository) GetActiveMarketsIsUsIndex() ([]model.Market, error) {
	var market []model.Market
	queryParams := model.Market{
		IsUsIndex: true,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market == nil {
		return nil, err
	}

	return market, nil
}

func (mr *marketRepository) GetActiveMarketIsAsiaIndex() (*model.Market, error) {
	var market model.Market
	queryParams := model.Market{
		IsAsiaIndex: true,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market.Code == common.STRING_EMPTY {
		return nil, err
	}

	return &market, nil
}

func (mr *marketRepository) GetActiveMarketsIsAsiaIndex() ([]model.Market, error) {
	var market []model.Market
	queryParams := model.Market{
		IsAsiaIndex: true,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market == nil {
		return nil, err
	}

	return market, nil
}

func (mr *marketRepository) GetActiveMarketIsOil() (*model.Market, error) {
	var market model.Market
	queryParams := model.Market{
		IsOil: true,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market.Code == common.STRING_EMPTY {
		return nil, err
	}

	return &market, nil
}

func (mr *marketRepository) GetActiveMarketsIsOil() ([]model.Market, error) {
	var market []model.Market
	queryParams := model.Market{
		IsOil: true,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&market, &queryParams).Error
	if err != nil || market == nil {
		return nil, err
	}

	return market, nil
}

func (mr *marketRepository) GetActiveMarkets() ([]model.Market, error) {
	var markets []model.Market
	queryParams := model.Market{
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := mr.db.Find(&markets, &queryParams).Error; err != nil {
		return nil, err
	}

	return markets, nil
}

func (mr *marketRepository) GetActiveMarketCountries() ([]model.MarketCountry, error) {
	var marketCountries []model.MarketCountry
	queryParams := model.MarketCountry{
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := mr.db.Find(&marketCountries, &queryParams).Error; err != nil {
		return nil, err
	}

	return marketCountries, nil
}

func (mr *marketRepository) GetMarketCurrencyRate(currencyBase, currencyQuote string) (*model.MarketCurrencyRate, error) {
	var marketCurrencyRate model.MarketCurrencyRate
	queryParams := model.MarketCurrencyRate{
		BaseCurrency:  currencyBase,
		QuoteCurrency: currencyQuote,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := mr.db.Find(&marketCurrencyRate, &queryParams).Error
	if err != nil || marketCurrencyRate.BaseCurrency == common.STRING_EMPTY {
		return nil, err
	}

	return &marketCurrencyRate, nil
}
