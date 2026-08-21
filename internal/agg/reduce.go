package agg

// Reduce 将多个快照按 kind 归约。
func Reduce(kind string, snaps []Snapshot) (Snapshot, error) {
	a, err := DefaultFactory{}.New(kind)
	if err != nil {
		return Snapshot{}, err
	}
	for _, s := range snaps {
		switch kind {
		case "sum", "":
			a.Add(s.Value)
		case "avg":
			// 近似：按 count 展开不可行时按均值加权
			if s.Count <= 0 {
				continue
			}
			for i := int64(0); i < s.Count; i++ {
				a.Add(s.Value)
			}
		case "max":
			a.Add(s.Value)
		}
	}
	return TakeSnapshot(a), nil
}

// EqualKind 判断两个聚合器 kind 相同。
func EqualKind(a, b Aggregator) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Kind() == b.Kind()
}
