package serieskey

import (
	"fmt"
	"strings"
)

// ParseLabels 从 Build 产物解析标签（简化解析，不处理转义嵌套）。
func ParseLabels(id string) (map[string]string, error) {
	i := strings.IndexByte(id, '{')
	if i < 0 {
		return map[string]string{}, nil
	}
	out := make(map[string]string)
	rest := id[i:]
	for len(rest) > 0 {
		if rest[0] != '{' {
			return nil, fmt.Errorf("bad label chunk")
		}
		end := strings.IndexByte(rest, '}')
		if end < 0 {
			return nil, fmt.Errorf("unclosed label")
		}
		pair := rest[1:end]
		eq := strings.IndexByte(pair, '=')
		if eq < 0 {
			return nil, fmt.Errorf("missing =")
		}
		out[pair[:eq]] = pair[eq+1:]
		rest = rest[end+1:]
	}
	return out, nil
}
