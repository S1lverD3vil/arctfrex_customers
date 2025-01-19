package market

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	user_mobile "arctfrex-customers/internal/user/mobile"
	"errors"
	"sort"
	"strings"
)

type MarketUsecase interface {
	Price(fwebSocketRequest WebSocketRequest) (*[]MarketPriceResponse, error)
	PriceUpdates() error
	LiveMarketUpdates() error
	UpdateWatchlist(userId string, market Market) error
	GetWatchlist(userId, marketCode string) error
	ConvertPrice(convertPrice ConvertPrice) (*ConvertPrice, error)
}

type marketUsecase struct {
	marketRepository MarketRepository
	marketApiclient  MarketApiclient
	userRepository   user_mobile.UserRepository
}

func NewMarketUsecase(
	mr MarketRepository,
	ma MarketApiclient,
	ur user_mobile.UserRepository,
) MarketUsecase {
	return &marketUsecase{
		marketRepository: mr,
		marketApiclient:  ma,
		userRepository:   ur,
	}
}

func (mu *marketUsecase) Price(webSocketRequest WebSocketRequest) (*[]MarketPriceResponse, error) {
	var marketPriceResponses []MarketPriceResponse
	marketCountries, err := mu.marketRepository.GetActiveMarketCountries()
	if err != nil {
		return nil, errors.New("not found")
	}

	switch strings.ToLower(webSocketRequest.FilterBy) {
	case "code":
		{
			market, err := mu.marketRepository.GetActiveMarketByCode(webSocketRequest.MarketCode)
			if err != nil {
				return nil, errors.New("not found")
			}

			if market == nil {
				return &marketPriceResponses, nil

			}
			marketCountryBase, countryBaseFound := getMarketCountry(marketCountries, market.BaseCurrency)
			marketCountryQuote, countryQuoteFound := getMarketCountry(marketCountries, market.QuoteCurrency)

			marketPriceResponse := MarketPriceResponse{
				Code:  market.Code,
				Price: market.Price,
				Ask:   market.Ask,
				Bid:   market.Bid,
				// Ask:              common.RoundTo4DecimalPlaces(market.Ask),
				// Bid:              common.RoundTo4DecimalPlaces(market.Bid),
				Change:           market.Change,
				ChangePercentage: market.ChangePercentage,
				IsWatchList:      market.IsWatchList,
			}

			if countryBaseFound {
				marketPriceResponse.Description = marketCountryBase.Description
				marketPriceResponse.BaseCurrencyImage = marketCountryBase.FlagImage
			}

			if countryQuoteFound {
				marketPriceResponse.Description += " vs " + marketCountryQuote.Description
				marketPriceResponse.QuoteCurrencyImage = marketCountryQuote.FlagImage
			}

			marketPriceResponses = append(marketPriceResponses, marketPriceResponse)

			// Sort the slice based on the MarketCode field
			sort.Slice(marketPriceResponses, func(i, j int) bool {
				return marketPriceResponses[i].Code < marketPriceResponses[j].Code
			})

			return &marketPriceResponses, nil
		}
	case "watchlist":
		{
			markets, err := mu.marketRepository.GetActiveMarketsIsWatchList(webSocketRequest.SessionId)
			if err != nil {
				return nil, errors.New("not found")
			}

			for _, market := range markets {
				marketCountryBase, countryBaseFound := getMarketCountry(marketCountries, market.BaseCurrency)
				marketCountryQuote, countryQuoteFound := getMarketCountry(marketCountries, market.QuoteCurrency)

				marketPriceResponse := MarketPriceResponse{
					Code:             market.Code,
					Price:            market.Price,
					Ask:              market.Ask,
					Bid:              market.Bid,
					Change:           market.Change,
					ChangePercentage: market.ChangePercentage,
					IsWatchList:      market.IsWatchList,
				}

				if countryBaseFound {
					marketPriceResponse.Description = marketCountryBase.Description
					marketPriceResponse.BaseCurrencyImage = marketCountryBase.FlagImage
				}

				if countryQuoteFound {
					marketPriceResponse.Description += " vs " + marketCountryQuote.Description
					marketPriceResponse.QuoteCurrencyImage = marketCountryQuote.FlagImage
				}

				marketPriceResponses = append(marketPriceResponses, marketPriceResponse)
			}

			// Sort the slice based on the MarketCode field
			sort.Slice(marketPriceResponses, func(i, j int) bool {
				return marketPriceResponses[i].Code < marketPriceResponses[j].Code
			})

			return &marketPriceResponses, nil
		}
	case "foreign_exchange":
		{
			markets, err := mu.marketRepository.GetActiveMarketsIsForeignExchange()
			if err != nil {
				return nil, errors.New("not found")
			}

			for _, market := range markets {
				marketCountryBase, countryBaseFound := getMarketCountry(marketCountries, market.BaseCurrency)
				marketCountryQuote, countryQuoteFound := getMarketCountry(marketCountries, market.QuoteCurrency)

				marketPriceResponse := MarketPriceResponse{
					Code:             market.Code,
					Price:            market.Price,
					Ask:              market.Ask,
					Bid:              market.Bid,
					Change:           market.Change,
					ChangePercentage: market.ChangePercentage,
					IsWatchList:      market.IsWatchList,
				}

				if countryBaseFound {
					marketPriceResponse.Description = marketCountryBase.Description
					marketPriceResponse.BaseCurrencyImage = marketCountryBase.FlagImage
				}

				if countryQuoteFound {
					marketPriceResponse.Description += " vs " + marketCountryQuote.Description
					marketPriceResponse.QuoteCurrencyImage = marketCountryQuote.FlagImage
				}

				marketPriceResponses = append(marketPriceResponses, marketPriceResponse)
			}

			// Sort the slice based on the MarketCode field
			sort.Slice(marketPriceResponses, func(i, j int) bool {
				return marketPriceResponses[i].Code < marketPriceResponses[j].Code
			})

			return &marketPriceResponses, nil

		}
	case "metals":
		{
			markets, err := mu.marketRepository.GetActiveMarketsIsMetals()
			if err != nil {
				return nil, errors.New("not found")
			}

			for _, market := range markets {
				marketCountryBase, countryBaseFound := getMarketCountry(marketCountries, market.BaseCurrency)
				marketCountryQuote, countryQuoteFound := getMarketCountry(marketCountries, market.QuoteCurrency)

				marketPriceResponse := MarketPriceResponse{
					Code:             market.Code,
					Price:            market.Price,
					Ask:              market.Ask,
					Bid:              market.Bid,
					Change:           market.Change,
					ChangePercentage: market.ChangePercentage,
					IsWatchList:      market.IsWatchList,
				}

				if countryBaseFound {
					marketPriceResponse.Description = marketCountryBase.Description
					marketPriceResponse.BaseCurrencyImage = marketCountryBase.FlagImage
				}

				if countryQuoteFound {
					marketPriceResponse.Description += " vs " + marketCountryQuote.Description
					marketPriceResponse.QuoteCurrencyImage = marketCountryQuote.FlagImage
				}

				marketPriceResponses = append(marketPriceResponses, marketPriceResponse)
			}

			// Sort the slice based on the MarketCode field
			sort.Slice(marketPriceResponses, func(i, j int) bool {
				return marketPriceResponses[i].Code < marketPriceResponses[j].Code
			})

			return &marketPriceResponses, nil
		}
	case "us_index":
		{
			markets, err := mu.marketRepository.GetActiveMarketsIsUsIndex()
			if err != nil {
				return nil, errors.New("not found")
			}

			for _, market := range markets {
				marketCountryBase, countryBaseFound := getMarketCountry(marketCountries, market.BaseCurrency)
				marketCountryQuote, countryQuoteFound := getMarketCountry(marketCountries, market.QuoteCurrency)

				marketPriceResponse := MarketPriceResponse{
					Code:             market.Code,
					Price:            market.Price,
					Ask:              market.Ask,
					Bid:              market.Bid,
					Change:           market.Change,
					ChangePercentage: market.ChangePercentage,
					IsWatchList:      market.IsWatchList,
				}

				if countryBaseFound {
					marketPriceResponse.Description = marketCountryBase.Description
					marketPriceResponse.BaseCurrencyImage = marketCountryBase.FlagImage
				}

				if countryQuoteFound {
					marketPriceResponse.Description += " vs " + marketCountryQuote.Description
					marketPriceResponse.QuoteCurrencyImage = marketCountryQuote.FlagImage
				}

				marketPriceResponses = append(marketPriceResponses, marketPriceResponse)
			}

			// Sort the slice based on the MarketCode field
			sort.Slice(marketPriceResponses, func(i, j int) bool {
				return marketPriceResponses[i].Code < marketPriceResponses[j].Code
			})

			return &marketPriceResponses, nil
		}
	case "asia_index":
		{
			markets, err := mu.marketRepository.GetActiveMarketsIsAsiaIndex()
			if err != nil {
				return nil, errors.New("not found")
			}

			for _, market := range markets {
				marketCountryBase, countryBaseFound := getMarketCountry(marketCountries, market.BaseCurrency)
				marketCountryQuote, countryQuoteFound := getMarketCountry(marketCountries, market.QuoteCurrency)

				marketPriceResponse := MarketPriceResponse{
					Code:             market.Code,
					Price:            market.Price,
					Ask:              market.Ask,
					Bid:              market.Bid,
					Change:           market.Change,
					ChangePercentage: market.ChangePercentage,
					IsWatchList:      market.IsWatchList,
				}

				if countryBaseFound {
					marketPriceResponse.Description = marketCountryBase.Description
					marketPriceResponse.BaseCurrencyImage = marketCountryBase.FlagImage
				}

				if countryQuoteFound {
					marketPriceResponse.Description += " vs " + marketCountryQuote.Description
					marketPriceResponse.QuoteCurrencyImage = marketCountryQuote.FlagImage
				}

				marketPriceResponses = append(marketPriceResponses, marketPriceResponse)
			}

			// Sort the slice based on the MarketCode field
			sort.Slice(marketPriceResponses, func(i, j int) bool {
				return marketPriceResponses[i].Code < marketPriceResponses[j].Code
			})

			return &marketPriceResponses, nil
		}
	case "oil":
		{
			markets, err := mu.marketRepository.GetActiveMarketsIsOil()
			if err != nil {
				return nil, errors.New("not found")
			}

			for _, market := range markets {
				marketCountryBase, countryBaseFound := getMarketCountry(marketCountries, market.BaseCurrency)
				marketCountryQuote, countryQuoteFound := getMarketCountry(marketCountries, market.QuoteCurrency)

				marketPriceResponse := MarketPriceResponse{
					Code:             market.Code,
					Price:            market.Price,
					Ask:              market.Ask,
					Bid:              market.Bid,
					Change:           market.Change,
					ChangePercentage: market.ChangePercentage,
					IsWatchList:      market.IsWatchList,
				}

				if countryBaseFound {
					marketPriceResponse.Description = marketCountryBase.Description
					marketPriceResponse.BaseCurrencyImage = marketCountryBase.FlagImage
				}

				if countryQuoteFound {
					marketPriceResponse.Description += " vs " + marketCountryQuote.Description
					marketPriceResponse.QuoteCurrencyImage = marketCountryQuote.FlagImage
				}

				marketPriceResponses = append(marketPriceResponses, marketPriceResponse)
			}

			// Sort the slice based on the MarketCode field
			sort.Slice(marketPriceResponses, func(i, j int) bool {
				return marketPriceResponses[i].Code < marketPriceResponses[j].Code
			})

			return &marketPriceResponses, nil
		}
	case "all":
		{
			markets, err := mu.marketRepository.GetActiveMarkets()
			if err != nil {
				return nil, errors.New("not found")
			}

			for _, market := range markets {
				marketCountryBase, countryBaseFound := getMarketCountry(marketCountries, market.BaseCurrency)
				marketCountryQuote, countryQuoteFound := getMarketCountry(marketCountries, market.QuoteCurrency)

				marketPriceResponse := MarketPriceResponse{
					Code:             market.Code,
					Price:            market.Price,
					Ask:              common.RoundTo4DecimalPlaces(market.Ask),
					Bid:              common.RoundTo4DecimalPlaces(market.Bid),
					Change:           market.Change,
					ChangePercentage: market.ChangePercentage,
					IsWatchList:      market.IsWatchList,
				}

				if countryBaseFound {
					marketPriceResponse.Description = marketCountryBase.Description
					marketPriceResponse.BaseCurrencyImage = marketCountryBase.FlagImage
				}

				if countryQuoteFound {
					marketPriceResponse.Description += " vs " + marketCountryQuote.Description
					marketPriceResponse.QuoteCurrencyImage = marketCountryQuote.FlagImage
				}

				marketPriceResponses = append(marketPriceResponses, marketPriceResponse)
			}

			// Sort the slice based on the MarketCode field
			sort.Slice(marketPriceResponses, func(i, j int) bool {
				return marketPriceResponses[i].Code < marketPriceResponses[j].Code
			})

			return &marketPriceResponses, nil
		}

	default:
		return &marketPriceResponses, errors.New("invalid filterby")
	}
}

func (mu *marketUsecase) PriceUpdates() error {
	forexDatas, err := mu.marketApiclient.GetLatestMarketPrice()
	if err != nil {
		return err
	}

	for _, forexData := range forexDatas.Data {
		// forexData.Symbol = strings.TrimSuffix(forexData.Symbol, ".fl")
		market, _ := mu.marketRepository.GetMarketByCode(forexData.Symbol)
		// changedAmount := float64(common.GenerateRandomNumber(1, 10))
		// forexData.Change += changedAmount / 100
		// forexData.ChangePercentage += changedAmount / 10

		if market != nil {
			if err := mu.marketRepository.UpdateMarket(mapForexToMarket(forexData)); err != nil {
				return err
			}
			continue
		}

		if err := mu.marketRepository.CreateMarket(mapForexToMarket(forexData)); err != nil {
			return err
		}

	}

	return nil
}

func (mu *marketUsecase) LiveMarketUpdates() error {
	liveMarketUpdates, err := mu.marketApiclient.GetLiveMarketUpdates()
	if err != nil {
		return err
	}

	if len(liveMarketUpdates.Quotes) < 1 {
		return errors.New("not found")
	}

	for _, quote := range liveMarketUpdates.Quotes {
		//changedAmount := float64(common.GenerateRandomNumber(1, 10)) / 100

		market := &Market{
			Code: quote.BaseCurrency + quote.QuoteCurrency,
			Ask:  quote.Ask,
			Bid:  quote.Bid,
			Mid:  quote.Mid,
			// Ask:  quote.Ask + changedAmount,
			// Bid:  quote.Bid + changedAmount,
			// Mid:  quote.Mid + changedAmount,
		}

		if err := mu.marketRepository.UpdateMarket(market); err != nil {
			return err
		}
	}

	return nil
}

func (mu *marketUsecase) UpdateWatchlist(userId string, market Market) error {
	marketDb, err := mu.marketRepository.GetActiveMarketByCode(market.Code)
	if marketDb == nil || err != nil {
		return errors.New("market not found")
	}

	user, err := mu.userRepository.GetActiveUserByUserId(userId)
	if user == nil || err != nil {
		return errors.New("user not found")
	}

	isWatchlist := common.Contains(user.Watchlist, market.Code)

	if !isWatchlist && market.IsWatchList {
		user.Watchlist = append(user.Watchlist, market.Code)

		return mu.userRepository.UpdateUserWatchlist(user)
	}

	if isWatchlist && !market.IsWatchList {
		user.Watchlist = common.Remove(user.Watchlist, market.Code)

		return mu.userRepository.UpdateUserWatchlist(user)
	}

	return nil
}

func (mu *marketUsecase) GetWatchlist(userId, marketCode string) error {
	user, err := mu.userRepository.GetActiveUserByUserId(userId)
	if user == nil || err != nil {
		return errors.New("user not found")
	}

	isWatchlist := common.Contains(user.Watchlist, marketCode)
	if !isWatchlist {

		return errors.New("not found")
	}

	return nil
}

func (mu *marketUsecase) ConvertPrice(convertPrice ConvertPrice) (*ConvertPrice, error) {
	marketCurrencyRate, err := mu.marketRepository.GetMarketCurrencyRate(convertPrice.CurrencyBase, convertPrice.CurrencyQuote)
	if marketCurrencyRate == nil || err != nil {
		return &convertPrice, errors.New("market rate not found")
	}
	convertPrice.AmountQuote = common.RoundTo4DecimalPlaces(convertPrice.AmountBase * marketCurrencyRate.Rate)
	return &convertPrice, nil
}

func getMarketCountry(marketCountries []MarketCountry, currencyCode string) (MarketCountry, bool) {
	for _, marketCountry := range marketCountries {
		if marketCountry.CurrencyCode == currencyCode {
			return marketCountry, true
		}
	}

	return MarketCountry{}, false
}

// Function to map ForexPriceResponse to Market
func mapForexToMarket(forexPrice ArcMetaIntegratorPriceData) *Market {
	return &Market{
		Code:          forexPrice.Symbol,
		BaseCurrency:  safeSubstring(forexPrice.Symbol, 0, 3),
		QuoteCurrency: safeSubstring(forexPrice.Symbol, 3, 6),
		Price:         forexPrice.Last,
		Ask:           forexPrice.Ask,
		Bid:           forexPrice.Bid,
		// Change:           forexPrice.Change,
		// ChangePercentage: forexPrice.ChangePercentage,
		// DayLow:      forexPrice.DayLow,
		// DayHigh:     forexPrice.DayHigh,
		// YearLow:     forexPrice.YearLow,
		// YearHigh:    forexPrice.YearHigh,
		// Avg50Price:  forexPrice.Avg50Price,
		// Avg200Price: forexPrice.Avg200Price,
		BaseModel: base.BaseModel{IsActive: true},
	}

}

// Helper function to safely extract a substring
func safeSubstring(input string, start, end int) string {
	if len(input) < start {
		return "" // Return empty string if the start index is out of range
	}
	if len(input) < end {
		return input[start:] // Return substring from start to the end of the string
	}
	return input[start:end] // Return the full slice
}

// // Function to map ForexPriceResponse to Market
// func mapForexToMarket(forexPrice ForexPriceResponse) *Market {
// 	return &Market{
// 		Code:             forexPrice.Ticker,
// 		BaseCurrency:     forexPrice.BaseCurrency,
// 		QuoteCurrency:    forexPrice.QuoteCurrency,
// 		Price:            forexPrice.Price,
// 		Change:           forexPrice.Change,
// 		ChangePercentage: forexPrice.ChangePercentage,
// 		DayLow:           forexPrice.DayLow,
// 		DayHigh:          forexPrice.DayHigh,
// 		YearLow:          forexPrice.YearLow,
// 		YearHigh:         forexPrice.YearHigh,
// 		Avg50Price:       forexPrice.Avg50Price,
// 		Avg200Price:      forexPrice.Avg200Price,
// 		BaseModel:        base.BaseModel{IsActive: true},
// 	}

// }
