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

// normalizeLabels 校验并深拷贝标签，确保调用方事后修改原 map
// 不会回流到 rec / liveSeries.labels / 倒排索引内部状态。
func normalizeLabels(labels map[string]string) (map[string]string, error) {
	if err := validate.Labels(labels); err != nil {
		return nil, err
	}
	return clone.Labels(labels), nil
}

func labelMatch(have, want map[string]string) bool {
	for k, v := range want {
		if have[k] != v {
			return false
		}
	}
	return true
}
