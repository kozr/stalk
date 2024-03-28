package follow_service

// array of followers and followees
var followMap = make(map[string][]string)
var aliveChannels = make(map[string]chan string)

// TODO: Integrate redis for followMap and aliveChannels
func Follow(followerId string, followeeId string) error {
	//  followMap: check if followerId is already following followeeId
	for _, followee := range followMap[followerId] {
		if followee == followeeId {
			return nil
		}
	}
	// add followeeId to followerId's list of followees
	followMap[followerId] = append(followMap[followerId], followeeId)
	// *** add followerId to followeeId's list of followers
	followMap[followeeId] = append(followMap[followeeId], followerId)

	return nil
}

// TODO: Integrate redis for followMap and aliveChannels
func Unfollow(followerId string, followeeId string) error {
	// remove followeeId from followerId's list of followees
	for i, followee := range followMap[followerId] {
		if followee == followeeId {
			followMap[followerId] = append(followMap[followerId][:i], followMap[followerId][i+1:]...)
			break
		}
	}

	// *** remove followerId from followeeId's list of followers
	for i, follower := range followMap[followeeId] {
		if follower == followerId {
			followMap[followeeId] = append(followMap[followeeId][:i], followMap[followeeId][i+1:]...)
			break
		}
	}

	return nil
}

func AddChannel(userId string, channel chan string) {
	aliveChannels[userId] = channel
}

func RemoveChannel(userId string) {
	delete(aliveChannels, userId)
}

func GetAliveFollowerChannels(userId string) ([]chan string, error) {
	var channels []chan string
	for _, follower := range followMap[userId] {
		if channel, ok := aliveChannels[follower]; ok {
			channels = append(channels, channel)
		}
	}
	return channels, nil
}
