package agg

// Sum 累加聚合。
type Sum struct {
	sum float64
	n   int64
}

func NewSum() *Sum { return &Sum{} }

func (s *Sum) Add(v float64) { s.sum += v; s.n++ }
func (s *Sum) Value() float64 { return s.sum }
func (s *Sum) Count() int64   { return s.n }
func (s *Sum) Kind() string   { return "sum" }
func (s *Sum) Reset()         { s.sum, s.n = 0, 0 }

func (s *Sum) Clone() Aggregator {
	c := *s
	return &c
}

func (s *Sum) Merge(other Aggregator) error {
	o, ok := other.(*Sum)
	if !ok {
		return errMergeType
	}
	s.sum += o.sum
	s.n += o.n
	return nil
}
