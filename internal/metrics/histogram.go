package metrics

import (
	"sort"
	"sync"
)

// Histogram 简易直方图（按桶上界计数）。
type Histogram struct {
	mu      sync.Mutex
	bounds  []float64
	counts  []int64
	sum     float64
	n       int64
}

func NewHistogram(bounds []float64) *Histogram {
	b := append([]float64(nil), bounds...)
	sort.Float64s(b)
	return &Histogram{bounds: b, counts: make([]int64, len(b)+1)}
}

func (h *Histogram) Observe(v float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	i := sort.SearchFloat64s(h.bounds, v)
	if i < len(h.bounds) && v > h.bounds[i] {
		i++
	}
	if i >= len(h.counts) {
		i = len(h.counts) - 1
	}
	h.counts[i]++
	h.sum += v
	h.n++
}

func (h *Histogram) Count() int64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.n
}

func (h *Histogram) Mean() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.n == 0 {
		return 0
	}
	return h.sum / float64(h.n)
}
