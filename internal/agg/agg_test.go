package agg_test

import (
	"math"
	"testing"

	"github.com/LYH2263/go-tsrollup/internal/agg"
)

func TestSumAvgMax(t *testing.T) {
	s := agg.NewSum()
	a := agg.NewAvg()
	m := agg.NewMax()
	for _, v := range []float64{2, 4, 6} {
		s.Add(v)
		a.Add(v)
		m.Add(v)
	}
	if s.Value() != 12 {
		t.Fatalf("sum %v", s.Value())
	}
	if math.Abs(a.Value()-4) > 1e-9 {
		t.Fatalf("avg %v", a.Value())
	}
	if m.Value() != 6 {
		t.Fatalf("max %v", m.Value())
	}
}

func TestFactory(t *testing.T) {
	for _, k := range agg.Kinds() {
		a, err := agg.DefaultFactory{}.New(k)
		if err != nil || a == nil {
			t.Fatalf("%s: %v", k, err)
		}
	}
}
