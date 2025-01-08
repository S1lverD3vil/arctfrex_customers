package order

import (
	"arctfrex-customers/internal/account"
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common/enums"
	"time"
)

type Order struct {
	ID              string            `json:"orderid" gorm:"primary_key"`
	MarketCode      string            `json:"market_code"`
	Lot             float64           `json:"lot"`
	Price           float64           `json:"price"`
	TotalPrice      float64           `json:"total_price"`
	TakeProfitPrice float64           `json:"take_profit_price"`
	StopLossPrice   float64           `json:"stop_loss_price"`
	AccountID       string            `json:"accountid"`
	UserID          string            `json:"userid"`
	Type            enums.OrderType   `json:"type"`
	Status          enums.OrderStatus `json:"status"`
	OrderTime       time.Time         `json:"order_time"`
	ExpiryDays      int               `json:"expiry"`
	OrderExpiryTime time.Time         `json:"order_expiry_time"`

	base.BaseModel
}
type Orders struct {
	Orderid                  string            `json:"orderid"`
	MarketCode               string            `json:"market_code"`
	MarketAsk                float64           `json:"market_ask"`
	MarketBid                float64           `json:"market_bid"`
	MarketDescription        string            `json:"market_description"`
	MarketBaseCurrencyImage  string            `json:"market_base_currency_image"`
	MarketQuoteCurrencyImage string            `json:"market_quote_currency_image"`
	Lot                      float64           `json:"lot"`
	Price                    float64           `json:"price"`
	TotalPrice               float64           `json:"total_price"`
	TakeProfitPrice          float64           `json:"take_profit_price"`
	StopLossPrice            float64           `json:"stop_loss_price"`
	TotalPL                  float64           `json:"total_pl"`
	TotalSLTP                float64           `json:"total_sl_tp"`
	Comission                float64           `json:"comission"`
	Tax                      float64           `json:"tax"`
	Accountid                string            `json:"accountid"`
	Userid                   string            `json:"userid"`
	Type                     enums.OrderType   `json:"type"`
	Status                   enums.OrderStatus `json:"status"`
	OrderTime                time.Time         `json:"order_time"`
	ExpiryDays               int               `json:"expiry"`
	OrderExpiryTime          time.Time         `json:"order_expiry_time"`
}
type WebSocketRequest struct {
	FilterBy  string `json:"filter_by"`
	OrderId   string `json:"order_id"`
	SessionId string `json:"session_id"`
	AccountId string `json:"account_id"`
}
type WebSocketResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    string `json:"time"`
}
type OrderData struct {
	AccountData   account.Account `json:"account_data"`
	OpenOrders    []Orders        `json:"open_orders"`
	PendingOrders []Orders        `json:"pending_orders"`
	HistoryOrders []Orders        `json:"history_orders"`
}
type OrderCloseAll struct {
	AccountID string            `json:"accountid"`
	UserID    string            `json:"userid"`
	Type      enums.OrderType   `json:"type"`
	Status    enums.OrderStatus `json:"status"`
}

type OrderRepository interface {
	Create(order *Order) error
	GetActiveOrderByIdSessionId(orderId, sessionId string) (*Orders, error)
	GetActiveOrderByIdAccountIdUserId(orderId, accountId, userId string) (*Order, error)
	GetActiveOrderByAccountIdUserId(accountId, userId string) ([]Orders, error)
	Update(order *Order) error
	UpdateOrderCloseAllByTypeStatus(orderCloseAllRequest *OrderCloseAll) error
	UpdateOrderCloseAllExpired() error
}
