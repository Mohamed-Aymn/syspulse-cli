package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// WebSocketClient handles WebSocket connection and message sending.
type WebSocketClient struct {
	Conn *websocket.Conn
}

// NewWebSocketClient creates and returns a WebSocket connection.
func NewWebSocketClient(url string) (*WebSocketClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("WebSocket connection failed: %v", err)
	}
	return &WebSocketClient{Conn: conn}, nil
}

// SendMessage sends a JSON message over the WebSocket.
func (ws *WebSocketClient) SendMessage(data map[string]string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Error encoding JSON: %v", err)
	}

	if err := ws.Conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		return fmt.Errorf("Error sending message: %v", err)
	}
	fmt.Printf("Sent: %s\n", jsonData)
	return nil
}

// Close closes the WebSocket connection.
func (ws *WebSocketClient) Close() {
	ws.Conn.Close()
}
