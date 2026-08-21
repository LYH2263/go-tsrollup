package buffer

import "time"

// TimeValue 带时间戳的值。
type TimeValue struct {
	Ts    time.Time
	Value float64
}

// TimeBuf 按时间追加的缓冲。
type TimeBuf struct {
	items []TimeValue
}

func (b *TimeBuf) Append(ts time.Time, v float64) {
	b.items = append(b.items, TimeValue{Ts: ts.UTC(), Value: v})
}

func (b *TimeBuf) Len() int { return len(b.items) }

func (b *TimeBuf) Range(from, to time.Time) []TimeValue {
	out := make([]TimeValue, 0, len(b.items))
	for _, it := range b.items {
		if !from.IsZero() && it.Ts.Before(from) {
			continue
		}
		if !to.IsZero() && !it.Ts.Before(to) {
			continue
		}
		out = append(out, it)
	}
	return out
}

func (b *TimeBuf) Clear() { b.items = b.items[:0] }

func (b *TimeBuf) Last() (TimeValue, bool) {
	if len(b.items) == 0 {
		return TimeValue{}, false
	}
	return b.items[len(b.items)-1], true
}
