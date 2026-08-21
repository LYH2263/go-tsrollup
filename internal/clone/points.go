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

// Points 拷贝点切片（元素值拷贝），返回与入参不共享底层的独立切片。
// 调用方对返回值的任何改动都不会影响入参，反之亦然。
func Points[T any](in []T) []T {
	if in == nil {
		return nil
	}
	out := make([]T, len(in))
	copy(out, in)
	return out
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
