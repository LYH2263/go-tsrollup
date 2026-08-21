package tsrollup

import (
	"context"
	"errors"
	"testing"
	"time"
)

// newTestEngine 打开一个临时引擎，测试结束自动清理。
func newTestEngine(t *testing.T, opts Options) *Engine {
	t.Helper()
	if opts.Root == "" {
		opts.Root = t.TempDir()
	}
	e, err := Open(opts)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = e.Close() })
	return e
}

func validSample() Sample {
	return Sample{
		Series: "cpu",
		Labels: map[string]string{"host": "a"},
		Ts:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Value:  1.0,
	}
}

// 已取消的 ctx 必须让 AppendContext 立即返回 ErrCanceled，
// 而不是像修复前那样丢弃 ctx 仍把记录写完。
func TestAppendContext_RespectsCanceledCtx(t *testing.T) {
	e := newTestEngine(t, Options{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := e.AppendContext(ctx, validSample())
	if !errors.Is(err, ErrCanceled) && !errors.Is(err, context.Canceled) {
		t.Fatalf("expected ErrCanceled/context.Canceled, got %v", err)
	}
	if e.walRecs.Load() != 0 {
		t.Fatalf("WAL record count should be 0 after cancel, got %d", e.walRecs.Load())
	}
}

// 已超时的 ctx 同样必须在写 WAL 前返回。
func TestAppendContext_RespectsDeadlineExceeded(t *testing.T) {
	e := newTestEngine(t, Options{})

	ctx, cancel := context.WithTimeout(context.Background(), -1*time.Millisecond)
	defer cancel()

	err := e.AppendContext(ctx, validSample())
	if !errors.Is(err, ErrCanceled) && !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected ErrCanceled/context.DeadlineExceeded, got %v", err)
	}
	if e.walRecs.Load() != 0 {
		t.Fatalf("WAL record count should be 0 after deadline, got %d", e.walRecs.Load())
	}
}

// 未取消的 ctx 仍可正常写入。
func TestAppendContext_NormalCtxAppends(t *testing.T) {
	e := newTestEngine(t, Options{})

	if err := e.AppendContext(context.Background(), validSample()); err != nil {
		t.Fatalf("AppendContext: %v", err)
	}
	if e.walRecs.Load() != 1 {
		t.Fatalf("WAL record count should be 1, got %d", e.walRecs.Load())
	}
}
