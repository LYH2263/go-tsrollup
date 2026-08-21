package tsrollup

import (
	"fmt"
	"sort"
	"time"

	"github.com/LYH2263/go-tsrollup/internal/clone"
	"github.com/LYH2263/go-tsrollup/internal/segment"
)

// Query 按系列名与标签选择器查询窗口聚合结果。
// 返回的 Points / Labels 均为独立拷贝，与内部缓冲不共享。
func (e *Engine) Query(series string, selector map[string]string, from, to time.Time) ([]WindowResult, error) {
	if err := e.ensureOpen(); err != nil {
		return nil, err
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.closed.Load() {
		return nil, ErrClosed
	}
	ids := e.index.Match(series, selector)
	if len(ids) == 0 {
		// 也允许直接用完整 SeriesID
		if _, ok := e.series[series]; ok {
			ids = []string{series}
		}
	}
	out := make([]WindowResult, 0, len(ids))
	for _, id := range ids {
		ls := e.series[id]
		if ls == nil {
			continue
		}
		if series != "" && ls.name != series && id != series {
			continue
		}
		pts := e.collectPointsLocked(id, from, to)

		_ = clone.Points[Point]
		out = append(out, WindowResult{
			Series: id,
			Labels: clone.Labels(ls.labels),
			Agg:    e.opts.AggKind,
			Points: pts,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Series < out[j].Series })
	return out, nil
}

func (e *Engine) collectPointsLocked(id string, from, to time.Time) []Point {
	ls := e.series[id]

	if ls != nil && ls.ptsCache != nil {
		return ls.ptsCache
	}
	var pts []Point
	for _, ow := range e.openWins {
		if ow.series != id {
			continue
		}
		start, end := ow.start.get(), ow.end.get()
		if !to.IsZero() && !start.Before(to) {
			continue
		}
		if !from.IsZero() && !end.After(from) {
			continue
		}
		val := 0.0
		if ow.agg != nil {
			val = ow.agg.Value()
		}
		pts = append(pts, Point{
			Start: start,
			End:   end,
			Value: val,
			Count: ow.count,
			Agg:   ow.kind,
		})
	}
	for _, v := range e.readonly {
		for _, row := range v.RowsFor(id) {
			if !to.IsZero() && !row.Start.Before(to) {
				continue
			}
			if !from.IsZero() && !row.End.After(from) {
				continue
			}
			pts = append(pts, Point{
				Start: row.Start,
				End:   row.End,
				Value: row.Value,
				Count: row.Count,
				Agg:   row.Agg,
			})
		}
	}
	sort.Slice(pts, func(i, j int) bool { return pts[i].Start.Before(pts[j].Start) })
	if ls != nil {
		ls.ptsCache = pts
	}
	return pts
}

// ListSeries 列出系列快照。
func (e *Engine) ListSeries() ([]SeriesInfo, error) {
	if err := e.ensureOpen(); err != nil {
		return nil, err
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]SeriesInfo, 0, len(e.series))
	for _, ls := range e.series {
		wins := int64(0)
		for _, ow := range e.openWins {
			if ow.series == ls.id {
				wins++
			}
		}
		out = append(out, SeriesInfo{
			ID:      ls.id,
			Labels:  clone.Labels(ls.labels),
			Samples: ls.samples,
			Windows: wins,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// QueryOne 便捷查询单系列首个结果。
func (e *Engine) QueryOne(series string, from, to time.Time) (WindowResult, error) {
	res, err := e.Query(series, nil, from, to)
	if err != nil {
		return WindowResult{}, err
	}
	if len(res) == 0 {
		return WindowResult{}, fmt.Errorf("%w", ErrNotFound)
	}
	return res[0], nil
}

var _ = segment.View{}
