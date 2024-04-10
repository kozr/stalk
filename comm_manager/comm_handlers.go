package comm_manager

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kozr/stalk/cache"
	db "github.com/kozr/stalk/database"
	"github.com/kozr/stalk/rsakey"
	follow_service "github.com/kozr/stalk/user_follow_service"
)

func HandleIncoming(ch chan string, conn Connection) {
	defer handleUserDisconnect(conn.GetUserId())

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break
		}
		fmt.Printf("Received message: %s", message)
		// TODO: Add RecencyLock to only execute the most recent request.
		go handleUserUrlChange(conn.GetUserId(), string(message), time.Now().Unix())
	}
}

func handleUserUrlChange(userId string, newUrl string, timestamp int64) {
	cache.UpdateUserHashedUrl(userId, newUrl)
	sameSiteFollowerIds := findFollowersOnSameUrl(userId, newUrl)
	broadcastToFollowers(sameSiteFollowerIds, userId)
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

func findFollowersOnSameUrl(userId string, message string) []string {
	followers, err := follow_service.GetOnlineFollowers(userId)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var sameSiteFollowerIds []string

	for _, follower := range followers {
		hashedUrl, err := cache.GetUserHashedUrl(follower)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !rsakey.CompareKeys(message, hashedUrl) {
			continue
		}
		sameSiteFollowerIds = append(sameSiteFollowerIds, follower)
	}

	return sameSiteFollowerIds
}

func broadcastToFollowers(followerIds []string, message string) {
	for _, followerId := range followerIds {
		go func() {
			channel, err := cache.GetUserChannel(followerId)
			if err != nil {
				fmt.Println("error getting channel for user id: ", followerId)
				return
			}
			select {
			case channel <- message:
			case <-time.After(1 * time.Second):
				fmt.Println("timeout sending message to user id: ", followerId)
			}
		}()
	}
}

func handleUserDisconnect(userId string) {
	db.UpdateUserOnlineStatus(userId, false)
	cache.RemoveUserChannel(userId)
	cache.RemoveUserHashedUrl(userId)
}
