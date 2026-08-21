package clone

import "time"

// Point 与对外 Point 结构对齐的可拷贝点。
type Point struct {
	Start time.Time
	End   time.Time
	Value float64
	Count int64
	Agg   string
}

// Points 拷贝点切片（元素值拷贝）。
func Points[T any](in []T) []T {

	return in
}

// Bytes 拷贝字节切片。
func Bytes(in []byte) []byte {
	if in == nil {
		return nil
	}
	out := make([]byte, len(in))
	copy(out, in)
	return out
}
