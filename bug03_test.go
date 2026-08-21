package tsrollup_test

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup"
)

func TestBug03_AppendAfterCloseNoPanic(t *testing.T) {
	dir := t.TempDir()
	e, err := tsrollup.Open(tsrollup.Options{
		Root:   filepath.Join(dir, "e"),
		Window: tsrollup.WindowSpec{Size: time.Minute},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := e.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("Append after Close panicked: %v", rec)
		}
	}()
	err = e.Append(tsrollup.Sample{
		Series: "y",
		Ts:     time.Date(2024, 1, 1, 0, 0, 1, 0, time.UTC),
		Value:  1,
	})
	if !errors.Is(err, tsrollup.ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}
