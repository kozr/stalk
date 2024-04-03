package redis

func userUrlPrefixKey(userId string) string {
	return "user_url:" + userId
}

func userChannelPrefixKey(userId string) string {
	return "user_channel:" + userId
}
