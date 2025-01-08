package order

import (
	websocketCache "arctfrex-customers/internal/cache"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/middleware"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
)

type orderHandler struct {
	jwtMiddleware *middleware.JWTMiddleware
	orderUsecase  OrderUsecase
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin (set restrictions for security)
	},
}

func NewOrderHandler(engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	orderUsecase OrderUsecase) *orderHandler {
	handler := &orderHandler{
		jwtMiddleware: jmw,
		orderUsecase:  orderUsecase,
	}

	unprotectedRoutes, protectedRoutes := engine.Group("/order"), engine.Group("/order")
	unprotectedRoutes.GET("/ws", handler.OrderWebSocket)
	protectedRoutes.Use(jmw.ValidateToken())
	{
		protectedRoutes.POST("/submit", handler.Submit)
		protectedRoutes.PATCH("/:orderId", handler.UpdateByOrderId)
		protectedRoutes.POST("/close/all", handler.CloseAllOrder)
		// protectedRoutes.GET("/pending/:accountId", handler.Pending)
		// protectedRoutes.GET("/:accountId", handler.DepositByAccountId)
		// protectedRoutes.GET("/detail/:depositId", handler.Detail)
	}

	return handler
}

func (oh *orderHandler) OrderWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) // Upgrade HTTP to WebSocket
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		conn.WriteJSON(WebSocketResponse{Message: "Failed to upgrade to WebSocket"})
		return
	}
	defer conn.Close()

	// Create a context to manage the goroutine lifecycle
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Generate a unique client ID (example using timestamp)
	clientID := time.Now().String()

	// This simulates sending automatic messages to the client
	ticker := time.NewTicker(2 * time.Hour)
	defer ticker.Stop()

	// Channel to receive messages from the client
	messageChan := make(chan []byte)
	//log.Println(messageChan)

	// Goroutine to read messages from the WebSocket
	go func(ctx context.Context) {
		defer close(messageChan)
		for {
			select {
			case <-ctx.Done():
				// Exit the goroutine if context is canceled
				log.Println("Goroutine stopped.")
				return
			default:
				messageType, msg, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
						log.Printf("WebSocket closed with status %d: %v", websocket.CloseNormalClosure, err)
					} else if websocket.IsUnexpectedCloseError(err) {
						log.Printf("Unexpected WebSocket closure: %v", err)
					} else {
						log.Println("Error while reading message:", err)
					}
					return
				}
				if messageType == websocket.TextMessage {
					messageChan <- msg
				}
			}
		}
	}(ctx)

	// Main loop to handle sending and optionally receiving messages
	for {
		select {
		case t := <-ticker.C:
			{
				//log.Println("masuk ticker")
				filterBy := WebSocketRequest{
					FilterBy: "all",
					OrderId:  common.STRING_EMPTY,
				}
				if cachedData, found := websocketCache.ClientCache.Get(clientID + "_filter"); found {
					// log.Println("cache data")
					// log.Println(cachedData)
					// Unmarshal JSON message into the struct
					if err := json.Unmarshal([]byte(cachedData.(string)), &filterBy); err != nil {
						log.Println("Error unmarshaling message:", err)
						continue
					}
				}

				//markets, err := mh.marketUsecase.Price(filterBy)
				orders, account, err := oh.orderUsecase.Orders(filterBy)
				//log.Println(filterBy)
				if err != nil {
					log.Printf("Error handling price update: %v", err)
					conn.WriteJSON(WebSocketResponse{Message: "Failed to process message from usecase"})
					continue
				}

				responseMessage := WebSocketResponse{
					Message: "Received your message",
					Time:    t.Format("2006-01-02 15:04:05"),
				}

				if strings.ToLower(filterBy.FilterBy) == "code" && orders != nil && len(*orders) > 0 {
					responseMessage.Data = (*orders)[0]
				} else {
					var openOrders []Orders
					var pendingOrders []Orders
					var historyOrders []Orders
					// for _, order := range *orders {
					for i := range *orders {
						order := (*orders)[i]
						switch order.Status {
						case enums.OrderStatusNew:
							{
								// changedAmount := float64(common.GenerateRandomNumber(1, 10)) / 100
								// order.Price += changedAmount
								order.Price = common.RoundTo4DecimalPlaces(order.Price)
								openOrders = append(openOrders, order)
							}
						case enums.OrderStatusPending:
							{
								// changedAmount := float64(common.GenerateRandomNumber(1, 10)) / 100
								// order.Price += changedAmount
								order.Price = common.RoundTo4DecimalPlaces(order.Price)
								pendingOrders = append(pendingOrders, order)
							}
						case enums.OrderStatusClosed:
							{
								historyOrders = append(historyOrders, order)
							}
						}
					}

					responseMessage.Data = &OrderData{
						AccountData:   account,
						OpenOrders:    openOrders,
						PendingOrders: pendingOrders,
						HistoryOrders: historyOrders,
					}
					// responseMessage.Data = orders
				}
				// if strings.ToLower(filterBy.FilterBy) == "code" && markets != nil && len(*markets) > 0 {
				// 	responseMessage.Data = (*markets)[0]
				// } else {
				// 	responseMessage.Data = markets
				// }
				// responseMessage.Data = orders

				if err := conn.WriteJSON(responseMessage); err != nil {
					log.Println("Error sending message:", err)
					return
				}

			}
		case msg, ok := <-messageChan:
			{
				if !ok {
					// Channel closed, exit loop
					return
				}

				log.Printf("%sReceived message: %s%s", common.LOG_CONSOLE_GREEN, msg, common.LOG_CONSOLE_RESET)

				// Unmarshal JSON message into the struct
				var webSocketRequest WebSocketRequest
				if err := json.Unmarshal(msg, &webSocketRequest); err != nil {
					log.Println("Error unmarshaling message:", err)
					continue
				}

				// Convert the struct to JSON
				jsonData, err := json.Marshal(webSocketRequest)
				if err != nil {
					fmt.Println("Error marshalling struct to JSON:", err)
					return
				}
				//log.Println("masuk message")
				// Cache the new price for this specific client with no expiration
				websocketCache.ClientCache.Set(clientID+"_filter", string(jsonData), cache.NoExpiration)

				//markets, err := mh.marketUsecase.Price(webSocketRequest)
				orders, account, err := oh.orderUsecase.Orders(webSocketRequest)
				if err != nil {
					log.Printf("Error handling price update: %v", err)
					conn.WriteJSON(WebSocketResponse{Message: "Failed to process message from usecase"})
					continue
				}

				responseMessage := WebSocketResponse{
					Message: "Received your message",
					// Data:    markets,
					Time: time.Now().Format("2006-01-02 15:04:05"),
				}

				if strings.ToLower(webSocketRequest.FilterBy) == "code" && orders != nil && len(*orders) > 0 {
					responseMessage.Data = (*orders)[0]
				} else {
					var openOrders []Orders
					var pendingOrders []Orders
					var historyOrders []Orders
					for _, order := range *orders {
						switch order.Status {
						case enums.OrderStatusNew:
							{
								openOrders = append(openOrders, order)
							}
						case enums.OrderStatusPending:
							{
								pendingOrders = append(pendingOrders, order)
							}
						case enums.OrderStatusClosed:
							{
								historyOrders = append(historyOrders, order)
							}
						}
					}

					responseMessage.Data = &OrderData{
						AccountData:   account,
						OpenOrders:    openOrders,
						PendingOrders: pendingOrders,
						HistoryOrders: historyOrders,
					}
				}
				//responseMessage.Data = orders

				if err := conn.WriteJSON(responseMessage); err != nil {
					log.Println("Error sending message:", err)
					return
				}
			}
		}
	}
}

func (oh *orderHandler) Submit(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
	}
	var order *Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.UserID = userId

	err := oh.orderUsecase.Submit(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//log.Println(order.ID)

	c.JSON(http.StatusCreated, order)
}
func (oh *orderHandler) UpdateByOrderId(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
	}
	var order *Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.UserID = userId
	order.ID = c.Param("orderId")

	if err := oh.orderUsecase.UpdateByOrderId(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
func (oh *orderHandler) CloseAllOrder(c *gin.Context) {
	// Retrieve the userID from context
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("userID not found in context")
	}

	// Convert userID to string
	userId, ok := userID.(string)
	if !ok {
		log.Println("userID is not of type string")
	}
	var orderCloseAll *OrderCloseAll
	if err := c.ShouldBindJSON(&orderCloseAll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderCloseAll.UserID = userId

	if err := oh.orderUsecase.CloseAllOrderByTypeStatus(orderCloseAll); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
