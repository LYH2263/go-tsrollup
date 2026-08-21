package queryplan

import (
	"sort"
	"time"
)

// Point 计划层点。
type Point struct {
	Start time.Time
	End   time.Time
	Value float64
	Count int64
}

// MergeSorted 合并已排序点列（同 Start 相加）。
func MergeSorted(a, b []Point) []Point {
	out := make([]Point, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i].Start.Before(b[j].Start):
			out = append(out, a[i])
			i++
		case b[j].Start.Before(a[i].Start):
			out = append(out, b[j])
			j++
		default:
			out = append(out, Point{
				Start: a[i].Start,
				End:   a[i].End,
				Value: a[i].Value + b[j].Value,
				Count: a[i].Count + b[j].Count,
			})
			i++
			j++
		}
	}
	out = append(out, a[i:]...)
	out = append(out, b[j:]...)
	return out
}

// SortByStart 就地排序。
func SortByStart(pts []Point) {
	sort.Slice(pts, func(i, j int) bool { return pts[i].Start.Before(pts[j].Start) })
}
