package special_locks

import "sync"

type LockResult int

const (
	LockSuccess LockResult = iota
	LockFailAbort
)

type RecencyLock struct {
	lock            sync.Mutex
	cond            *sync.Cond
	latestTimestamp int64 // Timestamp of the sole request we care about.
	active          bool  // Indicates whether the lock is currently held.
}

// NewRecencyLock creates a new RecencyLock instance.
// The RecencyLock will only be held by one request at a time.
// If a request with a lower timestamp than the current lock's timestamp is received,
// the lock will be aborted to prevent priority inversion.
func NewRecencyLock() *RecencyLock {
	pl := &RecencyLock{}
	pl.cond = sync.NewCond(&pl.lock)
	return pl
}

func (pl *RecencyLock) Lock(timestamp int64) LockResult {
	pl.lock.Lock()
	defer pl.lock.Unlock()

	// If the current timestamp is lower than the next queued priority,
	// recommend aborting to prevent priority inversion.
	if timestamp < pl.latestTimestamp {
		return LockFailAbort
	}

	pl.latestTimestamp = timestamp

	for pl.active {
		pl.cond.Wait()
		if timestamp < pl.latestTimestamp {
			return LockFailAbort
		}
	}

	pl.active = true
	return LockSuccess
}

func (pl *RecencyLock) Unlock() {
	pl.lock.Lock()
	defer pl.lock.Unlock()

	pl.active = false
	pl.cond.Broadcast()
}
