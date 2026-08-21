package tsrollup

import (
	"testing"
	"time"
)

// TestQuery_PointsNotSharedWithCache 回归：调用方改写返回 Points 的某个 Value 做高亮，
// 不应污染内部 ptsCache，下一次同进程 Query 同一位置不能出现 999。
func TestQuery_PointsNotSharedWithCache(t *testing.T) {
	dir := t.TempDir()
	e, err := Open(Options{Root: dir, AggKind: "sum"})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer e.Close()

	base := time.Date(2026, 8, 22, 9, 0, 0, 0, time.UTC)
	if err := e.Append(Sample{Series: "cpu", Ts: base, Value: 1}); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if err := e.Append(Sample{Series: "cpu", Ts: base.Add(30 * time.Second), Value: 3}); err != nil {
		t.Fatalf("Append: %v", err)
	}

	// 第一次查询，拿到 Points 后在调用侧改第一个 Value 做高亮。
	r1, err := e.Query("cpu", nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("Query 1: %v", err)
	}
	if len(r1) != 1 || len(r1[0].Points) != 1 {
		t.Fatalf("want 1 result with 1 point, got %+v", r1)
	}
	if r1[0].Points[0].Value != 4 {
		t.Fatalf("want agg value 4, got %v", r1[0].Points[0].Value)
	}
	r1[0].Points[0].Value = 999 // 模拟高亮污染

	// 紧接着再查一次，同位置不应为 999。
	r2, err := e.Query("cpu", nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("Query 2: %v", err)
	}
	if len(r2) != 1 || len(r2[0].Points) != 1 {
		t.Fatalf("want 1 result with 1 point on re-query, got %+v", r2)
	}
	if r2[0].Points[0].Value == 999 {
		t.Fatalf("cache leaked: second query saw mutated 999, want 4")
	}
	if r2[0].Points[0].Value != 4 {
		t.Fatalf("want 4 on re-query, got %v", r2[0].Points[0].Value)
	}
}
