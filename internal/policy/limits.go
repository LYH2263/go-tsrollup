package policy

// Limits 引擎运行限额。
type Limits struct {
	MaxSeries      int
	MaxLabelPairs  int
	MaxOpenWindows int
	MaxAppendBatch int
	MaxQueryPoints int
}

func DefaultLimits() Limits {
	return Limits{
		MaxSeries:      100000,
		MaxLabelPairs:  64,
		MaxOpenWindows: 4096,
		MaxAppendBatch: 1024,
		MaxQueryPoints: 10000,
	}
}

func (l Limits) AllowSeries(n int) bool {
	return l.MaxSeries <= 0 || n < l.MaxSeries
}

func (l Limits) AllowWindows(n int) bool {
	return l.MaxOpenWindows <= 0 || n < l.MaxOpenWindows
}

func (l Limits) ClampQuery(n int) int {
	if l.MaxQueryPoints <= 0 || n <= l.MaxQueryPoints {
		return n
	}
	return l.MaxQueryPoints
}
