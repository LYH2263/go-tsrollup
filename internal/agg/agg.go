package agg

import "fmt"

// Aggregator 窗口聚合状态机。
type Aggregator interface {
	Add(v float64)
	Value() float64
	Count() int64
	Merge(other Aggregator) error
	Reset()
	Kind() string
	Clone() Aggregator
}

// Factory 创建聚合器。
type Factory interface {
	New(kind string) (Aggregator, error)
}

// DefaultFactory 支持 sum/avg/max。
type DefaultFactory struct{}

func (DefaultFactory) New(kind string) (Aggregator, error) {
	switch kind {
	case "", "sum":
		return NewSum(), nil
	case "avg":
		return NewAvg(), nil
	case "max":
		return NewMax(), nil
	default:
		return nil, fmt.Errorf("unknown agg kind %q", kind)
	}
}
