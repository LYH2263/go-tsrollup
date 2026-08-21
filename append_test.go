package tsrollup

import (
	"testing"
	"time"
)

// TestAppend_LabelsIsolation 复现管理页“复用同一份 Labels map，Append
// 返回后改 host 打审计日志”的场景：外部事后改动 s.Labels 不得穿透进
// ListSeries / 倒排 byID / Match（系列 ID 须对得上）。
func TestAppend_LabelsIsolation(t *testing.T) {
	dir := t.TempDir()
	eng, err := Open(Options{Root: dir, Window: WindowSpec{Size: time.Minute}})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer eng.Close()

	labels := map[string]string{"host": "a"}
	wantID := SeriesID("cpu", labels)

	if err := eng.Append(Sample{
		Series: "cpu",
		Labels: labels,
		Ts:     time.Unix(0, 0).UTC(),
		Value:  1,
	}); err != nil {
		t.Fatalf("Append: %v", err)
	}

	// 模拟管理页“Append 返回后改 host 打审计日志”。
	labels["host"] = "MUTATED"

	// ListSeries 必须仍看到原始 host=a，且 ID 不变。
	got, err := eng.ListSeries()
	if err != nil {
		t.Fatalf("ListSeries: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("ListSeries: want 1 series, got %d", len(got))
	}
	if got[0].ID != wantID {
		t.Fatalf("ListSeries ID: want %q, got %q", wantID, got[0].ID)
	}
	if got[0].Labels["host"] != "a" {
		t.Fatalf("ListSeries host: want %q, got %q", "a", got[0].Labels["host"])
	}

	// 倒排索引 byID 不得被外部改动污染。
	if got := eng.index.LabelsOf(wantID); got["host"] != "a" {
		t.Fatalf("LabelsOf host: want %q, got %q", "a", got["host"])
	}

	// 用原始标签选择器必须仍能命中（系列 ID 对得上）。
	ids := eng.index.Match("cpu", map[string]string{"host": "a"})
	if len(ids) != 1 || ids[0] != wantID {
		t.Fatalf("Match host=a: want [%s], got %v", wantID, ids)
	}

	// MUTATED 不应作为已知取值出现。
	for _, v := range eng.index.Values("host") {
		if v == "MUTATED" {
			t.Fatalf("Values(host): leaked MUTATED into inverted index")
		}
	}
}
