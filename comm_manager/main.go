package comm_manager

import "net/http"

type Connection interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
	Close() error
	GetUserId() string
	SetUserId(userId string)
}

type ConnectionManager interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Connection, error)
}

func GetConnectionManager() ConnectionManager {
	return NewWebSocketConnectionManager()
}
