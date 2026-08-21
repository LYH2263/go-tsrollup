package wait

import (
	"context"
	"sync"
	"time"
)

// Lock 尝试获取互斥锁，同时尊重 ctx 取消（bug08 正确行为）。
func Lock(ctx context.Context, mu *sync.RWMutex) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if mu.TryLock() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Millisecond):
		}
	}
}

// Sleep 可取消睡眠。
func Sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
