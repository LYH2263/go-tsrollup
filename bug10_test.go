package tsrollup_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup"
	"github.com/LYH2263/go-tsrollup/internal/segment"
)

func TestBug10_CloseFlushesBeforeDropWindows(t *testing.T) {
	dir := t.TempDir()
	root := filepath.Join(dir, "e")
	e, err := tsrollup.Open(tsrollup.Options{
		Root:   root,
		Window: tsrollup.WindowSpec{Size: time.Hour},
	})
	if err != nil {
		t.Fatal(err)
	}
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	if err := e.Append(tsrollup.Sample{Series: "last", Ts: ts, Value: 7}); err != nil {
		t.Fatal(err)
	}
	if err := e.Close(); err != nil {
		t.Fatal(err)
	}
	paths, err := filepath.Glob(filepath.Join(root, "segments", "*.tsr"))
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatal("last window lost: Close dropped memory before flush (no segment)")
	}
	v, err := segment.OpenView(paths[0])
	if err != nil {
		t.Fatal(err)
	}
	rows := v.RowsFor("last")
	if len(rows) == 0 || rows[0].Value != 7 {
		t.Fatalf("segment rows=%+v", rows)
	}
}
