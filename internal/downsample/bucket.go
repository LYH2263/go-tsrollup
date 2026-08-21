package downsample

import "time"

// Bucket 降采样桶。
type Bucket struct {
	Start time.Time
	End   time.Time
	Sum   float64
	Max   float64
	Count int64
	set   bool
}

func (b *Bucket) Add(v float64) {
	b.Sum += v
	if !b.set || v > b.Max {
		b.Max = v
		b.set = true
	}
	b.Count++
}

func (b *Bucket) Avg() float64 {
	if b.Count == 0 {
		return 0
	}
	return b.Sum / float64(b.Count)
}
