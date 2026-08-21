package hashutil

import (
	"hash/fnv"
	"sort"
)

// Labels 计算标签稳定哈希。
func Labels(name string, labels map[string]string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(name))
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		_, _ = h.Write([]byte{0})
		_, _ = h.Write([]byte(k))
		_, _ = h.Write([]byte{1})
		_, _ = h.Write([]byte(labels[k]))
	}
	return h.Sum64()
}

// String 计算字符串哈希。
func String(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
