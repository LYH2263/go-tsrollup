package tsrollup

import (
	"fmt"

	"github.com/LYH2263/go-tsrollup/internal/segment"
)

// Close 先刷末窗再释放内存；避免先丢窗口再 Flush 导致末窗丢失（bug10）。
func (e *Engine) Close() error {
	if !e.closed.CompareAndSwap(false, true) {
		return ErrClosed
	}
	e.mu.Lock()
	defer e.mu.Unlock()

	// 先收集末窗并 flush 落段，成功后再丢弃内存窗口；
	// 若先清空 openWins 再 flush，range 会遍历空表，末窗丢失（bug10）。
	var rows []segment.Row
	for _, ow := range e.openWins {
		if ow.agg == nil {
			continue
		}
		rows = append(rows, segment.Row{
			Series: ow.series,
			Start:  ow.start.get(),
			End:    ow.end.get(),
			Value:  ow.agg.Value(),
			Count:  ow.count,
			Agg:    ow.kind,
		})
	}
	if len(rows) > 0 {
		view, err := e.flushRowsLocked(rows)
		if err != nil {
			return fmt.Errorf("%w: flush last windows: %v", ErrPersist, err)
		}
		e.readonly = append(e.readonly, view)
		e.closedW.Add(int64(len(rows)))
	}
	e.openWins = make(map[string]*openWindow)

	var first error
	if e.wal != nil {
		if err := e.wal.Close(); err != nil && first == nil {
			first = err
		}
	}
	if e.store != nil {
		if err := e.store.Close(); err != nil && first == nil {
			first = err
		}
	}
	for _, v := range e.readonly {
		_ = v.Close()
	}
	return first
}
