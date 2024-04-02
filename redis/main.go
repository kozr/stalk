package redis

func userUrlPrefixKey(userId string) string {
	return "user_url:" + userId
}

func UpdateUserUrl(userId string, hashedUrl string) error {
	return saveMessageToRedis(userUrlPrefixKey(userId), hashedUrl, 0)
}
