package market

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"

	"gorm.io/gorm"
)

type marketRepository struct {
	db *gorm.DB
}

func NewMarketRepository(db *gorm.DB) MarketRepository {
	return &marketRepository{db: db}
}

func (mr *marketRepository) CreateMarket(market *Market) error {
	return mr.db.Create(market).Error
}

func (mr *marketRepository) SaveMarket(market Market) error {
	return mr.db.Save(market).Error
}

func (mr *marketRepository) SaveMarkets(markets []Market) error {
	return mr.db.Save(markets).Error
}

func (mr *marketRepository) UpdateMarket(market *Market) error {
	return mr.db.Updates(market).Error
}

func (mr *marketRepository) UpdateMarketIsWatchList(market *Market) error {
	return mr.db.Select("IsWatchList").Updates(market).Error
}

func (mr *marketRepository) UpdateMarkets(market []Market) error {
	return mr.db.Updates(market).Error
}

func (mr *marketRepository) GetMarketByCode(code string) (*Market, error) {
	var market Market
	if err := mr.db.Where(&Market{Code: code}).First(&market).Error; err != nil {
		return nil, err
	}

	return &market, nil
}

func (mr *marketRepository) GetActiveMarketByCode(code string) (*Market, error) {
	var market Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarketIsWatchList() (*Market, error) {
	var market Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarketsIsWatchList(sessionId string) ([]Market, error) {
	var market []Market
	err := mr.db.Model(&Market{}).
		Joins("JOIN users ON markets.code = ANY(users.watchlist)").
		Where("users.session_id = ? AND users.is_active = ?", sessionId, true).
		Scan(&market).Error

	if err != nil || market == nil {
		return nil, err
	}

	return market, nil
}

func (mr *marketRepository) GetActiveMarketIsForeignExchange() (*Market, error) {
	var market Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarketsIsForeignExchange() ([]Market, error) {
	var market []Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarketIsMetals() (*Market, error) {
	var market Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarketsIsMetals() ([]Market, error) {
	var market []Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarketIsUsIndex() (*Market, error) {
	var market Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarketsIsUsIndex() ([]Market, error) {
	var market []Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarketIsAsiaIndex() (*Market, error) {
	var market Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarketsIsAsiaIndex() ([]Market, error) {
	var market []Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarketIsOil() (*Market, error) {
	var market Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarketsIsOil() ([]Market, error) {
	var market []Market
	queryParams := Market{
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

func (mr *marketRepository) GetActiveMarkets() ([]Market, error) {
	var markets []Market
	queryParams := Market{
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := mr.db.Find(&markets, &queryParams).Error; err != nil {
		return nil, err
	}

	return markets, nil
}

func (mr *marketRepository) GetActiveMarketCountries() ([]MarketCountry, error) {
	var marketCountries []MarketCountry
	queryParams := MarketCountry{
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}
	if err := mr.db.Find(&marketCountries, &queryParams).Error; err != nil {
		return nil, err
	}

	return marketCountries, nil
}

func (mr *marketRepository) GetMarketCurrencyRate(currencyBase, currencyQuote string) (*MarketCurrencyRate, error) {
	var marketCurrencyRate MarketCurrencyRate
	queryParams := MarketCurrencyRate{
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
