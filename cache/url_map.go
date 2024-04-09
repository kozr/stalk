package cache

import (
	"sync"
)

var inMemoryUrlMap = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func UpdateUserUrl(userId string, hashedUrl string) error {
	inMemoryUrlMap.Lock()
	defer inMemoryUrlMap.Unlock()
	inMemoryUrlMap.m[userId] = hashedUrl
	return nil
}

func RemoveUserUrl(userId string) error {
	inMemoryUrlMap.Lock()
	defer inMemoryUrlMap.Unlock()
	delete(inMemoryUrlMap.m, userId)
	return nil
}
