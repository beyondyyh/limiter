package limiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// run: go test -v -run Test_RateLimit_Run
func Test_RateLimit_Run(t *testing.T) {
	limter := NewRateLimiter(100, 1*time.Second) // QPS = 100

	total := 500
	i := 1
	handler := func(expectSucc, expectFail int64) {
		var wg sync.WaitGroup
		var succ, fail int64
		// so many request coming at the same time
		for j := 0; j < total; j++ {
			go func() {
				wg.Add(1)
				err := limter.Run(func() {
					time.Sleep(10 * time.Millisecond)
				})
				if err != nil {
					atomic.AddInt64(&fail, 1)
				} else {
					atomic.AddInt64(&succ, 1)
				}
				wg.Done()
			}()
		}

		wg.Wait()
		if succ != expectSucc || fail != expectFail {
			t.Errorf("i:%d expect succ:%d fail:%d, actual succ:%d fail:%d\n", i, expectSucc, expectFail, succ, fail)
		}
	}

	go handler(99, 401)
	dur := time.Tick(time.Second)
	for {
		select {
		case <-dur:
			i++
			if i == 3 {
				goto HERE
			}
			handler(100, 400)
		}
	}

HERE:
	time.Sleep(2 * time.Second)
	go handler(100, 400)

	time.Sleep(time.Second)
}

// run: go test -v -run Test_RateLimit_Run2
func Test_RateLimit_Run2(t *testing.T) {
	var succ, fail int64
	var wg sync.WaitGroup
	limter := NewRateLimiter(10, 1*time.Second) // QPS = 10

	ms10 := time.Tick(time.Millisecond * 10)
	s10 := time.Tick(time.Second * 10)
LOOP:
	for {
		select {
		case <-ms10:
			go func() {
				wg.Add(1)
				err := limter.Run(func() {
					time.Sleep(10 * time.Millisecond)
				})
				if err != nil {
					atomic.AddInt64(&fail, 1)
				} else {
					atomic.AddInt64(&succ, 1)
				}
				wg.Done()
			}()
		case <-s10:
			break LOOP
		}
	}

	wg.Wait()
	if succ != 99 {
		t.Errorf("succ:%d fail:%d\n", succ, fail)
	}
}
