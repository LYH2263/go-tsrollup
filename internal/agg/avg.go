package agg

// Avg 均值聚合（维护 sum 与 count）。
type Avg struct {
	sum float64
	n   int64
}

func NewAvg() *Avg { return &Avg{} }

func (a *Avg) Add(v float64) { a.sum += v; a.n++ }

func (a *Avg) Value() float64 {
	if a.n == 0 {
		return 0
	}
	return a.sum / float64(a.n)
}

func (a *Avg) Count() int64 { return a.n }
func (a *Avg) Kind() string { return "avg" }
func (a *Avg) Reset()       { a.sum, a.n = 0, 0 }

func (a *Avg) Clone() Aggregator {
	c := *a
	return &c
}

func (a *Avg) Merge(other Aggregator) error {
	o, ok := other.(*Avg)
	if !ok {
		return errMergeType
	}
	a.sum += o.sum
	a.n += o.n
	return nil
}
