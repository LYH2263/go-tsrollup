package clock

import "time"

// Clock 可注入时钟。
type Clock interface {
	Now() time.Time
}

type Real struct{}

func (Real) Now() time.Time { return time.Now().UTC() }

// Fake 测试时钟。
type Fake struct {
	T time.Time
}

func (f *Fake) Now() time.Time { return f.T.UTC() }

func (f *Fake) Advance(d time.Duration) { f.T = f.T.Add(d) }
