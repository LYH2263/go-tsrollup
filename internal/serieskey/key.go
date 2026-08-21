package serieskey

import (
	"sort"
	"strings"
)

// Build 规范化系列键。
func Build(name string, labels map[string]string) string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString(name)
	for _, k := range keys {
		b.WriteByte('{')
		b.WriteString(escape(k))
		b.WriteByte('=')
		b.WriteString(escape(labels[k]))
		b.WriteByte('}')
	}
	return b.String()
}

func escape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "{", "\\{")
	s = strings.ReplaceAll(s, "}", "\\}")
	s = strings.ReplaceAll(s, "=", "\\=")
	return s
}

// ParseName 从键提取指标名。
func ParseName(id string) string {
	if i := strings.IndexByte(id, '{'); i >= 0 {
		return id[:i]
	}
	return id
}
