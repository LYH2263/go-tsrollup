package tsrollup_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup"
)

func TestBug02_QueryPointsSliceAlias(t *testing.T) {
	dir := t.TempDir()
	e, err := tsrollup.Open(tsrollup.Options{
		Root:   filepath.Join(dir, "e"),
		Window: tsrollup.WindowSpec{Size: time.Minute},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer e.Close()
	ts := time.Date(2024, 1, 1, 0, 0, 1, 0, time.UTC)
	if err := e.Append(tsrollup.Sample{Series: "x", Ts: ts, Value: 2}); err != nil {
		t.Fatal(err)
	}
	res, err := e.Query("x", nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || len(res[0].Points) == 0 {
		t.Fatalf("got %+v", res)
	}
	res[0].Points[0].Value = 999
	res2, err := e.Query("x", nil, time.Time{}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if res2[0].Points[0].Value != 2 {
		t.Fatalf("query points aliased: got %v want 2", res2[0].Points[0].Value)
	}
}
