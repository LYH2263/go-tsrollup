package tsrollup_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup"
	"github.com/LYH2263/go-tsrollup/internal/agg"
	"github.com/LYH2263/go-tsrollup/internal/clock"
)

func TestBug04_NilLabelsQueryNoPanic(t *testing.T) {
	dir := t.TempDir()
	fake := &clock.Fake{T: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	e, err := tsrollup.Open(tsrollup.Options{
		Root:    filepath.Join(dir, "e"),
		Window:  tsrollup.WindowSpec{Size: time.Minute},
		Factory: agg.DefaultFactory{},
		Clock:   fake,
		AggKind: "sum",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer e.Close()
	if err := e.Append(tsrollup.Sample{
		Series: "cpu",
		Labels: nil,
		Ts:     fake.Now(),
		Value:  1,
	}); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("Query panicked on nil Labels series: %v", rec)
		}
	}()
	res, err := e.Query("cpu", nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("want 1 series result, got %+v", res)
	}
	if res[0].Labels == nil {
		t.Fatal("Labels should be non-nil map after query stamp")
	}
}
