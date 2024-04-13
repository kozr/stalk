package special_locks

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestRecencyLock(t *testing.T) {
	rl := NewRecencyLock()
	rl.Lock(1)
	if rl.Lock(0) != LockFailAbort {
		t.Error("Lock failed to abort")
	}
	rl.Unlock()
	if rl.Lock(0) != LockFailAbort {
		t.Error("Lock failed to abort")
	}
	if rl.Lock(2) != LockSuccess {
		t.Error("Lock failed to succeed")
	}
}

func TestRecencyLockConcurrently(t *testing.T) {
	const numGoroutines = 1000
	var wg sync.WaitGroup
	var successfulTimestamps []int64
	startSignal := make(chan struct{})

	wg.Add(numGoroutines)

	lock := NewRecencyLock()
	for i := 0; i < numGoroutines; i++ {
		go func(timestamp int64) {
			defer wg.Done()

			<-startSignal
			time.Sleep(time.Duration(rand.Intn(numGoroutines)/10) * time.Millisecond)

			response := lock.Lock(timestamp)
			if response == LockSuccess {
				successfulTimestamps = append(successfulTimestamps, timestamp)
				lock.Unlock()
			}
		}(int64(i))
	}

	close(startSignal)

	wg.Wait()

	// successful timestamps are in order
	for i := 0; i < len(successfulTimestamps)-1; i++ {
		if successfulTimestamps[i] > successfulTimestamps[i+1] {
			t.Error("Timestamps are not in order")
		}
	}
}
