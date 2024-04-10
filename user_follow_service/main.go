package follow_service

import (
	cache "github.com/kozr/stalk/cache"
	db "github.com/kozr/stalk/database"
)

func Follow(followeeId string, followerId string) error {
	err := db.AddFollower(followeeId, followerId)
	if err != nil {
		return err
	}
	return nil
}

func Unfollow(followeeId string, followerId string) error {
	err := db.RemoveFollower(followeeId, followerId)
	if err != nil {
		return err
	}
	return nil
}

func AddChannel(userId string, channel chan string) {
	cache.UpdateUserChannel(userId, channel)
}

func RemoveChannel(userId string) {
	cache.RemoveUserChannel(userId)
}

func GetOnlineFollowers(userId string) ([]string, error) {
	onlineFollowers, err := db.GetOnlineFollowers(userId)
	return onlineFollowers, err
}
