package tsrollup

import (
	"github.com/LYH2263/go-tsrollup/internal/clone"
)

// Snapshot 导出当前打开窗口摘要（独立拷贝）。
type WindowSnap struct {
	Series string
	Start  int64
	End    int64
	Value  float64
	Count  int64
	Agg    string
}

func (e *Engine) Snapshot() ([]WindowSnap, error) {
	if err := e.ensureOpen(); err != nil {
		return nil, err
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]WindowSnap, 0, len(e.openWins))
	for _, ow := range e.openWins {
		val := 0.0
		if ow.agg != nil {
			val = ow.agg.Value()
		}
		out = append(out, WindowSnap{
			Series: ow.series,
			Start:  ow.start.get().UnixNano(),
			End:    ow.end.get().UnixNano(),
			Value:  val,
			Count:  ow.count,
			Agg:    clone.String(ow.kind),
		})
	}
	return out, nil
}
