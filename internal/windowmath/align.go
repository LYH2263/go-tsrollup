package windowmath

import "time"

// Floor 对齐到 window 下界。
func Floor(ts time.Time, size time.Duration, offset time.Duration) time.Time {
	t := ts.Add(-offset).UnixNano()
	sz := size.Nanoseconds()
	base := (t / sz) * sz
	return time.Unix(0, base).Add(offset).UTC()
}

// Ceil 对齐到 window 上界起点。
func Ceil(ts time.Time, size time.Duration, offset time.Duration) time.Time {
	t := ts.Add(-offset).UnixNano()
	sz := size.Nanoseconds()
	base := ((t + sz - 1) / sz) * sz
	return time.Unix(0, base).Add(offset).UTC()
}

// Span 返回 [start,end)。
func Span(ts time.Time, size time.Duration, offset time.Duration) (time.Time, time.Time) {
	start := Floor(ts, size, offset)
	return start, start.Add(size)
}

// Overlaps 判断两窗是否相交。
func Overlaps(a0, a1, b0, b1 time.Time) bool {
	return a0.Before(b1) && b0.Before(a1)
}

// ClampRange 将 [from,to) 裁剪到 [lo,hi)。
func ClampRange(from, to, lo, hi time.Time) (time.Time, time.Time, bool) {
	if from.Before(lo) {
		from = lo
	}
	if to.After(hi) {
		to = hi
	}
	if !from.Before(to) {
		return time.Time{}, time.Time{}, false
	}
	return from, to, true
}
