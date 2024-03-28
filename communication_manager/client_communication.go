package connection_manager

import (
	"fmt"

	"github.com/gorilla/websocket"
	follow_service "github.com/kozr/stalk/follow_service"
	"github.com/kozr/stalk/rsakey"
)

func HandleIncoming(ch chan string, conn Connection) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break
		}
		if rsakey.CompareKeys(string(message), "example.com") {
			broadcastToFollowers(conn.GetUserId(), string(message))
		}
		fmt.Printf("Received message: %s", message)
	}
}

func HandleOutgoing(ch chan string, conn Connection) {
	for {
		message := <-ch
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			conn.Close()
			break
		}
		fmt.Printf("Sent message: %s", message)
	}
}

func broadcastToFollowers(userId string, message string) {
	followers, err := follow_service.GetAliveFollowerChannels(userId)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, follower := range followers {
		follower <- message
	}
}
