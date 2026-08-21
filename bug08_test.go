package tsrollup

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func TestBug08_CompactHonorsContextCancel(t *testing.T) {
	dir := t.TempDir()
	e, err := Open(Options{
		Root:   filepath.Join(dir, "e"),
		Window: WindowSpec{Size: time.Minute},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer e.Close()
	e.mu.Lock()
	unlocked := false
	unlock := func() {
		if !unlocked {
			unlocked = true
			e.mu.Unlock()
		}
	}
	defer unlock()

	done := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		done <- e.Compact(ctx)
	}()
	time.Sleep(30 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected cancel error")
		}
		if !errors.Is(err, ErrCanceled) && !errors.Is(err, context.Canceled) {
			t.Fatalf("want canceled, got %v", err)
		}
	case <-time.After(300 * time.Millisecond):
		// 释放锁让卡死的 Compact 退出，避免 go test 进程挂死
		unlock()
		t.Fatal("Compact ignored ctx cancel while waiting for lock")
	}
}
