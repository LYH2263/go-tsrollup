package clone

import "time"

// SampleLike 可拷贝采样字段。
type SampleLike struct {
	Series string
	Labels map[string]string
	Ts     time.Time
	Value  float64
}

func Sample(s SampleLike) SampleLike {
	return SampleLike{
		Series: s.Series,
		Labels: Labels(s.Labels),
		Ts:     s.Ts,
		Value:  s.Value,
	}
}
