package tsrollup_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup"
	"github.com/LYH2263/go-tsrollup/internal/agg"
	"github.com/LYH2263/go-tsrollup/internal/clock"
)

type nilAggFactory struct{}

func (nilAggFactory) New(kind string) (agg.Aggregator, error) { return nil, nil }

func TestCompactNilAggregatorReturnsErrNoAgg(t *testing.T) {
	dir := t.TempDir()
	fake := &clock.Fake{T: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	e, err := tsrollup.Open(tsrollup.Options{
		Root:    filepath.Join(dir, "e"),
		Window:  tsrollup.WindowSpec{Size: time.Second},
		Factory: nilAggFactory{},
		Clock:   fake,
		AggKind: "sum",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer e.Close()
	if err := e.Append(tsrollup.Sample{
		Series: "n",
		Ts:     fake.Now(),
		Value:  1,
	}); err != nil {
		t.Fatal(err)
	}
	fake.Advance(2 * time.Second)
	err = e.Compact(context.Background())
	if err != tsrollup.ErrNoAgg {
		t.Fatalf("want ErrNoAgg got %v", err)
	}
}
