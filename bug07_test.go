package tsrollup_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup"
)

func TestBug07_AppendHonorsContextCancel(t *testing.T) {
	dir := t.TempDir()
	e, err := tsrollup.Open(tsrollup.Options{
		Root:   filepath.Join(dir, "e"),
		Window: tsrollup.WindowSpec{Size: time.Minute},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer e.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = e.AppendContext(ctx, tsrollup.Sample{
		Series: "c",
		Ts:     time.Date(2024, 1, 1, 0, 0, 1, 0, time.UTC),
		Value:  1,
	})
	if err == nil {
		t.Fatal("expected cancel error")
	}
	if !errors.Is(err, tsrollup.ErrCanceled) && !errors.Is(err, context.Canceled) {
		t.Fatalf("want canceled, got %v", err)
	}
}
