package tsrollup_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup"
	"github.com/LYH2263/go-tsrollup/internal/clock"
)

func TestBug06_CompactFlushRollback(t *testing.T) {
	dir := t.TempDir()
	root := filepath.Join(dir, "e")
	fake := &clock.Fake{T: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	e, err := tsrollup.Open(tsrollup.Options{
		Root:   root,
		Window: tsrollup.WindowSpec{Size: time.Second},
		Clock:  fake,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer e.Close()
	if err := e.Append(tsrollup.Sample{
		Series: "m",
		Ts:     fake.Now(),
		Value:  5,
	}); err != nil {
		t.Fatal(err)
	}
	fake.Advance(2 * time.Second)
	segDir := filepath.Join(root, "segments")
	ents, _ := os.ReadDir(segDir)
	for _, ent := range ents {
		_ = os.Remove(filepath.Join(segDir, ent.Name()))
	}
	_ = os.Remove(segDir)
	if err := os.WriteFile(segDir, []byte("not-a-dir"), 0o644); err != nil {
		t.Fatal(err)
	}
	err = e.Compact(context.Background())
	if err == nil {
		t.Fatal("expected persist error")
	}
	_ = os.Remove(segDir)
	_ = os.MkdirAll(segDir, 0o755)
	res, err := e.Query("m", nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) == 0 || len(res[0].Points) == 0 {
		t.Fatal("open window lost after failed Compact flush")
	}
	if res[0].Points[0].Value != 5 {
		t.Fatalf("value=%v", res[0].Points[0].Value)
	}
}
