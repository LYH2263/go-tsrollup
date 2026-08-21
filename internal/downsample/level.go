package downsample

import (
	"sort"
	"time"
)

// Level 单层降采样。
type Level struct {
	step   time.Duration
	factor int
	series map[string]map[int64]*Bucket
}

type LevelInfo struct {
	Step   time.Duration
	Factor int
	Buckets int
}

func NewLevel(step time.Duration, factor int) *Level {
	return &Level{
		step:   step,
		factor: factor,
		series: make(map[string]map[int64]*Bucket),
	}
}

func (l *Level) Observe(series string, ts time.Time, v float64) {
	m, ok := l.series[series]
	if !ok {
		m = make(map[int64]*Bucket)
		l.series[series] = m
	}
	key := ts.UnixNano() / l.step.Nanoseconds()
	b, ok := m[key]
	if !ok {
		start := time.Unix(0, key*l.step.Nanoseconds()).UTC()
		b = &Bucket{Start: start, End: start.Add(l.step)}
		m[key] = b
	}
	b.Add(v)
}

func (l *Level) Range(series string, from, to time.Time) []Bucket {
	m := l.series[series]
	if m == nil {
		return nil
	}
	out := make([]Bucket, 0, len(m))
	for _, b := range m {
		if !to.IsZero() && !b.Start.Before(to) {
			continue
		}
		if !from.IsZero() && !b.End.After(from) {
			continue
		}
		cp := *b
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Start.Before(out[j].Start) })
	return out
}

func (l *Level) Info() LevelInfo {
	n := 0
	for _, m := range l.series {
		n += len(m)
	}
	return LevelInfo{Step: l.step, Factor: l.factor, Buckets: n}
}
