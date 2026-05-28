package pkg

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	request_context "main/internal/context"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ConnectionManager struct {
	connection map[uuid.UUID]*websocket.Conn
	mu         sync.RWMutex
}

var Manager = &ConnectionManager{
	connection: make(map[uuid.UUID]*websocket.Conn),
}

func (c *ConnectionManager) HandleConnection(w http.ResponseWriter, r *http.Request) {
	userID, ok := request_context.UserID(r.Context())
	if !ok || userID == uuid.Nil {
		http.Error(w, "websocket: invalid user context", http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	defer c.disconnect(userID)
	Manager.mu.Lock()
	Manager.connection[userID] = conn
	Manager.mu.Unlock()
	c.readMessages(userID, conn)
}

func (c *ConnectionManager) readMessages(userID uuid.UUID, conn *websocket.Conn) {
	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message := string(messageBytes)
		log.Println("Message: ", message)
		c.writeMessage(userID, message)
	}
}

func (c *ConnectionManager) writeMessage(userID uuid.UUID, message string) error {
	Manager.mu.RLock()
	conn, ok := Manager.connection[userID]
	if !ok {
		return fmt.Errorf("user %s not connected", userID)
	}
	Manager.mu.RUnlock()
	Manager.mu.Lock()
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return fmt.Errorf("ws: write error %w", err)
	}
	defer Manager.mu.Unlock()
	return nil
}

func (c *ConnectionManager) disconnect(userID uuid.UUID) {
	Manager.mu.Lock()
	defer Manager.mu.Unlock()
	if conn, ok := Manager.connection[userID]; ok {
		conn.Close()
		delete(Manager.connection, userID)
	}
}
