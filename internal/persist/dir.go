package persist

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Dir 引擎根目录管理。
type Dir struct {
	mu   sync.Mutex
	root string
	seq  int
}

func Open(root string) (*Dir, error) {
	for _, sub := range []string{"segments", "wal"} {
		if err := os.MkdirAll(filepath.Join(root, sub), 0o755); err != nil {
			return nil, err
		}
	}
	d := &Dir{root: root}
	paths, _ := d.ListSegments()
	d.seq = len(paths)
	return d, nil
}

func (d *Dir) Root() string { return d.root }

func (d *Dir) NextSegmentPath() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seq++
	name := fmt.Sprintf("seg-%06d-%d.tsr", d.seq, time.Now().UnixNano())
	return filepath.Join(d.root, "segments", name)
}

func (d *Dir) ListSegments() ([]string, error) {
	dir := filepath.Join(d.root, "segments")
	ents, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []string
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) == ".tsr" {
			out = append(out, filepath.Join(dir, name))
		}
	}
	return out, nil
}

func (d *Dir) Close() error { return nil }
