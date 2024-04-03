package comm_manager

import "net/http"

type Connection interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
	ReadJSON(interface{}) error
	WriteJSON(interface{}) error
	Close() error
	SetCloseHandler(func())
	GetUserId() string
	SetUserId(userId string)
}

type ConnectionManager interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Connection, error)
}

func GetConnectionManager() ConnectionManager {
	return NewWebSocketConnectionManager()
}
