package segment_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup/internal/segment"
)

func TestBug09_SegmentSyncBeforeRename(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.tsr")
	w, err := segment.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Unix(100, 0).UTC()
	if err := w.WriteRow(segment.Row{
		Series: "s",
		Start:  start,
		End:    start.Add(time.Minute),
		Value:  3,
		Count:  1,
		Agg:    "sum",
	}); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close failed (likely Sync after Close): %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) < 16 {
		t.Fatalf("segment too small: %d", len(got))
	}
	v, err := segment.OpenView(path)
	if err != nil {
		t.Fatal(err)
	}
	rows := v.RowsFor("s")
	if len(rows) != 1 || rows[0].Value != 3 {
		t.Fatalf("%+v", rows)
	}
}
