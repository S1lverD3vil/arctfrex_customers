package market

import (
	websocketCache "arctfrex-customers/internal/cache"
	"arctfrex-customers/internal/common"
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

type marketHandler struct {
	jwtMiddleware *middleware.JWTMiddleware
	marketUsecase MarketUsecase
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin (set restrictions for security)
	},
}

func NewMarketHandler(
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	mu MarketUsecase,
) *marketHandler {
	handler := &marketHandler{
		jwtMiddleware: jmw,
		marketUsecase: mu,
	}

	unprotectedRoutes, protectedRoutes := engine.Group("/markets"), engine.Group("/markets")
	//unprotectedRoutes.GET("/price", handler.Price)
	unprotectedRoutes.GET("/price/ws", handler.PriceWebSocket)
	unprotectedRoutes.POST("convert/price", handler.ConvertPrice)
	protectedRoutes.Use(jmw.ValidateToken())
	{
		protectedRoutes.GET("/price", handler.Price)
		protectedRoutes.PATCH("/watchlist", handler.UpdateWatchList)
		protectedRoutes.GET("/watchlist/:marketCode", handler.GetWatchList)
	}

	return handler
}

func (mh *marketHandler) PriceWebSocket(c *gin.Context) {
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
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	// Channel to receive messages from the client
	messageChan := make(chan []byte)

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
				filterBy := WebSocketRequest{
					FilterBy:   "all",
					MarketCode: common.STRING_EMPTY,
				}
				if cachedData, found := websocketCache.ClientCache.Get(clientID + "_filter"); found {
					// Unmarshal JSON message into the struct
					if err := json.Unmarshal([]byte(cachedData.(string)), &filterBy); err != nil {
						log.Println("Error unmarshaling message:", err)
						continue
					}
				}

				markets, err := mh.marketUsecase.Price(filterBy)
				if err != nil {
					log.Printf("Error handling price update: %v", err)
					conn.WriteJSON(WebSocketResponse{Message: "Failed to process message from usecase"})
					continue
				}

				responseMessage := WebSocketResponse{
					Message: "Received your message",
					Time:    t.Format("2006-01-02 15:04:05"),
				}

				if strings.ToLower(filterBy.FilterBy) == "code" && markets != nil && len(*markets) > 0 {
					responseMessage.Data = (*markets)[0]
				} else {
					responseMessage.Data = markets
				}

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

				// Cache the new price for this specific client with no expiration
				websocketCache.ClientCache.Set(clientID+"_filter", string(jsonData), cache.NoExpiration)

				markets, err := mh.marketUsecase.Price(webSocketRequest)
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

				if strings.ToLower(webSocketRequest.FilterBy) == "code" && markets != nil && len(*markets) > 0 {
					responseMessage.Data = (*markets)[0]
				} else {
					responseMessage.Data = markets
				}

				if err := conn.WriteJSON(responseMessage); err != nil {
					log.Println("Error sending message:", err)
					return
				}
			}
		}
	}
}

func (mh *marketHandler) Price(c *gin.Context) {
	markets, err := mh.marketUsecase.Price(WebSocketRequest{FilterBy: "all"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, markets)
}

func (mh *marketHandler) UpdateWatchList(c *gin.Context) {
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

	var market *Market
	if err := c.ShouldBindJSON(&market); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := mh.marketUsecase.UpdateWatchlist(userId, *market); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (mh *marketHandler) GetWatchList(c *gin.Context) {
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

	if err := mh.marketUsecase.GetWatchlist(userId, c.Param("marketCode")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (mh *marketHandler) ConvertPrice(c *gin.Context) {
	var convertPrice *ConvertPrice
	if err := c.ShouldBindJSON(&convertPrice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	convertPrice, err := mh.marketUsecase.ConvertPrice(*convertPrice)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, convertPrice)
}
