package policy

import "time"

// Retention 段保留策略。
type Retention struct {
	MaxAge      time.Duration
	MaxSegments int
	MaxBytes    int64
}

func DefaultRetention() Retention {
	return Retention{
		MaxAge:      24 * time.Hour,
		MaxSegments: 1024,
		MaxBytes:    1 << 30,
	}
}

func (r Retention) Keep(age time.Duration, segs int, bytes int64) bool {
	if r.MaxAge > 0 && age > r.MaxAge {
		return false
	}
	if r.MaxSegments > 0 && segs > r.MaxSegments {
		return false
	}
	if r.MaxBytes > 0 && bytes > r.MaxBytes {
		return false
	}
	return true
}

func (r Retention) ClampSegments(n int) int {
	if r.MaxSegments <= 0 {
		return n
	}
	if n > r.MaxSegments {
		return r.MaxSegments
	}
	return n
}
