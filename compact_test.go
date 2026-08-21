package tsrollup

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup/internal/clock"
)

// TestCompactFlushFailureKeepsOpenWindows 覆盖 bug06：
// 压测把 segments 目录替换成普通文件，使 flush 失败。Compact 必须返回
// ErrPersist，但不得删除 openWins，也不得提升只读视图——否则同进程再
// Query 已读不到未刷盘的聚合点。
func TestCompactFlushFailureKeepsOpenWindows(t *testing.T) {
	root := t.TempDir()
	segDir := filepath.Join(root, "segments")

	// 固定时钟：让窗口落在 [t0, t0+1m)，推进后该窗口即到期关闭。
	t0 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	clk := &clock.Fake{T: t0}

	e, err := Open(Options{
		Root:    root,
		Window:  WindowSpec{Size: time.Minute, Align: AlignFloor},
		AggKind: "sum",
		Clock:   clk,
	})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() {
		// Compact 失败后窗口仍留在 openWins，Close 会再次尝试 flush。
		// 先恢复 segments 目录，让 Close 的末窗 flush 成功并正常关闭 WAL，
		// 避免 Windows 下 WAL 句柄泄漏卡住 TempDir 清理。
		_ = os.Remove(segDir)
		_ = os.MkdirAll(segDir, 0o755)
		_ = e.Close()
	})

	// 写入一个落在当前窗口内的采样。
	if err := e.Append(Sample{
		Series: "cpu",
		Ts:     t0.Add(10 * time.Second),
		Value:  42,
	}); err != nil {
		t.Fatalf("Append: %v", err)
	}

	// 推进时钟越过窗口终点，使窗口在 Compact 时被判为已关闭。
	clk.Advance(2 * time.Minute)

	// 破坏 segments 目录：替换为普通文件，使 segment.Create 内部的
	// os.MkdirAll 失败 -> flushRowsLocked 返回错误。
	if err := os.RemoveAll(segDir); err != nil {
		t.Fatalf("remove segments dir: %v", err)
	}
	if err := os.WriteFile(segDir, []byte("not a dir"), 0o644); err != nil {
		t.Fatalf("plant segments file: %v", err)
	}

	// Compact 必须返回 ErrPersist（持久化失败是符合预期的）。
	err = e.Compact(context.Background())
	if !errors.Is(err, ErrPersist) {
		t.Fatalf("Compact err = %v, want ErrPersist", err)
	}

	// 关键不变量：flush 失败时不得删除 openWins。
	stats := e.Stats()
	if got := stats.WindowsClosed; got != 0 {
		t.Errorf("WindowsClosed after failed Compact = %d, want 0 (openWins must not be deleted)", got)
	}
	if got := stats.Segments; got != 0 {
		t.Errorf("Segments after failed Compact = %d, want 0 (readonly view must not be promoted)", got)
	}

	// 同进程再 Query 必须仍能读到未刷盘的聚合点。
	res, err := e.Query("cpu", nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("Query results = %d, want 1", len(res))
	}
	pts := res[0].Points
	if len(pts) != 1 {
		t.Fatalf("Points = %d, want 1 (unpersisted aggregation must still be readable)", len(pts))
	}
	if pts[0].Value != 42 {
		t.Errorf("Point value = %v, want 42", pts[0].Value)
	}
	if pts[0].Count != 1 {
		t.Errorf("Point count = %d, want 1", pts[0].Count)
	}
}
