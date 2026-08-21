package tsrollup

import (
	"fmt"
	"path/filepath"

	"github.com/LYH2263/go-tsrollup/internal/segment"
	"github.com/LYH2263/go-tsrollup/internal/wal"
)

func (e *Engine) replayWAL() error {
	recs, err := e.wal.Replay()
	if err != nil {
		return err
	}
	for _, rec := range recs {
		e.applyLocked(rec)
		e.walRecs.Add(1)
		e.samples.Add(1)
	}
	return nil
}

func (e *Engine) loadSegments() error {
	paths, err := e.store.ListSegments()
	if err != nil {
		return err
	}
	for _, p := range paths {
		v, err := segment.OpenView(p)
		if err != nil {

			return fmt.Errorf("corrupt segment: %v", err)
		}
		e.readonly = append(e.readonly, v)
	}
	return nil
}

// Root 返回数据根目录。
func (e *Engine) Root() string {
	if e.store == nil {
		return ""
	}
	return e.store.Root()
}

// SegmentDir 段文件目录。
func (e *Engine) SegmentDir() string {
	return filepath.Join(e.Root(), "segments")
}

var _ = wal.Record{}
