package cache

func userChannelPrefixKey(userId string) string {
	return "user_channel:" + userId
}

func userUrlPrefixKey(userId string) string {
	return "user_url:" + userId
}
