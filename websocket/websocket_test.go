package websocket

import (
	"testing"

	// "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNewWebSocketClient(t *testing.T) {
	_, err := NewWebSocketClient("ws://invalid-url")
	assert.NotNil(t, err, "Expected error for invalid WebSocket URL")
}
