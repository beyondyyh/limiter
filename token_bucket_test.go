package limiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// run: go test -v -run Test_TokenBucket_Run
func Test_TokenBucket_Run(t *testing.T) {
	tb := NewTokenBucket(&Config{
		QPS:      100,
		MaxCount: 200,
	})

	dur := time.Tick(time.Second)

	total := 500
	i := 0

	var handler = func(expectSucc, expectFail int64) {
		var wg sync.WaitGroup
		var succ, fail int64
		// so many requests coming at the same time
		for j := 0; j < total; j++ {
			go func() {
				wg.Add(1)

				err := tb.Run(func() {
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
			t.Errorf("%d expect succ:%d fail:%d, actual succ:%d fail:%d\n", i, expectSucc, expectFail, succ, fail)
		}
	}

	go handler(100, 400)

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
	go handler(200, 300)

	time.Sleep(3 * time.Second)
	go handler(200, 300)

	time.Sleep(time.Second)
}

// run: go test -v -run Test_TokenBucket_Run2
func Test_TokenBucket_Run2(t *testing.T) {
	var succ, fail int64
	var wg sync.WaitGroup

	tb := NewTokenBucket(&Config{
		QPS:      1.2,
		MaxCount: 2,
	})

	ms20 := time.Tick(time.Millisecond * 20)
	s20 := time.Tick(time.Second * 20)
LOOP:
	for {
		select {
		case <-ms20:
			go func() {
				wg.Add(1)

				err := tb.Run(func() {
					time.Sleep(10 * time.Millisecond)
				})
				if err != nil {
					atomic.AddInt64(&fail, 1)
				} else {
					atomic.AddInt64(&succ, 1)
				}

				wg.Done()
			}()
		case <-s20:
			break LOOP
		}
	}

	wg.Wait()
	if succ != 24 {
		t.Errorf("succ:%d fail:%d\n", succ, fail)
	}
}
