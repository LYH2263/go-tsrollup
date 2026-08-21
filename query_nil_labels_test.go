package tsrollup

import (
	"testing"
	"time"
)

// 复现：采集端只给系列名（Labels=nil）建系成功，Query 盖 __queried__ 戳时不再 nil map panic。
func TestQueryNilLabelsNoPanic(t *testing.T) {
	dir := t.TempDir()
	eng, err := Open(Options{Root: dir, Window: WindowSpec{Size: time.Minute}})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer eng.Close()

	// 只有 Series 名，Labels=nil —— 之前 normalizeLabels 返回 nil，建系成功但埋下 nil labels。
	if err := eng.Append(Sample{Series: "cpu", Ts: time.Unix(0, 0), Value: 1}); err != nil {
		t.Fatalf("Append: %v", err)
	}

	// 盖戳路径：曾经 ls.labels["__queried__"]="1" 对 nil map 写入 panic。
	res, err := eng.Query("cpu", nil, time.Time{}, time.Unix(60, 0))
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("want 1 result, got %d", len(res))
	}
	if got := res[0].Labels["__queried__"]; got != "1" {
		t.Fatalf("want __queried__=1, got %q (Labels=%v)", got, res[0].Labels)
	}
}

// 保证空标签系列自身的 SeriesID 稳定、与带标签系列区分。
func TestNormalizeLabelsReturnsNonNil(t *testing.T) {
	out, err := normalizeLabels(nil)
	if err != nil {
		t.Fatalf("normalizeLabels(nil): %v", err)
	}
	if out == nil {
		t.Fatal("normalizeLabels(nil) returned nil map")
	}
	out["k"] = "v" // 写入不应 panic
	if out["k"] != "v" {
		t.Fatalf("expect k=v, got %q", out["k"])
	}
}
