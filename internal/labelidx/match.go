package labelidx

import "sort"

// Match 按指标名与标签选择器匹配系列 ID。
func (idx *Index) Match(name string, selector map[string]string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var cand map[string]struct{}
	if name != "" {
		cand = copySet(idx.byName[name])
		// 也尝试直接把 name 当完整 id
		if _, ok := idx.byID[name]; ok {
			if cand == nil {
				cand = make(map[string]struct{})
			}
			cand[name] = struct{}{}
		}
	} else {
		cand = make(map[string]struct{}, len(idx.byID))
		for id := range idx.byID {
			cand[id] = struct{}{}
		}
	}
	for k, v := range selector {
		set := idx.byLabel[k][v]
		cand = intersect(cand, set)
	}
	out := make([]string, 0, len(cand))
	for id := range cand {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

func copySet(in map[string]struct{}) map[string]struct{} {
	if in == nil {
		return nil
	}
	out := make(map[string]struct{}, len(in))
	for k := range in {
		out[k] = struct{}{}
	}
	return out
}

func intersect(a map[string]struct{}, b map[string]struct{}) map[string]struct{} {
	if a == nil {
		return nil
	}
	if b == nil {
		return map[string]struct{}{}
	}
	out := make(map[string]struct{})
	for id := range a {
		if _, ok := b[id]; ok {
			out[id] = struct{}{}
		}
	}
	return out
}
