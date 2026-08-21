package windowmath

import "time"

// Iter 按 tumbling 窗口迭代。
type Iter struct {
	cur  time.Time
	end  time.Time
	size time.Duration
}

func NewIter(from, to time.Time, size time.Duration) *Iter {
	start := Floor(from, size, 0)
	return &Iter{cur: start, end: to, size: size}
}

func (it *Iter) Next() (start, end time.Time, ok bool) {
	if !it.cur.Before(it.end) {
		return time.Time{}, time.Time{}, false
	}
	start = it.cur
	end = start.Add(it.size)
	it.cur = end
	return start, end, true
}

// CountWindows 估算区间内窗口数。
func CountWindows(from, to time.Time, size time.Duration) int {
	if !from.Before(to) || size <= 0 {
		return 0
	}
	n := 0
	it := NewIter(from, to, size)
	for {
		_, _, ok := it.Next()
		if !ok {
			break
		}
		n++
	}
	return n
}
