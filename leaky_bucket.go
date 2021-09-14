package limiter

import (
	"math"
	"sync"
	"time"
)

// 漏桶算法
// 基本思路：
// 一个固定大小的桶，请求按照固定的速率流出；
// 请求数大于桶的容量时，丢弃多余的请求。

type LeakyBucket struct {
	rate       float64    // 每秒固定流出速率
	capacity   float64    // 桶的容量
	water      float64    // 当前桶中请求量
	lastLeakMs int64      // 上次桶漏水的时间戳，单位：微秒
	mu         sync.Mutex // 互斥锁
}

func NewLeakyBucket(rate, cap float64) *LeakyBucket {
	return &LeakyBucket{
		rate:       rate,
		capacity:   cap,
		water:      0,
		lastLeakMs: time.Now().UnixNano() / 1e6,
	}
}

func (l *LeakyBucket) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	// 计算剩余水量，两次执行时间中需要漏掉的水
	leakyWater := l.water - (float64(now-l.lastLeakMs) * l.rate / 1000)
	l.water = math.Max(0, leakyWater)
	l.lastLeakMs = now
	if l.water+1 <= l.capacity {
		l.water++
		return true
	}

	return false
}

func (l *LeakyBucket) Run(fn func()) error {
	if !l.Allow() {
		return ErrFreqExceed
	}

	// Invoke fn
	fn()

	return nil
}
