package cache

import (
	"sync"
)

var inMemoryChannelMap = struct {
	sync.RWMutex
	m map[string]chan string
}{m: make(map[string]chan string)}

func UpdateUserChannel(userId string, channel chan string) error {
	inMemoryChannelMap.Lock()
	defer inMemoryChannelMap.Unlock()
	inMemoryChannelMap.m[userId] = channel
	return nil
}

func RemoveUserChannel(userId string) error {
	inMemoryChannelMap.Lock()
	defer inMemoryChannelMap.Unlock()
	delete(inMemoryChannelMap.m, userId)
	return nil
}

func GetUserChannel(userId string) (chan string, error) {
	inMemoryChannelMap.RLock()
	defer inMemoryChannelMap.RUnlock()
	channel, ok := inMemoryChannelMap.m[userId]
	if !ok {
		return nil, ErrCacheMiss
	}
	return channel, nil
}
