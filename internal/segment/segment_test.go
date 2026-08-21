package segment_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup/internal/errs"
	"github.com/LYH2263/go-tsrollup/internal/segment"
)

func TestWriteRead(t *testing.T) {
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
		t.Fatal(err)
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

func TestCorruptWraps(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.tsr")
	if err := os.WriteFile(path, []byte("not-a-segment"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := segment.OpenView(path)
	if err == nil || !errors.Is(err, errs.ErrCorrupt) {
		t.Fatalf("want corrupt got %v", err)
	}
}
