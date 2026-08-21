package tsrollup_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/LYH2263/go-tsrollup"
)

func TestBug01_AppendLabelsSliceAlias(t *testing.T) {
	dir := t.TempDir()
	e, err := tsrollup.Open(tsrollup.Options{
		Root:   filepath.Join(dir, "e"),
		Window: tsrollup.WindowSpec{Size: time.Minute},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer e.Close()
	labels := map[string]string{"host": "a", "zone": "z1"}
	if err := e.Append(tsrollup.Sample{
		Series: "cpu",
		Labels: labels,
		Ts:     time.Date(2024, 1, 1, 0, 0, 1, 0, time.UTC),
		Value:  1,
	}); err != nil {
		t.Fatal(err)
	}
	labels["host"] = "MUTATED"
	list, err := e.ListSeries()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("series=%d", len(list))
	}
	if list[0].Labels["host"] != "a" {
		t.Fatalf("labels polluted by caller alias: %+v", list[0].Labels)
	}
}
