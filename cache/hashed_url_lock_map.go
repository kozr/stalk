package cache

import (
	"sync"

	"github.com/kozr/stalk/special_locks"
)

var locksMap sync.Map

func GetOrLoadUserLock(userId string) *special_locks.RecencyLock {
	lock, _ := locksMap.LoadOrStore(userId, special_locks.NewRecencyLock())
	return lock.(*special_locks.RecencyLock)
}

// Be extra careful, only remove if you are sure that the lock is not being used anymore.
func RemoveUserLock(userId string) {
	locksMap.Delete(userId)
}
