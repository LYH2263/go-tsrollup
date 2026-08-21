package tsrollup_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup"
)

func TestBug05_CorruptSegmentWraps(t *testing.T) {
	dir := t.TempDir()
	root := filepath.Join(dir, "e")
	if err := os.MkdirAll(filepath.Join(root, "segments"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "wal"), 0o755); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(root, "segments", "seg-000001-1.tsr")
	if err := os.WriteFile(bad, []byte("not-a-segment"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := tsrollup.Open(tsrollup.Options{
		Root:   root,
		Window: tsrollup.WindowSpec{Size: time.Minute},
	})
	if err == nil {
		t.Fatal("expected corrupt error")
	}
	if !errors.Is(err, tsrollup.ErrCorruptSegment) {
		t.Fatalf("want ErrCorruptSegment wrap, got %v", err)
	}
}
