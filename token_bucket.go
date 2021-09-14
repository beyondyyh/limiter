package limiter

import (
	"sync"
	"time"
)

// 令牌桶算法
// 基本思路：
// 固定的速率向桶中放入令牌，桶满了则不再加入令牌；
// 服务受到请求时尝试从桶中取一个令牌，如果取到则继续后续逻辑，取不到则直接返回频率超限错误码。

type TokenBucket struct {
	rate         int64      // 固定的token放入速率，req/s
	capacity     int64      // 桶的容量
	tokens       int64      // 桶中当前token数量
	lastTokenSec int64      // 上次向桶中放令牌的时间的时间戳，单位：秒
	mu           sync.Mutex // 互斥锁
}

func NewTokenBucket(rate, cap int64) *TokenBucket {
	return &TokenBucket{
		rate:         rate,
		capacity:     cap,
		tokens:       0,
		lastTokenSec: time.Now().Unix(),
	}
}

// func (l *TokenBucket) Take() bool {
func (l *TokenBucket) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().Unix()
	// 时间差*放入速率得到时间间隔内需要放入的令牌数，一次添加多个令牌
	l.tokens += (now - l.lastTokenSec) * l.rate
	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}
	l.lastTokenSec = now
	// 还有令牌，领取令牌
	if l.tokens > 0 {
		l.tokens--
		return true
	}
	// 没有令牌，则拒绝
	return false
}
