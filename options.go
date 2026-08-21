package tsrollup

import (
	"time"

	"github.com/LYH2263/go-tsrollup/internal/agg"
	"github.com/LYH2263/go-tsrollup/internal/clock"
	"github.com/LYH2263/go-tsrollup/internal/metrics"
)

// Options 打开引擎的配置。
type Options struct {
	Root           string
	Window         WindowSpec
	AggKind        string // sum | avg | max
	Factory        agg.Factory
	Clock          clock.Clock
	Metrics        *metrics.Registry
	MaxOpenWindows int
	WALSyncEvery   int
	SegmentMaxRows int
	Downsample     bool
}

func (o *Options) withDefaults() Options {
	out := *o
	if out.Window.Size <= 0 {
		out.Window = WindowSpec{Size: time.Minute, Align: AlignFloor}
	}
	if out.AggKind == "" {
		out.AggKind = "sum"
	}
	if out.Factory == nil {
		out.Factory = agg.DefaultFactory{}
	}
	if out.Clock == nil {
		out.Clock = clock.Real{}
	}
	if out.Metrics == nil {
		out.Metrics = metrics.NewRegistry()
	}
	if out.MaxOpenWindows <= 0 {
		out.MaxOpenWindows = 4096
	}
	if out.WALSyncEvery <= 0 {
		out.WALSyncEvery = 32
	}
	if out.SegmentMaxRows <= 0 {
		out.SegmentMaxRows = 1024
	}
	return out
}
