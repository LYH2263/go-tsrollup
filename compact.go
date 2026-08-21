package tsrollup

import (
	"context"
	"fmt"
	"time"

	"github.com/LYH2263/go-tsrollup/internal/segment"
	"github.com/LYH2263/go-tsrollup/internal/wait"
)

// Compact 将已关闭窗口刷入段文件，并提升只读视图。
func (e *Engine) Compact(ctx context.Context) error {
	if err := e.ensureOpen(); err != nil {
		return err
	}
	if err := wait.Lock(ctx, &e.mu); err != nil {
		return fmt.Errorf("%w: %v", ErrCanceled, err)
	}
	defer e.mu.Unlock()
	if e.closed.Load() {
		return ErrClosed
	}
	if e.factory == nil {
		return ErrNoAgg
	}
	now := e.clock.Now().UTC()
	var rows []segment.Row
	var closeKeys []string
	for k, ow := range e.openWins {
		if ow.agg == nil {
			return ErrNoAgg
		}
		end := ow.end.get()
		if !end.After(now) {
			rows = append(rows, segment.Row{
				Series: ow.series,
				Start:  ow.start.get(),
				End:    end,
				Value:  ow.agg.Value(),
				Count:  ow.count,
				Agg:    ow.kind,
			})
			closeKeys = append(closeKeys, k)
		}
	}
	if len(rows) == 0 {
		return nil
	}
	// 刷盘失败时不得删除 openWins，也不得提升只读视图：
	// 否则未落盘的内存聚合点将丢失，同进程再 Query 也读不到（bug06 目标）。
	// 窗口留在 openWins 内，等待下一次 Compact 重试刷盘。
	view, err := e.flushRowsLocked(rows)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPersist, err)
	}
	e.readonly = append(e.readonly, view)
	for _, k := range closeKeys {
		delete(e.openWins, k)
		e.closedW.Add(1)
	}
	e.compacts.Add(1)
	e.metrics.IncCompactions(1)
	_ = e.store.WriteManifest(len(e.readonly))
	return nil
}

func (e *Engine) flushRowsLocked(rows []segment.Row) (*segment.View, error) {
	path := e.store.NextSegmentPath()
	w, err := segment.Create(path)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		if err := w.WriteRow(r); err != nil {
			_ = w.Abort()
			return nil, err
		}
	}
	// Create→Close 内部 Sync 后再 Rename（bug09 目标）。
	if err := w.Close(); err != nil {
		return nil, err
	}
	v, err := segment.OpenView(path)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// ForceCloseWindow 测试辅助：强制关闭指定系列当前窗。
func (e *Engine) ForceCloseWindow(series string, at time.Time) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	start, _ := e.opts.Window.Bounds(at)
	k := windowKey(series, start)
	ow, ok := e.openWins[k]
	if !ok {
		return ErrNotFound
	}
	if ow.agg == nil {
		return ErrNoAgg
	}
	rows := []segment.Row{{
		Series: ow.series,
		Start:  ow.start.get(),
		End:    ow.end.get(),
		Value:  ow.agg.Value(),
		Count:  ow.count,
		Agg:    ow.kind,
	}}
	view, err := e.flushRowsLocked(rows)
	if err != nil {
		return err
	}
	e.readonly = append(e.readonly, view)
	delete(e.openWins, k)
	e.closedW.Add(1)
	return nil
}
