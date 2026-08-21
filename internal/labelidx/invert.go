package labelidx

import "sort"

// Keys 返回所有标签键。
func (idx *Index) Keys() []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	out := make([]string, 0, len(idx.byLabel))
	for k := range idx.byLabel {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Values 返回某标签键的全部取值。
func (idx *Index) Values(key string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	m := idx.byLabel[key]
	out := make([]string, 0, len(m))
	for v := range m {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}
