package redis

func UpdateUserUrl(userId string, hashedUrl string) error {
	return saveKeyToRedis(userUrlPrefixKey(userId), hashedUrl, 0)
}

func RemoveUserUrl(userId string) error {
	return removeKeyFromRedis(userUrlPrefixKey(userId))
}

func UpdateUserChannel(userId string, channel string) error {
	return saveKeyToRedis(userChannelPrefixKey(userId), channel, 0)
}

func RemoveUserChannel(userId string) error {
	return removeKeyFromRedis(userChannelPrefixKey(userId))
}
