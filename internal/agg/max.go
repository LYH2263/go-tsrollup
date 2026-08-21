package agg

import "math"

// Max 最大值聚合。
type Max struct {
	max float64
	n   int64
	set bool
}

func NewMax() *Max { return &Max{max: math.Inf(-1)} }

func (m *Max) Add(v float64) {
	if !m.set || v > m.max {
		m.max = v
		m.set = true
	}
	m.n++
}

func (m *Max) Value() float64 {
	if !m.set {
		return 0
	}
	return m.max
}

func (m *Max) Count() int64 { return m.n }
func (m *Max) Kind() string { return "max" }

func (m *Max) Reset() {
	m.max = math.Inf(-1)
	m.n = 0
	m.set = false
}

func (m *Max) Clone() Aggregator {
	c := *m
	return &c
}

func (m *Max) Merge(other Aggregator) error {
	o, ok := other.(*Max)
	if !ok {
		return errMergeType
	}
	if o.set && (!m.set || o.max > m.max) {
		m.max = o.max
		m.set = true
	}
	m.n += o.n
	return nil
}
