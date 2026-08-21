package clone

// Labels 深拷贝标签 map，使调用方事后修改原 map 不影响副本。
// 对 nil 输入返回 nil，保持“零值无副作用”语义。
func Labels(in map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// String 返回字符串副本（避免共享底层？对 string 无必要，保持 API）。
func String(s string) string { return s + "" }
