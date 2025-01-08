package websocket

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	Conn *websocket.Conn
}

// NewWebSocketClient creates a new WebSocket client and connects to the given URL
func NewWebSocketClient(wsURL string) (*WebSocketClient, error) {
	u, err := url.Parse(wsURL)
	if err != nil {
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return &WebSocketClient{Conn: conn}, nil
}

// Send sends a message to the WebSocket server
func (client *WebSocketClient) Send(message []byte) error {
	return client.Conn.WriteMessage(websocket.TextMessage, message)
}

// Read listens for messages from the WebSocket server
func (client *WebSocketClient) Read() ([]byte, error) {
	_, message, err := client.Conn.ReadMessage()
	return message, err
}

// Close closes the WebSocket connection
func (client *WebSocketClient) Close() error {
	return client.Conn.Close()
}

// Reconnect tries to reconnect in case of failure
func (client *WebSocketClient) Reconnect(wsURL string, reconnectInterval time.Duration) error {
	var err error
	for {
		client.Conn, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
		if err == nil {
			log.Println("Reconnected to WebSocket server")
			return nil
		}
		log.Printf("Failed to reconnect: %v, retrying in %v", err, reconnectInterval)
		time.Sleep(reconnectInterval)
	}
}

// // WebSocketClient represents a WebSocket client
// type WebSocketClient struct {
// 	URL  string
// 	Conn *websocket.Conn
// }

// // NewWebSocketClient creates a new WebSocket client
// func NewWebSocketClient(url string) *WebSocketClient {
// 	return &WebSocketClient{URL: url}
// }

// // Connect establishes the WebSocket connection
// func (client *WebSocketClient) Connect() error {
// 	var err error
// 	client.Conn, _, err = websocket.DefaultDialer.Dial(client.URL, nil)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Close closes the WebSocket connection
// func (client *WebSocketClient) Close() error {
// 	if client.Conn != nil {
// 		return client.Conn.Close()
// 	}
// 	return nil
// }

// // ReadMessages reads messages from the WebSocket connection
// func (client *WebSocketClient) ReadMessages(onMessage func(message []byte)) {
// 	for {
// 		_, msg, err := client.Conn.ReadMessage()
// 		if err != nil {
// 			log.Printf("Error reading message: %v", err)
// 			break
// 		}
// 		onMessage(msg)
// 	}
// }
