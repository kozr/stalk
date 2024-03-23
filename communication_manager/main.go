package connection_manager

import "net/http"

type Connection interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
	Close() error
}

type ConnectionManager interface {
	Upgrade(http.ResponseWriter, *http.Request) (Connection, error)
}

func GetConnectionManager() ConnectionManager {
	return NewWebSocketConnectionManager()
}
