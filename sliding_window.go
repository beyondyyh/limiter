package limiter

// 滑动时间窗口算法

import (
	"sync"
	"time"
)

type timeSlot struct {
	begin time.Time // 此time slot的时间起点
	count int       // 落在此time slot内的请求数
}

type SlidingWindowLimiter struct {
	slotDuration time.Duration // time slot的长度
	winDuration  time.Duration // sliding window的长度
	maxReq       int           // 大窗口时间内允许的最大请求数
	windows      []*timeSlot   // window内slot的合集
	mu           sync.Mutex    // 互斥锁保护其他字段
}

func NewSlidingWindowLimiter(slotDuration, winDuration time.Duration, maxReq int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		slotDuration: slotDuration,
		winDuration:  winDuration,
		maxReq:       maxReq,
	}
}

func (l *SlidingWindowLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	timeoutOffset := -1
	// 将已经过期的time slot移出时间窗口
	for i, ts := range l.windows {
		if ts.begin.Add(l.winDuration).After(now) {
			break
		}
		timeoutOffset = i
	}
	if timeoutOffset > -1 {
		l.windows = l.windows[timeoutOffset+1:]
	}

	// 判断请求是否超限
	var ok bool
	if l.countReq() < l.maxReq {
		ok = true
	}

	// 记录这次的请求数
	var lastSlot *timeSlot
	if len(l.windows) > 0 {
		lastSlot = l.windows[len(l.windows)-1]
		// 如果当前时间已经超过这个时间插槽的跨度，那么新建一个时间插槽
		if lastSlot.begin.Add(l.slotDuration).Before(now) {
			lastSlot = &timeSlot{begin: now, count: 1}
			l.windows = append(l.windows, lastSlot)
		} else {
			lastSlot.count++
		}
	} else {
		lastSlot = &timeSlot{begin: now, count: 1}
		l.windows = append(l.windows, lastSlot)
	}

	return ok
}

func (l *SlidingWindowLimiter) countReq() int {
	var count int
	for _, ts := range l.windows {
		count += ts.count
	}
	return count
}

func (l *SlidingWindowLimiter) Run(fn func()) error {
	if !l.Allow() {
		return ErrFreqExceed
	}

	// Invoke fn
	fn()

	return nil
}
