package tsrollup

import (
	"context"
	"fmt"
	"time"

	"github.com/LYH2263/go-tsrollup/internal/clone"
	"github.com/LYH2263/go-tsrollup/internal/validate"
	"github.com/LYH2263/go-tsrollup/internal/wal"
)

// Append 写入一条采样；会深拷贝 Labels，调用方事后修改不影响索引。
func (e *Engine) Append(s Sample) error {
	return e.AppendContext(context.Background(), s)
}

// AppendContext 带取消的写入；WAL 写前检查 ctx。
func (e *Engine) AppendContext(ctx context.Context, s Sample) error {

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: %v", ErrCanceled, err)
	}
	if err := validate.Sample(s.Series, s.Ts, s.Value); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	labels, err := normalizeLabels(s.Labels)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	labels = clone.Labels(labels)
	id := SeriesID(s.Series, labels)
	rec := wal.Record{
		Series: id,
		Name:   s.Series,
		Labels: labels,
		Ts:     s.Ts.UTC(),
		Value:  s.Value,
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: %v", ErrCanceled, err)
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if err := e.wal.Append(ctx, rec); err != nil {
		return fmt.Errorf("%w: %v", ErrWAL, err)
	}
	e.walRecs.Add(1)
	e.applyLocked(rec)
	e.metrics.IncSamples(1)
	e.samples.Add(1)
	return nil
}

func (e *Engine) applyLocked(rec wal.Record) {
	ls, ok := e.series[rec.Series]
	if !ok {
		ls = &liveSeries{
			id:     rec.Series,
			name:   rec.Name,
			labels: clone.Labels(rec.Labels),
		}
		e.series[rec.Series] = ls
		e.index.Add(rec.Series, rec.Labels)
	}
	ls.samples++
	start, end := e.opts.Window.Bounds(rec.Ts)
	wkey := windowKey(rec.Series, start)
	ow, ok := e.openWins[wkey]
	if !ok {
		a, err := e.factory.New(e.opts.AggKind)
		if err != nil || a == nil {
			// 工厂缺失时延后到 Compact 报 ErrNoAgg；此处仍建窗口占位。
			a = nil
		}
		ow = &openWindow{
			key:    wkey,
			series: rec.Series,
			kind:   e.opts.AggKind,
			agg:    a,
		}
		ow.start.set(start)
		ow.end.set(end)
		e.openWins[wkey] = ow
	}
	if ow.agg != nil {
		ow.agg.Add(rec.Value)
	}
	ow.count++
	if e.opts.Downsample {
		e.pyramid.Observe(rec.Series, rec.Ts, rec.Value)
	}
}

func windowKey(series string, start time.Time) string {
	return series + "@" + start.UTC().Format(time.RFC3339Nano)
}

func (t *atomicTime) set(tm time.Time) { t.nsec = tm.UnixNano() }
func (t *atomicTime) get() time.Time   { return time.Unix(0, t.nsec).UTC() }
