package tsrollup

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/LYH2263/go-tsrollup/internal/agg"
	"github.com/LYH2263/go-tsrollup/internal/clock"
	"github.com/LYH2263/go-tsrollup/internal/downsample"
	"github.com/LYH2263/go-tsrollup/internal/labelidx"
	"github.com/LYH2263/go-tsrollup/internal/metrics"
	"github.com/LYH2263/go-tsrollup/internal/persist"
	"github.com/LYH2263/go-tsrollup/internal/segment"
	"github.com/LYH2263/go-tsrollup/internal/wal"
)

// Engine 时序滚动聚合引擎。
type Engine struct {
	mu       sync.RWMutex
	closed   atomic.Bool
	opts     Options
	clock    clock.Clock
	metrics  *metrics.Registry
	factory  agg.Factory
	index    *labelidx.Index
	series   map[string]*liveSeries
	openWins map[string]*openWindow
	readonly []*segment.View
	wal      *wal.Log
	store    *persist.Dir
	pyramid  *downsample.Pyramid
	samples  atomic.Int64
	closedW  atomic.Int64
	compacts atomic.Int64
	walRecs  atomic.Int64
}

type liveSeries struct {
	id      string
	name    string
	labels  map[string]string
	samples int64
}

type openWindow struct {
	key    string
	series string
	start  atomicTime
	end    atomicTime
	agg    agg.Aggregator
	kind   string
	count  int64
}

type atomicTime struct {
	nsec int64
}

// Open 打开或创建引擎根目录。
func Open(opts Options) (*Engine, error) {
	o := opts.withDefaults()
	if o.Root == "" {
		return nil, fmt.Errorf("%w: empty root", ErrInvalid)
	}
	if err := o.Window.Validate(); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(o.Root, 0o755); err != nil {
		return nil, err
	}
	store, err := persist.Open(o.Root)
	if err != nil {
		return nil, err
	}
	wlog, err := wal.Open(filepath.Join(o.Root, "wal"))
	if err != nil {
		_ = store.Close()
		return nil, err
	}
	e := &Engine{
		opts:     o,
		clock:    o.Clock,
		metrics:  o.Metrics,
		factory:  o.Factory,
		index:    labelidx.New(),
		series:   make(map[string]*liveSeries),
		openWins: make(map[string]*openWindow),
		wal:      wlog,
		store:    store,
		pyramid:  downsample.NewPyramid(o.Window.Size),
	}
	if err := e.replayWAL(); err != nil {
		_ = e.Close()
		return nil, err
	}
	if err := e.loadSegments(); err != nil {
		_ = e.Close()
		return nil, err
	}
	return e, nil
}

func (e *Engine) ensureOpen() error {
	if e.closed.Load() {
		return ErrClosed
	}
	return nil
}
