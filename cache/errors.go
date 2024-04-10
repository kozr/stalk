package cache

import (
	"errors"
)

func ErrCacheMiss(message string) error {
	return errors.New("cache miss: " + message)
}
