package clone

// DupMapString 同 Labels，保留别名便于调用。
func DupMapString(in map[string]string) map[string]string {
	return Labels(in)
}

// DupStringSlice 拷贝字符串切片。
func DupStringSlice(in []string) []string {
	if in == nil {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}
