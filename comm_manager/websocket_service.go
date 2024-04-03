package comm_manager

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	conn   *websocket.Conn
	userId string
}

func (wsc *WebSocketConnection) SetUserId(userId string) {
	wsc.userId = userId
}

func (wsc *WebSocketConnection) GetUserId() string {
	return wsc.userId
}

func (wsc *WebSocketConnection) ReadMessage() (int, []byte, error) {
	return wsc.conn.ReadMessage()
}

func (wsc *WebSocketConnection) ReadJSON(v interface{}) error {
	return wsc.conn.ReadJSON(v)
}

func (wsc *WebSocketConnection) WriteMessage(messageType int, data []byte) error {
	return wsc.conn.WriteMessage(messageType, data)
}

func (wsc *WebSocketConnection) SetCloseHandler(callback func()) {
	wsc.conn.SetCloseHandler(func(code int, text string) error {
		callback()
		return nil
	})
}

func (wsc *WebSocketConnection) WriteJSON(v interface{}) error {
	return wsc.conn.WriteJSON(v)
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

func (cm *WebSocketConnectionManager) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Connection, error) {
	conn, err := cm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	// Here, we return a *WebSocketConnection as a Connection interface type.
	// This is a valid operation since *WebSocketConnection implements all methods of Connection.
	return &WebSocketConnection{conn: conn}, nil
}
