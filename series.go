package tsrollup

import (
	"sort"
	"strings"

	"github.com/LYH2263/go-tsrollup/internal/clone"
	"github.com/LYH2263/go-tsrollup/internal/validate"
)

// SeriesID 由指标名与规范化标签派生。
func SeriesID(name string, labels map[string]string) string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString(name)
	for _, k := range keys {
		b.WriteByte('{')
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(labels[k])
		b.WriteByte('}')
	}
	return b.String()
}

func normalizeLabels(labels map[string]string) (map[string]string, error) {
	if err := validate.Labels(labels); err != nil {
		return nil, err
	}
	out := clone.Labels(labels)
	if out == nil {
		// 调用方只给系列名（Labels=nil）时 clone.Labels 返回 nil；
		// 收成非空 map，避免后续盖 __queried__ 戳时对 nil map 写入 panic。
		out = make(map[string]string)
	}
	return out, nil
}

func labelMatch(have, want map[string]string) bool {
	for k, v := range want {
		if have[k] != v {
			return false
		}
	}
	return true
}
