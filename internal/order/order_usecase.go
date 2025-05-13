package order

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/model"
	"arctfrex-customers/internal/repository"
)

type OrderUsecase interface {
	Orders(webSocketRequest WebSocketRequest) (*[]Orders, model.Account, error)
	Submit(order *Order) error
	UpdateByOrderId(order *Order) error
	CloseAllOrderByTypeStatus(orderCloseAllRequest *OrderCloseAll) error
	CloseAllExpiredOrder() error
}

type orderUsecase struct {
	orderRepository   OrderRepository
	accountRepository repository.AccountRepository
	marketRepository  repository.MarketRepository
}

func NewOrderUsecase(
	or OrderRepository,
	ar repository.AccountRepository,
	mr repository.MarketRepository,
) *orderUsecase {
	return &orderUsecase{
		orderRepository:   or,
		accountRepository: ar,
		marketRepository:  mr,
	}
}

func (ou *orderUsecase) Orders(webSocketRequest WebSocketRequest) (*[]Orders, model.Account, error) {
	var orders []Orders
	marketCountries, err := ou.marketRepository.GetActiveMarketCountries()
	if err != nil {
		return nil, model.Account{}, errors.New("not found")
	}
	account, err := ou.accountRepository.GetAccountsById(webSocketRequest.AccountId)
	if err != nil || account == nil || account.ID == common.STRING_EMPTY {
		return nil, *account, errors.New("record not found")
	}
	// changedAmount := float64(common.GenerateRandomNumber(1, 1000))
	// account.Equity += changedAmount
	//log.Println(webSocketRequest)
	switch strings.ToLower(webSocketRequest.FilterBy) {
	case "code":
		{
			order, err := ou.orderRepository.GetActiveOrderByIdSessionId(webSocketRequest.OrderId, webSocketRequest.SessionId)
			// log.Println("cek setelah get data")
			if err != nil {
				return nil, *account, errors.New("not found")
			}

			if order == nil {
				return &orders, *account, nil
			}

			market, err := ou.marketRepository.GetActiveMarketByCode(order.MarketCode)
			if err != nil {
				return nil, *account, errors.New("not found")
			}

			order.MarketAsk = common.RoundTo4DecimalPlaces(market.Ask)
			order.MarketBid = common.RoundTo4DecimalPlaces(market.Bid)
			marketCountryBase, countryBaseFound := getMarketCountry(marketCountries, market.BaseCurrency)
			marketCountryQuote, countryQuoteFound := getMarketCountry(marketCountries, market.QuoteCurrency)

			if countryBaseFound {
				order.MarketDescription = marketCountryBase.Description
				order.MarketBaseCurrencyImage = marketCountryBase.FlagImage
			}

			if countryQuoteFound {
				order.MarketDescription += " vs " + marketCountryQuote.Description
				order.MarketQuoteCurrencyImage = marketCountryQuote.FlagImage
			}

			orders = append(orders, *order)

			// // Sort the slice based on the MarketCode field
			// sort.Slice(orders, func(i, j int) bool {
			// 	return orders[i].OrderTime.After(orders[j].OrderTime)
			// })
			// log.Println("cek final")
			return &orders, *account, nil
		}
	case "all":
		{
			//log.Println(webSocketRequest.AccountId)
			orders, err := ou.orderRepository.GetActiveOrderByAccountIdUserId(webSocketRequest.AccountId, webSocketRequest.SessionId)
			//log.Println(orders)
			if err != nil {
				return nil, *account, errors.New("not found")
			}

			for i := range orders {
				order := &orders[i]
				// log.Println(order.MarketCode)
				market, err := ou.marketRepository.GetActiveMarketByCode(order.MarketCode)
				if err != nil {
					return nil, *account, errors.New("not found")
				}
				// log.Println(market)
				order.MarketAsk = common.RoundTo4DecimalPlaces(market.Ask)
				order.MarketBid = common.RoundTo4DecimalPlaces(market.Bid)
				marketCountryBase, countryBaseFound := getMarketCountry(marketCountries, market.BaseCurrency)
				marketCountryQuote, countryQuoteFound := getMarketCountry(marketCountries, market.QuoteCurrency)

				if countryBaseFound {
					order.MarketDescription = marketCountryBase.Description
					order.MarketBaseCurrencyImage = marketCountryBase.FlagImage
				}

				if countryQuoteFound {
					order.MarketDescription += " vs " + marketCountryQuote.Description
					order.MarketQuoteCurrencyImage = marketCountryQuote.FlagImage
				}
			}

			// Sort the slice based on the MarketCode field
			sort.Slice(orders, func(i, j int) bool {
				return orders[i].OrderTime.After(orders[j].OrderTime)
			})

			return &orders, *account, nil
		}

	default:
		return &orders, *account, errors.New("invalid filterby")
	}
}

func (ou *orderUsecase) Submit(order *Order) error {
	account, err := ou.accountRepository.GetAccountsByIdUserId(order.UserID, order.AccountID)
	if err != nil || account == nil || account.ID == common.STRING_EMPTY {
		return errors.New("record not found")
	}
	order.TotalPrice = order.Lot * order.Price
	if order.Type == enums.OrderTypeBuy && account.Balance < order.TotalPrice {
		return errors.New("insufficient balance")
	}

	orderID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	order.ID = common.UUIDNormalizer(orderID)
	//order.Status = enums.OrderStatusNew
	order.OrderTime = time.Now()
	order.OrderExpiryTime = order.OrderTime.AddDate(0, 0, int(order.ExpiryDays))
	order.IsActive = true

	return ou.orderRepository.Create(order)
}

func (ou *orderUsecase) UpdateByOrderId(order *Order) error {
	account, err := ou.accountRepository.GetAccountsByIdUserId(order.UserID, order.AccountID)
	if err != nil || account == nil || account.ID == common.STRING_EMPTY {
		return errors.New("record not found")
	}

	orderDb, err := ou.orderRepository.GetActiveOrderByIdAccountIdUserId(order.ID, order.AccountID, order.UserID)
	if err != nil || orderDb == nil {
		return errors.New("record not found")
	}

	orderDb.Price = order.Price
	orderDb.Status = order.Status
	orderDb.TotalPrice = orderDb.Lot * orderDb.Price
	orderDb.StopLossPrice = order.StopLossPrice
	orderDb.TakeProfitPrice = order.TakeProfitPrice

	if orderDb.Type == enums.OrderTypeBuy && account.Balance < orderDb.TotalPrice {
		return errors.New("insufficient balance")
	}

	return ou.orderRepository.Update(orderDb)
}

func (ou *orderUsecase) CloseAllOrderByTypeStatus(orderCloseAllRequest *OrderCloseAll) error {
	account, err := ou.accountRepository.GetAccountsByIdUserId(orderCloseAllRequest.UserID, orderCloseAllRequest.AccountID)
	if err != nil || account == nil || account.ID == common.STRING_EMPTY {
		return errors.New("record not found")
	}

	return ou.orderRepository.UpdateOrderCloseAllByTypeStatus(orderCloseAllRequest)
}
func (ou *orderUsecase) CloseAllExpiredOrder() error {
	return ou.orderRepository.UpdateOrderCloseAllExpired()
}

func getMarketCountry(marketCountries []model.MarketCountry, currencyCode string) (model.MarketCountry, bool) {
	for _, marketCountry := range marketCountries {
		if marketCountry.CurrencyCode == currencyCode {
			return marketCountry, true
		}
	}

	return model.MarketCountry{}, false
}
