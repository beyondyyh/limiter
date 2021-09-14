package limiter

// 静态时间窗口，也叫访问计数

import (
	"sync"
	"time"
)

type RateLimiter struct {
	rate   int           // 阈值
	begin  time.Time     // 计数开始时间
	period time.Duration // 计数周期
	count  int           // 收到的请求数
	mu     sync.Mutex    // 锁
}

func (l *RateLimiter) NewRateLimiter(rate int, period time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:   rate,
		begin:  time.Now(),
		period: period,
		count:  0,
	}
}

func (l *RateLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 判断收到的请求数是否达到阈值
	if l.count == l.rate-1 {
		now := time.Now()
		// 达到阈值后，当前请求如果不在计数周期内，重置计数器
		if now.Sub(l.begin) >= l.period {
			l.Reset(now)
			return true
		}
		return false
	}
	// 请求未达到阈值，计数加一
	l.count++
	return true
}

func (l *RateLimiter) Reset(begin time.Time) {
	l.begin = begin
	l.count = 0
}
