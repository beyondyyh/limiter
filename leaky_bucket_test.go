package limiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// run: go test -v -run Test_LeakyBucket_Run
func Test_LeakyBucket_Run(t *testing.T) {
	limter := NewLeakyBucket(10, 100)

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

	go handler(100, 400)
	dur := time.Tick(time.Second)
	for {
		select {
		case <-dur:
			i++
			if i == 3 {
				goto HERE
			}
			handler(10, 490)
		}
	}

HERE:
	time.Sleep(2 * time.Second)
	go handler(30, 470)

	time.Sleep(time.Second)
}
