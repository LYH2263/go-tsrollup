package tsrollup

import (
	"fmt"
	"time"

	"github.com/LYH2263/go-tsrollup/internal/validate"
)

// AlignMode 控制窗口对齐方式。
type AlignMode int

const (
	AlignFloor AlignMode = iota
	AlignCeil
)

// WindowSpec 描述 tumbling 窗口。
type WindowSpec struct {
	Size  time.Duration
	Align AlignMode
	Offset time.Duration
}

// Bounds 返回包含 ts 的窗口 [start, end)。
func (w WindowSpec) Bounds(ts time.Time) (start, end time.Time) {
	if w.Size <= 0 {
		return ts, ts
	}
	t := ts.Add(-w.Offset).UnixNano()
	sz := w.Size.Nanoseconds()
	var base int64
	switch w.Align {
	case AlignCeil:
		base = ((t + sz - 1) / sz) * sz
	default:
		base = (t / sz) * sz
	}
	start = time.Unix(0, base).Add(w.Offset).UTC()
	end = start.Add(w.Size)
	return start, end
}

// Validate 校验窗口规格。
func (w WindowSpec) Validate() error {
	if err := validate.WindowSize(w.Size); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	return nil
}

// Contains 判断 ts 是否落在 [start, end)。
func Contains(start, end, ts time.Time) bool {
	return !ts.Before(start) && ts.Before(end)
}
