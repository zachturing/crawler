package collect

import (
	"sync"
	"testing"
)

func BenchmarkIncWithLock(b *testing.B) {
	var counter int64
	var lock sync.Mutex
	for i := 0; i < b.N; i++ {
		incWithLock(&counter, &lock)
	}
	// b.Logf("%v", counter)
}

func BenchmarkIncWithAtomic(b *testing.B) {
	var counter int64
	for i := 0; i < b.N; i++ {
		incWithAtomic(&counter)
	}
	// b.Logf("%v", counter)

}

func TestPrint(t *testing.T) {
	tests := []struct {
		name string
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			send()
		})
	}
}
