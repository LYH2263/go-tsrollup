package encoding

import (
	"bytes"
	"encoding/json"
)

// CompactJSON 压缩 JSON。
func CompactJSON(in []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.Compact(&buf, in); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// PrettyJSON 缩进 JSON。
func PrettyJSON(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// DecodeMap 解码为 map。
func DecodeMap(b []byte) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}
