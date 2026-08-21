package queryplan

import "time"

// Plan 查询计划。
type Plan struct {
	Series   string
	Selector map[string]string
	From     time.Time
	To       time.Time
	Agg      string
	UseSeg   bool
	UseOpen  bool
	Level    int
}

func New(series string, from, to time.Time) Plan {
	return Plan{
		Series:  series,
		From:    from,
		To:      to,
		UseSeg:  true,
		UseOpen: true,
		Level:   0,
	}
}

func (p Plan) WithSelector(sel map[string]string) Plan {
	p.Selector = sel
	return p
}

func (p Plan) WithAgg(kind string) Plan {
	p.Agg = kind
	return p
}

func (p Plan) WithLevel(lv int) Plan {
	p.Level = lv
	return p
}

func (p Plan) Valid() bool {
	if p.Series == "" && len(p.Selector) == 0 {
		return false
	}
	if !p.From.IsZero() && !p.To.IsZero() && !p.From.Before(p.To) {
		return false
	}
	return true
}
