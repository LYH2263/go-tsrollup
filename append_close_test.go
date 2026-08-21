package tsrollup

import (
	"context"
	"errors"
	"testing"
	"time"
)

// 回归 bug3：Close 后再 AppendContext 必须返回 ErrClosed，而不是在
// 已被置空的 wal 上空指针 panic。复现运维"先 Close 再误点补点"的路径。
func TestAppendContextAfterCloseReturnsErrClosed(t *testing.T) {
	dir := t.TempDir()
	eng, err := Open(Options{Root: dir})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	// 先写一条正常样本，确保 wal 路径本身可用。
	if err := eng.AppendContext(context.Background(), Sample{
		Series: "cpu.load",
		Ts:     time.Now().UTC(),
		Value:  1.0,
	}); err != nil {
		t.Fatalf("append before close: %v", err)
	}

	if err := eng.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// Close 已把 e.wal 置 nil；这里历史上会 panic（nil 解引用）。
	err = eng.AppendContext(context.Background(), Sample{
		Series: "cpu.load",
		Ts:     time.Now().UTC(),
		Value:  2.0,
	})
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("after close: want ErrClosed, got %v", err)
	}
}

// 并发场景：Close 与 AppendContext 竞态。即便 Close 正在执行（已置 closed
// 但尚未执行到 e.wal=nil，或反之），AppendContext 也不应踩到 nil wal。
func TestAppendContextConcurrentCloseNoPanic(t *testing.T) {
	dir := t.TempDir()
	eng, err := Open(Options{Root: dir})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	stop := make(chan struct{})
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		for {
			err := eng.AppendContext(context.Background(), Sample{
				Series: "mem.used",
				Ts:     time.Now().UTC(),
				Value:  1.0,
			})
			// 关闭后唯一可接受的错误是 ErrClosed；任何 panic 都会让测试炸掉。
			if err != nil && !errors.Is(err, ErrClosed) && !errors.Is(err, ErrCanceled) {
				t.Errorf("append during close: unexpected err %v", err)
				return
			}
			select {
			case <-stop:
				return
			default:
			}
		}
	}()

	// 给写入 goroutine 一点起跑时间，再关闭。
	time.Sleep(10 * time.Millisecond)
	if err := eng.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	// 收尾：再补几条，触发"关闭后写入"路径。
	for i := 0; i < 5; i++ {
		_ = eng.AppendContext(context.Background(), Sample{
			Series: "mem.used",
			Ts:     time.Now().UTC(),
			Value:  float64(i),
		})
	}
	close(stop)
	<-finished
}
