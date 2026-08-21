package tsrollup_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup"
	"github.com/LYH2263/go-tsrollup/internal/agg"
	"github.com/LYH2263/go-tsrollup/internal/clock"
)

func openTest(t *testing.T, kind string, win time.Duration) *tsrollup.Engine {
	t.Helper()
	dir := t.TempDir()
	fake := &clock.Fake{T: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	e, err := tsrollup.Open(tsrollup.Options{
		Root:    dir,
		Window:  tsrollup.WindowSpec{Size: win},
		AggKind: kind,
		Clock:   fake,
		Factory: agg.DefaultFactory{},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = e.Close() })
	return e
}

func TestSumAvgMaxWindows(t *testing.T) {
	for _, kind := range []string{"sum", "avg", "max"} {
		kind := kind
		t.Run(kind, func(t *testing.T) {
			e := openTest(t, kind, time.Minute)
			base := time.Date(2024, 1, 1, 0, 0, 10, 0, time.UTC)
			vals := []float64{1, 3, 5}
			for i, v := range vals {
				if err := e.Append(tsrollup.Sample{
					Series: "cpu",
					Labels: map[string]string{"host": "a"},
					Ts:     base.Add(time.Duration(i) * time.Second),
					Value:  v,
				}); err != nil {
					t.Fatal(err)
				}
			}
			res, err := e.Query("cpu", map[string]string{"host": "a"}, time.Time{}, time.Time{})
			if err != nil {
				t.Fatal(err)
			}
			if len(res) != 1 || len(res[0].Points) != 1 {
				t.Fatalf("got %+v", res)
			}
			got := res[0].Points[0].Value
			var want float64
			switch kind {
			case "sum":
				want = 9
			case "avg":
				want = 3
			case "max":
				want = 5
			}
			if math.Abs(got-want) > 1e-9 {
				t.Fatalf("%s want %v got %v", kind, want, got)
			}
		})
	}
}

func TestAppendClonesLabels(t *testing.T) {
	e := openTest(t, "sum", time.Minute)
	labels := map[string]string{"host": "a"}
	if err := e.Append(tsrollup.Sample{
		Series: "mem",
		Labels: labels,
		Ts:     time.Date(2024, 1, 1, 0, 0, 1, 0, time.UTC),
		Value:  1,
	}); err != nil {
		t.Fatal(err)
	}
	labels["host"] = "mutated"
	list, err := e.ListSeries()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Labels["host"] != "a" {
		t.Fatalf("labels polluted: %+v", list)
	}
}

func TestQueryPointsIndependent(t *testing.T) {
	e := openTest(t, "sum", time.Minute)
	ts := time.Date(2024, 1, 1, 0, 0, 1, 0, time.UTC)
	_ = e.Append(tsrollup.Sample{Series: "x", Ts: ts, Value: 2})
	res, err := e.Query("x", nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	res[0].Points[0].Value = 999
	res2, err := e.Query("x", nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if res2[0].Points[0].Value != 2 {
		t.Fatalf("shared buffer: %v", res2[0].Points[0].Value)
	}
}

func TestAppendAfterClose(t *testing.T) {
	dir := t.TempDir()
	e, err := tsrollup.Open(tsrollup.Options{
		Root:   dir,
		Window: tsrollup.WindowSpec{Size: time.Minute},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := e.Close(); err != nil {
		t.Fatal(err)
	}
	err = e.Append(tsrollup.Sample{
		Series: "y",
		Ts:     time.Now().UTC(),
		Value:  1,
	})
	if err != tsrollup.ErrClosed {
		t.Fatalf("want ErrClosed got %v", err)
	}
}

func TestAppendContextCanceled(t *testing.T) {
	e := openTest(t, "sum", time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := e.AppendContext(ctx, tsrollup.Sample{
		Series: "c",
		Ts:     time.Date(2024, 1, 1, 0, 0, 1, 0, time.UTC),
		Value:  1,
	})
	if err == nil {
		t.Fatal("expected cancel error")
	}
}

func TestCloseFlushesLastWindow(t *testing.T) {
	dir := t.TempDir()
	e, err := tsrollup.Open(tsrollup.Options{
		Root:   dir,
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
	e2, err := tsrollup.Open(tsrollup.Options{
		Root:   dir,
		Window: tsrollup.WindowSpec{Size: time.Hour},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer e2.Close()
	res, err := e2.Query("last", nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	// WAL replay 也应恢复；段文件在 Close 时写入
	if len(res) == 0 {
		t.Fatal("expected recovered series")
	}
}

func TestWindowBounds(t *testing.T) {
	w := tsrollup.WindowSpec{Size: time.Minute}
	ts := time.Date(2024, 1, 1, 0, 0, 30, 0, time.UTC)
	start, end := w.Bounds(ts)
	if !start.Equal(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("start %v", start)
	}
	if !end.Equal(start.Add(time.Minute)) {
		t.Fatalf("end %v", end)
	}
}
