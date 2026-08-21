package downsample

import (
	"sync"
	"time"
)

// Pyramid 多级降采样金字塔。
type Pyramid struct {
	mu     sync.Mutex
	base   time.Duration
	levels []*Level
}

func NewPyramid(base time.Duration) *Pyramid {
	if base <= 0 {
		base = time.Minute
	}
	p := &Pyramid{base: base}
	// 1x, 5x, 15x, 60x
	for _, m := range []int{1, 5, 15, 60} {
		p.levels = append(p.levels, NewLevel(base*time.Duration(m), m))
	}
	return p
}

func (p *Pyramid) Observe(series string, ts time.Time, v float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, lv := range p.levels {
		lv.Observe(series, ts, v)
	}
}

func (p *Pyramid) Levels() []LevelInfo {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]LevelInfo, len(p.levels))
	for i, lv := range p.levels {
		out[i] = lv.Info()
	}
	return out
}

func (p *Pyramid) Query(series string, level int, from, to time.Time) []Bucket {
	p.mu.Lock()
	defer p.mu.Unlock()
	if level < 0 || level >= len(p.levels) {
		return nil
	}
	return p.levels[level].Range(series, from, to)
}
