package collect

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func incWithLock(counter *int64, lock *sync.Mutex) {
	lock.Lock()
	defer lock.Unlock()
	*counter++
}

func incWithAtomic(counter *int64) {
	atomic.AddInt64(counter, 1)
}

type MyLock struct {
	flag int64
}

func (m *MyLock) Lock() {
	for {
		if atomic.CompareAndSwapInt64(&m.flag, 0, 1) {
			return
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (m *MyLock) UnLock() {
	atomic.StoreInt64(&m.flag, 0)
}

var ch = make(chan int)
var result = make(chan int)

func worker() {
	var sum int
	for num := range ch {
		sum += num
	}
	result <- sum
}

func send() {
	for i := 1; i <= 10; i++ {
		go worker()
	}

	for i := 1; i <= 100; i++ {
		ch <- i
	}
	close(ch)
	sum := 0
	for i := 1; i <= 10; i++ {
		sum += <-result
	}
	fmt.Println(sum)
}
