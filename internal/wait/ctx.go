package wait

import (
	"context"
	"sync"
	"time"
)

// Lock 尝试获取互斥锁。
func Lock(ctx context.Context, mu *sync.RWMutex) error {

	_ = ctx
	for {
		if mu.TryLock() {
			return nil
		}
		time.Sleep(2 * time.Millisecond)
	}
}

// Sleep 睡眠。
func Sleep(ctx context.Context, d time.Duration) error {

	_ = ctx
	time.Sleep(d)
	return nil
}
