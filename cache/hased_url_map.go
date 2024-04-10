package cache

import (
	"sync"
)

var inMemoryUrlMap = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func UpdateUserHashedUrl(userId string, hashedUrl string) error {
	inMemoryUrlMap.Lock()
	defer inMemoryUrlMap.Unlock()
	inMemoryUrlMap.m[userId] = hashedUrl
	return nil
}

func RemoveUserHashedUrl(userId string) error {
	inMemoryUrlMap.Lock()
	defer inMemoryUrlMap.Unlock()
	delete(inMemoryUrlMap.m, userId)
	return nil
}

func GetUserHashedUrl(userId string) (string, error) {
	inMemoryUrlMap.RLock()
	defer inMemoryUrlMap.RUnlock()
	url, ok := inMemoryUrlMap.m[userId]
	if !ok {
		return "", ErrCacheMiss("url not found for user id: " + userId)
	}
	return url, nil
}
