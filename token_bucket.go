package limiter

import (
	"errors"
	"sync"
	"time"
)

var ErrNoEnoughToken = errors.New("no enouth token")

// TokenBucket token bucket
// Ref: https://en.wikipedia.org/wiki/Token_bucket
type TokenBucket struct {
	*Config
	lastAccessTime time.Time
	surplus        int64
	m              sync.Mutex
}

var _ Limiter = &TokenBucket{}

// NewTokenBucket create one token bucket by config
func NewTokenBucket(conf *Config) Limiter {
	return &TokenBucket{
		Config:         conf,
		lastAccessTime: time.Now(),
		surplus:        int64(conf.QPS),
	}
}

// Get get count token from bucke
func (tb *TokenBucket) Get(count int64) error {
	tb.m.Lock()
	defer tb.m.Unlock()

	if tb.surplus >= count {
		tb.surplus -= count
		return nil
	}

	tb.refreshUnSafe()

	if tb.surplus >= count {
		tb.surplus -= count
		return nil
	}

	return ErrNoEnoughToken
}

// Put put count token to bucket
func (tb *TokenBucket) Put(count int64) error {
	tb.m.Lock()
	defer tb.m.Unlock()

	tb.refreshUnSafe()

	tb.surplus += count
	if tb.MaxCount > 0 && tb.surplus > tb.MaxCount {
		tb.surplus = tb.MaxCount
	}

	return nil
}

// Run Invoke fn if get one token success, otherwise do nothing and returns error
func (tb *TokenBucket) Run(fn func()) error {
	if err := tb.Get(1); err != nil {
		return err
	}

	// Invoke fn
	fn()

	return nil
}

func (tb *TokenBucket) refreshUnSafe() {
	now := time.Now()
	span := now.Sub(tb.lastAccessTime).Seconds()

	surplus := span * tb.QPS
	if surplus < 1 {
		return
	}

	tb.surplus += int64(surplus)
	// if maxCount be set, limit the surplus
	if tb.MaxCount > 0 && tb.surplus > tb.MaxCount {
		tb.surplus = tb.MaxCount
	}

	tb.lastAccessTime = now
}
