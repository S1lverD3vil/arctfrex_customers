package order

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"time"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (or *orderRepository) Create(order *Order) error {
	return or.db.Create(order).Error
}

func (or *orderRepository) GetActiveOrderByIdSessionId(orderId, sessionId string) (*Orders, error) {
	var orders Orders
	// log.Println("sebelum  query")
	if err := or.db.Table("orders").
		Joins("JOIN users ON users.id = orders.user_id").
		Joins("JOIN accounts ON accounts.id = orders.account_id").
		Select(`
			orders.id as orderid,
			orders.market_code as market_code,
			orders.lot,
			orders.price,
			orders.total_price,
			orders.take_profit_price,
			orders.stop_loss_price,
			accounts.total_pl,
			orders.account_id as accountid, 
			orders.user_id as userid,
			orders.type,
			orders.status,
			orders.order_time,
			orders.order_expiry_time
		`).
		Where(`
			orders.is_active = ? 
			AND users.session_id = ? 
			AND orders.id = ?`,
			true, sessionId, orderId,
		).
		Find(&orders).Error; err != nil {
		// Scan(&orders).Error; err != nil {

		return nil, err
	}
	// log.Println("setelah  query")
	// queryParams := Order{
	// 	ID: orderId,
	// 	BaseModel: base.BaseModel{
	// 		IsActive: true,
	// 	},
	// }

	// err := or.db.Find(&order, &queryParams).Error
	// if err != nil || order.MarketCode == common.STRING_EMPTY {
	// 	return nil, err
	// }

	return &orders, nil
}

func (or *orderRepository) GetActiveOrderByIdAccountIdUserId(orderId, accountId, userId string) (*Order, error) {
	var order Order
	queryParams := Order{
		ID:        orderId,
		AccountID: accountId,
		UserID:    userId,
		BaseModel: base.BaseModel{
			IsActive: true,
		},
	}

	err := or.db.Find(&order, &queryParams).Error
	if err != nil || order.MarketCode == common.STRING_EMPTY {
		return nil, err
	}

	return &order, nil
}

func (or *orderRepository) GetActiveOrderByAccountIdUserId(accountId, userId string) ([]Orders, error) {
	var orders []Orders
	if err := or.db.Table("orders").
		Joins("JOIN users ON users.id = orders.user_id").
		Joins("JOIN accounts ON accounts.id = orders.account_id").
		Select(`
			orders.id as orderid,
			orders.market_code as market_code,
			orders.lot,
			orders.price,
			orders.total_price,
			orders.take_profit_price,
			orders.stop_loss_price,
			accounts.total_pl,
			orders.account_id as accountid, 
			orders.user_id as userid,
			orders.type,
			orders.status,
			orders.order_time,
			orders.order_expiry_time
		`).
		Where(`
			orders.is_active = ? 
			AND users.session_id = ? 
			AND accounts.id = ?`,
			true, userId, accountId,
		).
		Scan(&orders).Error; err != nil {

		return nil, err
	}

	return orders, nil
}

func (or *orderRepository) Update(order *Order) error {
	return or.db.Model(order).Updates(order).Error
}

func (or *orderRepository) UpdateOrderCloseAllByTypeStatus(orderCloseAllRequest *OrderCloseAll) error {
	return or.db.Model(&Order{}).
		Where(`account_id = ? 
				AND user_id = ? 
				AND status = ?`,
			orderCloseAllRequest.AccountID,
			orderCloseAllRequest.UserID,
			orderCloseAllRequest.Status).
		Updates(Order{Status: enums.OrderStatusClosed}).Error
}

func (or *orderRepository) UpdateOrderCloseAllExpired() error {
	return or.db.Model(&Order{}).
		Where(`DATE(order_expiry_time) < ? 
				AND status = ?`,
			time.Now().Format("2006-01-02"),
			enums.OrderStatusPending).
		Updates(Order{Status: enums.OrderStatusClosed}).Error
}
