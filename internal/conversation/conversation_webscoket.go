package conversation

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	sessions map[string][]*websocket.Conn // Map of sessionID to connections
	mu       sync.Mutex                   // To synchronize access to the sessions map
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		sessions: make(map[string][]*websocket.Conn),
	}
}

func (cm *ConnectionManager) AddConnection(sessionID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.sessions[sessionID] = append(cm.sessions[sessionID], conn)
}

func (cm *ConnectionManager) RemoveConnection(sessionID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	connections := cm.sessions[sessionID]
	for i, c := range connections {
		if c == conn {
			cm.sessions[sessionID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}

func (cm *ConnectionManager) Broadcast(sessionID string, message interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	connections := cm.sessions[sessionID]
	for _, conn := range connections {
		if err := conn.WriteJSON(message); err != nil {
			log.Println("Broadcast error:", err)
		}
	}
}
