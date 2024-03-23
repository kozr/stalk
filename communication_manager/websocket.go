package connection_manager

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	conn *websocket.Conn
}

func (wsc *WebSocketConnection) ReadMessage() (int, []byte, error) {
	return wsc.conn.ReadMessage()
}

func (wsc *WebSocketConnection) WriteMessage(messageType int, data []byte) error {
	return wsc.conn.WriteMessage(messageType, data)
}

func (wsc *WebSocketConnection) Close() error {
	return wsc.conn.Close()
}

type WebSocketConnectionManager struct {
	upgrader websocket.Upgrader
}

func NewWebSocketConnectionManager() *WebSocketConnectionManager {
	return &WebSocketConnectionManager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (cm *WebSocketConnectionManager) Upgrade(w http.ResponseWriter, r *http.Request) (Connection, error) {
	conn, err := cm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &WebSocketConnection{conn}, nil
}
