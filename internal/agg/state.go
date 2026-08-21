package agg

import "errors"

var errMergeType = errors.New("agg: merge type mismatch")

// Snapshot 只读聚合快照。
type Snapshot struct {
	Kind  string
	Value float64
	Count int64
}

func TakeSnapshot(a Aggregator) Snapshot {
	if a == nil {
		return Snapshot{}
	}
	return Snapshot{Kind: a.Kind(), Value: a.Value(), Count: a.Count()}
}

// Combine 将同 kind 的聚合器合并到 dst。
func Combine(dst Aggregator, srcs ...Aggregator) error {
	for _, s := range srcs {
		if s == nil {
			continue
		}
		if err := dst.Merge(s); err != nil {
			return err
		}
	}
	return nil
}
